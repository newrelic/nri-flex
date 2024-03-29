/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package inputs

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	xj "github.com/basgys/goxml2json"
	"github.com/newrelic/nri-flex/internal/formatter"
	"github.com/newrelic/nri-flex/internal/load"
	"github.com/sirupsen/logrus"
)

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func checkOS(os string) bool {
	if os == "" || runtime.GOOS == os {
		return true
	}
	return false
}

// RunCommands executes the given commands to create one merged sampled
func RunCommands(dataStore *[]interface{}, yml *load.Config, apiNo int) {
	startTime := makeTimestamp()
	api := yml.APIs[apiNo]

	load.Logrus.WithFields(logrus.Fields{
		"name":  yml.Name,
		"count": len(api.Commands),
	}).Debug("commands: executing")

	dataSample := map[string]interface{}{}
	processType := ""
	for _, command := range api.Commands {
		if command.Run != "" && command.Dial == "" && checkOS(command.OS) {
			commandRun(dataStore, yml, command, api, startTime, dataSample, processType)
		} else if command.Cache != "" {
			if yml.Datastore[command.Cache] != nil {
				for _, cache := range yml.Datastore[command.Cache] {
					switch sample := cache.(type) {
					case map[string]interface{}:
						if sample["http"] != nil {
							load.Logrus.WithFields(logrus.Fields{
								"cache": command.Cache,
							}).Debug("command: processing http cache with command processor")

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
					load.Logrus.Debug("command: not handled yet")
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

func commandRun(dataStore *[]interface{}, yml *load.Config, command load.Command, api load.API, startTime int64, dataSample map[string]interface{}, processType string) {
	command.Run = envCommandCheck(command.Run)
	runCommand := command.Run
	if command.Output == load.Jmx {
		SetJMXCommand(&runCommand, command, api, yml)
		command.Run = runCommand
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

	// Create the command with our context
	cmd := buildCommand(ctx, api, command)

	// https://golang.org/pkg/os/exec/#Cmd.StdinPipe
	if load.Args.StdinPipe {
		_, err := cmd.StdinPipe()
		if err != nil {
			load.Logrus.WithFields(logrus.Fields{
				"exec":       command.Run,
				"err":        err,
				"suggestion": "StdinPipe failed",
			}).Debug("command: failed")
		}
	}

	output, err := cmd.CombinedOutput()

	// check if a assertion is defined and successfully passes before continuing, see function for detailed comments
	if !checkAssertion(command.Assert, output) {
		load.Logrus.WithFields(logrus.Fields{
			"name": yml.Name,
			"exe":  command.Run,
		}).Debug("commands: assertion failed will not process")
		return
	}

	contextError := ctx.Err()

	if err != nil || contextError != nil {
		contextErrorStr := ""
		if contextError != nil {
			contextErrorStr = contextError.Error()
		}

		load.Logrus.WithFields(logrus.Fields{
			"exec":        command.Run,
			"err":         err,
			"context_err": contextErrorStr,
			"suggestion":  "if you are handling this error case, ignore",
		}).Debug("command: failed")

		if command.HideErrorExec {
			errorSample := map[string]interface{}{
				"error":         err,
				"error_msg":     string(output),
				"context_error": contextErrorStr,
				"error_exec":    "COMMAND HIDDEN!",
			}
			*dataStore = append(*dataStore, errorSample)
			return
		}

		errorSample := map[string]interface{}{
			"error":         err,
			"error_msg":     string(output),
			"context_error": contextErrorStr,
			"error_exec":    command.Run,
		}
		*dataStore = append(*dataStore, errorSample)
		return
	}

	if len(output) > 0 {
		if command.SplitOutput != "" {
			splitOutput(dataStore, string(output), command, startTime)
		} else {
			processOutput(dataStore, string(output), &dataSample, command, api, &processType)
		}
	}
}

// checks if explicitedly enabled log
func envCommandCheck(commandStr string) string {
	if load.Args.AllowEnvCommands {
		load.Logrus.WithFields(logrus.Fields{
			"command": commandStr,
			"prepend": os.Getenv("FLEX_CMD_PREPEND"),
			"append":  os.Getenv("FLEX_CMD_APPEND"),
		}).Info("command: environment modification enabled")

		switch os.Getenv("FLEX_CMD_WRAP") {
		case "\"":
			commandStr = fmt.Sprintf("\"%v\"", commandStr)
		case "'":
			commandStr = fmt.Sprintf("'%v'", commandStr)
		case "`":
			commandStr = fmt.Sprintf("`%v`", commandStr)
		case "true":
			commandStr = fmt.Sprintf("\"%v\"", commandStr)
		}

		return os.Getenv("FLEX_CMD_PREPEND") + commandStr + os.Getenv("FLEX_CMD_APPEND")
	}
	return commandStr
}

func splitOutput(dataStore *[]interface{}, output string, command load.Command, startTime int64) {
	lines := strings.Split(strings.TrimSuffix(output, "\n"), "\n")
	var outputBlocks [][]string
	startSplit := -1 // initialize
	endSplit := 0

	if len(lines) == 1 {
		outputBlocks = append(outputBlocks, lines[0:1])
	} else {
		for i, line := range lines {
			if formatter.KvFinder("regex", line, command.SplitOutput) {
				if startSplit == -1 {
					startSplit = 0
					endSplit = i
					outputBlocks = append(outputBlocks, lines[startSplit:endSplit])
					startSplit = i
				} else {
					endSplit = i
					outputBlocks = append(outputBlocks, lines[startSplit:endSplit])
					startSplit = i
				}
			}

			//create the last block
			if i+1 == len(lines) && startSplit != -1 {
				outputBlocks = append(outputBlocks, lines[startSplit:i+1])
			}
		}
	}

	processBlocks(dataStore, outputBlocks, command, startTime)
}

func processBlocks(dataStore *[]interface{}, blocks [][]string, command load.Command, startTime int64) {
	for _, block := range blocks {
		sample := map[string]interface{}{}

		if len(command.RegexMatches) > 0 {
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

		} else {
			processRaw(&sample, "", block, command)
		}

		// do not add empty samples
		if len(sample) > 0 {
			sample["flex.commandTimeMs"] = makeTimestamp() - startTime
			*dataStore = append(*dataStore, sample)
		}

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

			load.Logrus.Debugf("command: running %v", cmd)

			if command.Split == "" { // default vertical split
				applyCustomAttributes(dataSample, &command.CustomAttributes)
				processRaw(dataSample, dataOutput, []string{}, command)
			} else if command.Split == load.TypeColumns || command.Split == "horizontal" {
				if *processType == load.TypeColumns {
					load.Logrus.Debugf("command: horizontal split only allowed once per command set %v %v", api.Name, command.Name)
				} else {
					*processType = "columns"
					processRawCol(dataStore, dataSample, dataOutput, command)
				}
			}
		case load.TypeJSON:
			// load.StoreAppend(dataInterface)
			*dataStore = append(*dataStore, dataInterface)
		case load.TypeXML:
			*dataStore = append(*dataStore, dataInterface)
		case load.Jmx:
			*processType = "jmx"
			ParseJMX(dataStore, dataInterface, command, dataSample)
		case load.TypeCSV:
			err := processCsv(dataStore, "", "command output", &dataOutput, command.SetHeader)
			if err != nil {
				load.Logrus.WithError(err).Errorf("Failed to process text/csv body")
			}

		}
	}
}

// processRaw processes a raw data output
func processRaw(dataSample *map[string]interface{}, dataOutput string, lines []string, command load.Command) {
	splitBy := command.SplitBy
	lineStart := command.LineStart
	lineEnd := command.LineEnd

	// if no lines exist, check the dataOutput to split lines from
	if len(lines) == 0 && dataOutput != "" {
		lines = strings.Split(strings.TrimSuffix(dataOutput, "\n"), "\n")
	}

	// SplitBy key is required else we cannot easily distinguish between keys and values
	for i, line := range lines {
		if i >= lineStart {
			if i >= lineEnd && lineEnd != 0 {
				load.Logrus.Debugf("command: reached line limit %d", lineEnd)
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
	startLine := 0

	if command.RowHeader != 0 {
		headerLine = command.RowHeader
	}
	// this is buggy? If we set RowHeader to 10, RowStart should be after RowHeader, no?
	// ie, RowStart > RowHeader, always?
	if command.RowStart != headerLine && command.RowStart >= 1 {
		startLine = command.RowStart
	}

	// same comment as above. but why not just ignore LineStart and use only RowStart?
	// this is a bit confusing to have different options for the same thing.
	if command.LineStart != headerLine && command.LineStart >= 1 {
		startLine = command.LineStart
	}

	lines := strings.Split(strings.TrimSuffix(dataOutput, "\n"), "\n")
	header := lines[headerLine]
	var keys []string

	// set header keys
	if len(command.SetHeader) > 0 {
		keys = command.SetHeader
		headerLine = -1
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
				load.Logrus.Debugf("command: reached line limit %d", command.LineEnd)
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
	if commandOutput == load.TypeCSV {
		return "csv", nil
	}
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

			load.Logrus.WithFields(logrus.Fields{
				"err": err,
			}).Error("commands: failed to unmarshal jmx output")
		}
		return load.Jmx, map[string]interface{}{"error": "Failed to process JMX Data"}
	}

	// check json
	var f interface{}
	err := json.Unmarshal([]byte(dataOutput), &f)
	if err == nil {
		return load.TypeJSON, f
	}
	// check xml
	xmlSignature := `<?xml version=`
	if strings.HasPrefix(strings.TrimSpace(dataOutput), xmlSignature) || commandOutput == load.TypeXML {
		// return load.TypeXML, nil
		xmlBody := strings.NewReader(dataOutput)
		jsonBody, err := xj.Convert(xmlBody)

		if err != nil {
			load.Logrus.WithFields(logrus.Fields{
				"err": err,
			}).Errorf("Failed to convert XML to Json ")
		} else {
			var f interface{}
			err := json.Unmarshal(jsonBody.Bytes(), &f)
			if err == nil {
				return load.TypeXML, f
			}
		}

	}

	// default raw
	return "raw", nil
}

// in *nix, the command will be run with "/bin/sh"
// in windows it will be run under "cmd". some windows set powershell as the default
// to override the defaults, set the shell to run either at the API level or command level.
// for *unix append the "-c", for windows "/c" unless we override the shell. in that case flags should be provided
func buildCommand(ctx context.Context, api load.API, command load.Command) *exec.Cmd {
	commandShell := load.DefaultShell
	// not sure we should keep this for other shells
	secondParameter := "-c"

	if runtime.GOOS == "windows" {
		commandShell = "cmd"
		secondParameter = "/C"
	}

	if api.Shell != "" {
		commandShell = api.Shell
	}
	if command.Shell != "" {
		commandShell = command.Shell
	}

	return exec.CommandContext(ctx, commandShell, secondParameter, command.Run)
}

// checkAssertion perform output based assertions
// when a match or not_match value is defined in the command section an assertion will be performed
// if only match is defined, and the output successfully matches it will continue
// if only not_match is defined, and the output successfully does not contain the output it will continue
// if both match and not_match is defined, both the above must be true to continue
func checkAssertion(assert load.Assert, output []byte) bool {

	// return true if no matches defined
	if assert.Match == "" && assert.NotMatch == "" {
		return true
	}

	strOutput := string(output)

	if assert.Match != "" && assert.NotMatch != "" {
		return formatter.KvFinder("regex", strOutput, assert.Match) && !formatter.KvFinder("regex", strOutput, assert.NotMatch)
	} else if assert.Match != "" {
		return formatter.KvFinder("regex", strOutput, assert.Match)
	} else if assert.NotMatch != "" {
		return !formatter.KvFinder("regex", strOutput, assert.NotMatch)
	}

	return false
}
