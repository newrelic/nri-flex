package parser

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/newrelic/nri-flex/internal/formatter"
	"github.com/newrelic/nri-flex/internal/load"
	"github.com/newrelic/nri-flex/internal/logger"

	"github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// RunCommands executes the given commands to create one merged sampled
func RunCommands(yml load.Config, api load.API, dataStore *[]interface{}) {
	commandShell := load.DefaultShell
	dataSample := map[string]interface{}{}
	processedCols := false
	processedJMX := false
	for _, command := range api.Commands {
		if command.Run != "" {
			runCommand := command.Run
			if command.Output == load.Jmx {
				SetJMXCommand(&runCommand, command, api, yml)
			}
			// Create a new context and add a timeout to it
			ctx, cancel := context.WithTimeout(context.Background(), load.DefaultTimeout)
			defer cancel() // The cancel should be deferred so resources are cleaned up

			if api.Shell != "" {
				commandShell = api.Shell
			}
			if command.Shell != "" {
				commandShell = command.Shell
			}

			secondParameter := "-c"

			// windows commands are untested currently
			if runtime.GOOS == "windows" {
				commandShell = "cmd"
				secondParameter = "/C"
			}

			// Create the command with our context
			cmd := exec.CommandContext(ctx, commandShell, secondParameter, fmt.Sprintf("%v", runCommand))
			output, err := cmd.CombinedOutput()

			if err != nil {
				message := "command failed: " + command.Run
				if output != nil {
					message = message + " " + string(output)
				}
				logger.Flex("debug", err, message, false)
			} else if ctx.Err() == context.DeadlineExceeded {
				logger.Flex("debug", ctx.Err(), "command timed out", false)
			} else if ctx.Err() != nil {
				logger.Flex("debug", err, "command execution failed", false)
			} else {
				dataOutput := string(output)
				commandOutput, dataInterface := detectCommandOutput(dataOutput, command.Output)
				if !command.IgnoreOutput {
					switch commandOutput {
					case "raw":
						logger.Flex("info", nil, fmt.Sprintf("running %v", command.Run), false)
						if command.Split == "" { // default vertical split
							processRaw(&dataSample, dataOutput, command.SplitBy, command.LineLimit)
						} else if command.Split == "column" || command.Split == "horizontal" {
							if processedCols {
								logger.Flex("debug", fmt.Errorf("horizonal split only allowed once per command set %v %v", api.Name, command.Name), "", false)
							} else {
								processedCols = true
								processRawCol(dataStore, &dataSample, dataOutput, command)
							}
						}
					case load.JSONType:
						*dataStore = append(*dataStore, dataInterface)
					case load.Jmx:
						processedJMX = true
						ParseJMX(dataInterface, dataStore, command, dataSample)
					}
				}
			}
		}
	}
	// only send dataSample back, not if horizontal split or jmx was processed
	// this can probably be shuffled elsewhere
	if len(dataSample) > 0 && !processedCols && !processedJMX {
		*dataStore = append(*dataStore, dataSample)
	}
}

// processRaw processes a raw data output
func processRaw(dataSample *map[string]interface{}, dataOutput string, splitBy string, lineLimit int) {
	// SplitBy key is required else we cannot easily distinguish between keys and values
	for i, line := range strings.Split(strings.TrimSuffix(dataOutput, "\n"), "\n") {
		if i >= lineLimit && lineLimit != 0 {
			logger.Flex("info", nil, fmt.Sprintf("reached line limit %d", lineLimit), false)
			break
		}
		key, val, success := formatter.SplitKey(line, splitBy)
		if success {
			(*dataSample)[key] = strings.TrimRight(val, "\r\n") //line endings appear so we trim them
		}
	}
}

func processRawCol(dataStore *[]interface{}, dataSample *map[string]interface{}, dataOutput string, command load.Command) {
	headerLine := 0
	startLine := 1

	if command.RowHeader != 0 {
		headerLine = command.RowHeader
	}
	if command.RowStart != headerLine && command.RowStart >= 1 {
		startLine = command.RowStart
	}

	lines := strings.Split(strings.TrimSuffix(dataOutput, "\n"), "\n")
	header := lines[headerLine]
	var keys []string

	// set header keys
	if len(command.SetHeader) > 0 {
		keys = command.SetHeader
	} else {
		if command.HeaderRegexMatch {
			keys = append(keys, formatter.RegMatch(header, command.HeaderSplitBy)...)
		} else {
			keys = append(keys, formatter.RegSplit(header, command.HeaderSplitBy)...)
		}
	}

	for i, line := range lines {
		if i != headerLine && i >= startLine {
			if i >= command.LineLimit && command.LineLimit != 0 {
				logger.Flex("info", nil, fmt.Sprintf("reached line limit %d", command.LineLimit), false)
				break
			}

			cmdSample := map[string]interface{}{}

			// values contains the row values split
			var values []string
			if command.RegexMatch {
				values = formatter.RegMatch(line, command.SplitBy)
			} else {
				values = formatter.RegSplit(line, command.SplitBy)
			}

			// loop through header keys to apply values
			for index, key := range keys {
				if index+1 <= len(values) { // there can be items that exist past this, added this in because of docker ps example whilst testing
					cmdSample[key] = values[index]
				}
			}

			// add attributes from previously run commands into this cmdSample
			for key, val := range *dataSample {
				cmdSample[key] = val
			}

			*dataStore = append(*dataStore, cmdSample)
		}
	}
}

// detectCommandOutput currently only supports checking if json output
func detectCommandOutput(dataOutput string, commandOutput string) (string, interface{}) {
	if commandOutput == load.Jmx {
		dataOutputLines := strings.Split(strings.TrimSuffix(dataOutput, "\n"), "\n")
		startLine := 0
		endLine := 0
		startSet := false
		endSet := false
		for i, line := range dataOutputLines {
			if strings.HasPrefix(line, `{"`) {
				startLine = i
				startSet = true
			}
			if strings.HasSuffix(line, `}`) {
				endLine = i
				endSet = true
			}
		}
		if startLine == endLine && startSet && endSet {
			jmxDataOutput := dataOutputLines[startLine]
			var f interface{}
			err := json.Unmarshal([]byte(jmxDataOutput), &f)
			if err == nil {
				return load.Jmx, f
			}
			logger.Flex("debug", err, "failed to unmarshal jmx output", false)
		}
		return load.Jmx, map[string]interface{}{"error": "Failed to process JMX Data"}
	}

	// check json
	var f interface{}
	err := json.Unmarshal([]byte(dataOutput), &f)
	if err == nil {
		return load.JSONType, f
	}

	// default raw
	return "raw", nil
}
