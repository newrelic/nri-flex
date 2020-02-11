package gofile

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

// Run executes the go source file (which must contain a main() function) with the
// provided arguments.
// At returns the standard output of the command and logs (debug mode) the standard error
func Run(filePath string, args ...string) ([]byte, error) {
	gocmd, err := exec.LookPath("go")
	if err != nil {
		return nil, err
	}

	args = append([]string{"run", filePath}, args...)
	cmd := exec.Command(gocmd, args...)

	sp, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	so, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}
	stderr, err := ioutil.ReadAll(sp)
	if err != nil {
		return nil, err
	}
	fmt.Fprint(os.Stderr, string(stderr))

	stdout, err := ioutil.ReadAll(so)
	if err != nil {
		return nil, err
	}
	if err := cmd.Wait(); err != nil {
		return stdout, err
	}

	return stdout, nil
}
