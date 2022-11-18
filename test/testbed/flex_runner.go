package testbed

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

var (
	errNoResults = fmt.Errorf("no nri-flex results available, the process might have failed")
)

type FlexRunner interface {
	Run() error
	// returns the results of the run process
	Results() (string, string, error)
}

// childInfraRunner implements the FlexRunner interface as a child process on the same machine
type ChildFlexRunner struct {
	// Path to nri-flex executable
	BinPath string

	// Configuration file path
	ConfigPath string

	// captures stdout
	LogStdout strings.Builder
	// captures sterr
	LogStderr strings.Builder

	// Command to execute
	cmd *exec.Cmd

	isStarted  bool
	isFinished bool
}

func NewChildFlexRunner(binPath, configPath string) *ChildFlexRunner {
	return &ChildFlexRunner{
		BinPath:    binPath,
		ConfigPath: configPath,
	}
}

func (cr *ChildFlexRunner) Run() error {

	args := []string{"-config_path", cr.ConfigPath}
	// #nosec G204
	cr.cmd = exec.Command(cr.BinPath, args...)

	// Capture standard output and standard error.
	cr.cmd.Stdout = &cr.LogStdout
	cr.cmd.Stderr = &cr.LogStderr

	cr.isStarted = true

	// Run the process.
	err := cr.cmd.Run()
	if err != nil {
		log.Printf("Error running nri-flex (%s)", err)
		log.Printf("Stderr: (%s)", cr.LogStderr.String())
		return err
	}

	cr.isFinished = true
	return nil
}

// return stdout and stderr files
func (cr *ChildFlexRunner) Results() (string, string, error) {
	if !cr.isFinished {
		return "", "", errNoResults
	}
	return cr.LogStdout.String(), cr.LogStderr.String(), nil
}
