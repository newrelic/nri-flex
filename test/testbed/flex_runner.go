package testbed

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"
)

var (
	errNoFlexRunning = fmt.Errorf("No nri-flex running")
	errNoResults     = fmt.Errorf("No nri-flex results available, the process might have failed")
)

type FlexRunner interface {
	Run() error
	Stop() error
	// returns the results of the run process
	Results() (string, string, error)
}

// childInfraRunner implements the FlexRunner interface as a child process on the same machine
type childFlexRunner struct {
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
	stopOnce   sync.Once
	isStopped  bool
	isFinished bool
}

func NewChildFlexRunner(binPath, configPath string) *childFlexRunner {
	return &childFlexRunner{
		BinPath:    binPath,
		ConfigPath: configPath,
	}
}

func (cr *childFlexRunner) Run() error {
	//log.Printf("Starting nri-flex (%s) (%s)", cr.BinPath, cr.ConfigPath)

	// prepare command
	args := []string{"-config_path", cr.ConfigPath}
	cr.cmd = exec.Command(cr.BinPath, args...)
	//log.Printf("Command (%s)", cr.cmd.String())

	// Capture standard output and standard error.
	cr.cmd.Stdout = &cr.LogStdout
	cr.cmd.Stderr = &cr.LogStderr

	cr.isStarted = true

	// Run the process.
	err := cr.cmd.Run()
	if err != nil {
		log.Printf("Error running nri-flex (%s)", err)
		return err
	}

	cr.isFinished = true
	return nil
}

func (cr *childFlexRunner) Stop() (err error) {
	if !cr.isStarted || cr.isStopped {
		return errNoFlexRunning
	}
	cr.stopOnce.Do(func() {
		cr.isStopped = true

		log.Printf("Gracefully terminating nri-flex pid=%d, sending SIGTEM...", cr.cmd.Process.Pid)

		// Gracefully signal process to stop.
		if err := cr.cmd.Process.Signal(syscall.SIGTERM); err != nil {
			log.Printf("Cannot send SIGTEM: %s", err.Error())
		}

		finished := make(chan struct{})

		// Setup a goroutine to wait a while for process to finish and send kill signal
		// to the process if it doesn't finish.
		go func() {
			// Wait 15 seconds.
			t := time.After(15 * time.Second)
			select {
			case <-t:
				log.Printf("nri-flex pid=%d is not responding to SIGTERM. Sending SIGKILL to kill forcedly.",
					cr.cmd.Process.Pid)
				if err = cr.cmd.Process.Signal(syscall.SIGKILL); err != nil {
					log.Printf("Cannot send SIGKILL: %s", err.Error())
				}
			case <-finished:
			}
		}()

		// Wait for process to terminate
		err = cr.cmd.Wait()

	})
	return
}

// return stdout and stderr files
func (cr *childFlexRunner) Results() (string, string, error) {
	if !cr.isFinished {
		return "", "", errNoResults
	}
	return cr.LogStdout.String(), cr.LogStderr.String(), nil
}
