package inputs

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/newrelic/nri-flex/internal/formatter"
	"github.com/newrelic/nri-flex/internal/load"
	"github.com/newrelic/nri-flex/internal/logger"
)

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// RunCommands executes the given commands to create one merged sampled
func RunCommands(dataStore *[]interface{}, yml *load.Config, apiNo int) {
	startTime := makeTimestamp()
	logger.Flex("debug", nil, fmt.Sprintf("%v - running commands", yml.Name), false)
	api := yml.APIs[apiNo]
	commandShell := load.DefaultShell
	dataSample := map[string]interface{}{}
	processType := ""
	for _, command := range api.Commands {
		if command.Run != "" && command.Dial == "" {
			runCommand := command.Run
			if command.Output == load.Jmx {
				SetJMXCommand(dataStore, &runCommand, command, api, yml)
			}
			commandTimeout := load.DefaultTimeout
			if api.Timeout > 0 {
				commandTimeout = time.Duration(api.Timeout) * time.Millisecond
			}
			if command.Timeout > 0 {
				commandTimeout = time.Duration(command.Timeout) * time.Millisecond
			}

			// Create a new context and add a timeout to it
			ctx, cancel := context.WithTimeout(context.Background(), commandTimeout)
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
				if command.SplitOutput != "" {
					splitOutput(dataStore, string(output), command, startTime)
				} else {
					processOutput(dataStore, string(output), &dataSample, command, api, &processType)
				}
			}
		} else if command.Cache != "" {
			if yml.Datastore[command.Cache] != nil {
				for _, cache := range yml.Datastore[command.Cache] {
					switch sample := cache.(type) {
					case map[string]interface{}:
						if sample["http"] != nil {
							logger.Flex("debug", nil, fmt.Sprintf("processing http cache with command processor %v", command.Cache), false)
							if command.SplitOutput != "" {
								splitOutput(dataStore, sample["http"].(string), command, startTime)
							} else {
								processOutput(dataStore, sample["http"].(string), &dataSample, command, api, &processType)
							}
						}
					}
				}
			}
		} else if command.Dial != "" {
			NetDialWithTimeout(dataStore, command, &dataSample, api, &processType)
		} else if command.ContainerExec != "" {
			// handle commands against containers
			if yml.CustomAttributes != nil {
				if yml.CustomAttributes["containerId"] != "" {
					logger.Flex("debug", nil, "not handled yet", false)
				}
			}
		}
	}
	// only send dataSample back, not if horizontal (columns) split or jmx was processed
	// this can probably be shuffled elsewhere
	if len(dataSample) > 0 && processType != load.TypeColumns && processType != "jmx" {
		dataSample["flex.commandTimeMs"] = makeTimestamp() - startTime
		*dataStore = append(*dataStore, dataSample)
	}
}

func splitOutput(dataStore *[]interface{}, output string, command load.Command, startTime int64) {
	lines := strings.Split(strings.TrimSuffix(output, "\n"), "\n")
	outputBlocks := [][]string{}
	startSplit := -1
	endSplit := 0
	for i, line := range lines {
		if formatter.KvFinder("regex", line, command.SplitOutput) {
			if startSplit == -1 {
				startSplit = i
			} else {
				endSplit = i
				outputBlocks = append(outputBlocks, lines[startSplit:endSplit])
				startSplit = i
			}
		}
		//create the last block
		if i+1 == len(lines) {
			outputBlocks = append(outputBlocks, lines[startSplit:i+1])
		}
	}
	processBlocks(dataStore, outputBlocks, command, startTime)
}

func processBlocks(dataStore *[]interface{}, blocks [][]string, command load.Command, startTime int64) {
	for _, block := range blocks {
		sample := map[string]interface{}{}
		regmatchCount := 0
		for _, regmatch := range command.RegexMatches {
			for _, line := range block {
				matches := formatter.RegMatch(line, regmatch.Expression)
				if len(matches) > 0 {
					for i, match := range matches {
						if len(regmatch.Keys) > 0 {
							key := regmatch.Keys[i]
							if len(regmatch.KeysMulti) > 0 {
								key = regmatch.KeysMulti[regmatchCount] + key
							}
							sample[key] = match
						}
					}
					regmatchCount++
				}
			}
			regmatchCount = 0
		}
		sample["flex.commandTimeMs"] = makeTimestamp() - startTime
		*dataStore = append(*dataStore, sample)
	}
}

func processOutput(dataStore *[]interface{}, output string, dataSample *map[string]interface{}, command load.Command, api load.API, processType *string) {
	dataOutput := output
	commandOutput, dataInterface := detectCommandOutput(dataOutput, command.Output)
	if !command.IgnoreOutput {
		switch commandOutput {
		case "raw":
			cmd := command.Run
			if command.Cache != "" {
				cmd = "cache - " + command.Cache
			}
			logger.Flex("debug", nil, fmt.Sprintf("running %v", cmd), false)
			if command.Split == "" { // default vertical split
				applyCustomAttributes(dataSample, &command.CustomAttributes)
				processRaw(dataSample, dataOutput, command.SplitBy, command.LineStart, command.LineEnd)
			} else if command.Split == load.TypeColumns || command.Split == "horizontal" {
				if *processType == load.TypeColumns {
					logger.Flex("debug", fmt.Errorf("horizonal split only allowed once per command set %v %v", api.Name, command.Name), "", false)
				} else {
					*processType = "columns"
					processRawCol(dataStore, dataSample, dataOutput, command)
				}
			}
		case load.TypeJSON:
			// load.StoreAppend(dataInterface)
			*dataStore = append(*dataStore, dataInterface)
		case load.Jmx:
			*processType = "jmx"
			ParseJMX(dataStore, dataInterface, command, dataSample)
		}
	}
}

// processRaw processes a raw data output
func processRaw(dataSample *map[string]interface{}, dataOutput string, splitBy string, lineStart int, lineEnd int) {
	// SplitBy key is required else we cannot easily distinguish between keys and values
	for i, line := range strings.Split(strings.TrimSuffix(dataOutput, "\n"), "\n") {
		if i >= lineStart {
			if i >= lineEnd && lineEnd != 0 {
				logger.Flex("debug", nil, fmt.Sprintf("reached line limit %d", lineEnd), false)
				break
			}
			key, val, success := formatter.SplitKey(line, splitBy)
			if success {
				(*dataSample)[key] = strings.TrimRight(val, "\r\n") //line endings appear so we trim them
			}
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
	if command.LineStart != headerLine && command.LineStart >= 1 {
		startLine = command.LineStart
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
		if (i != headerLine && i >= startLine) || len(lines) == 1 {
			if i >= command.LineEnd && command.LineEnd != 0 {
				logger.Flex("debug", nil, fmt.Sprintf("reached line limit %d", command.LineEnd), false)
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

			if len(cmdSample) > 0 {
				applyCustomAttributes(&cmdSample, &command.CustomAttributes)
				// load.StoreAppend(cmdSample)
				*dataStore = append(*dataStore, cmdSample)
			}
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
		return load.TypeJSON, f
	}

	// default raw
	return "raw", nil
}
