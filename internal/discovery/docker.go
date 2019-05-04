package discovery

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/newrelic/nri-flex/internal/load"
	"github.com/newrelic/nri-flex/internal/logger"

	"github.com/docker/docker/api"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

var cli *client.Client

// setDockerClient Sets the docker client
// There can be edge cases when the integration API version may need a matching or lower API version then the hosts docker API version
func setDockerClient() (*client.Client, error) {
	var out []byte
	var err error

	if load.Args.DockerAPIVersion != "" {
		cli, err = client.NewClientWithOpts(client.WithVersion(load.Args.DockerAPIVersion))
	} else {
		logger.Flex("info", nil, fmt.Sprintf("GOOS: %v", runtime.GOOS), false)

		if runtime.GOOS == "windows" {
			out, err = exec.Command("cmd", "/C", `docker`, `version`, `--format`, `"{{json .Client.APIVersion}}"`).Output()
		} else {
			out, err = exec.Command(`docker`, `version`, `--format`, `"{{json .Client.APIVersion}}"`).Output()
			if err != nil {
				out, err = exec.Command(`/host/usr/local/bin/docker`, `version`, `--format`, `"{{json .Client.APIVersion}}"`).Output()
			}
		}

		if err != nil {
			logger.Flex("error", err, "unable to fetch Docker API version - setting client with NewClientWithOpts()", false)
			cli, err = client.NewClientWithOpts()
		} else {
			cmdOut := string(out)
			clientAPIVersion := strings.TrimSpace(strings.Replace(cmdOut, `"`, "", -1))
			clientVer, _ := strconv.ParseFloat(clientAPIVersion, 64)
			apiVer, _ := strconv.ParseFloat(api.DefaultVersion, 64)

			if clientVer <= apiVer {
				logger.Flex("info", nil, fmt.Sprintf("Setting client with version:%v", clientAPIVersion), false)
				cli, err = client.NewClientWithOpts(client.WithVersion(clientAPIVersion))
			} else {
				logger.Flex("info", nil, fmt.Sprintf("Client API Version %v is higher then integration version %v", clientAPIVersion, api.DefaultVersion), false)
				logger.Flex("info", nil, "Setting client with NewClientWithOpts()", false)
				cli, err = client.NewClientWithOpts()
			}
		}
	}

	return cli, err
}

func execContainerCommand(containerID string, command []string) (string, error) {
	if cli == nil {
		var err error
		cli, err = setDockerClient()
		if err != nil {
			return "", err
		}
	}

	ctx := context.Background()
	execConfig := types.ExecConfig{
		AttachStderr: true,
		AttachStdin:  true,
		AttachStdout: true,
		Cmd:          command,
		Tty:          true,
		Detach:       false,
	}
	//set target container
	exec, err := cli.ContainerExecCreate(ctx, containerID, execConfig)
	if err != nil {
		return "", err
	}
	execAttachConfig := types.ExecStartCheck{
		Detach: false,
		Tty:    true,
	}
	containerConn, err := cli.ContainerExecAttach(ctx, exec.ID, execAttachConfig)
	if err != nil {
		return "", err
	}
	defer containerConn.Close()
	data, err := Readln(containerConn.Reader)
	if err != nil {
		return "", err
	}
	return data, nil
}

// Readln from bufioReader
func Readln(r *bufio.Reader) (string, error) {
	var (
		isPrefix = true
		err      error
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return string(ln), err
}
