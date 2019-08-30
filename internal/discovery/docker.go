/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package discovery

import (
	"bufio"
	"context"
	"fmt"

	"github.com/newrelic/nri-flex/internal/load"
	"github.com/newrelic/nri-flex/internal/logger"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

var cli *client.Client

// setDockerClient Sets the docker client
// There can be edge cases when the integration API version may need a matching or lower API version then the hosts docker API version
func setDockerClient() (*client.Client, error) {
	// var out []byte
	var err error
	if load.Args.DockerAPIVersion != "" {
		logger.Flex("debug", nil, fmt.Sprintf("setting docker client via argument %v", load.Args.DockerAPIVersion), false)
		cli, err = client.NewClientWithOpts(client.WithVersion(load.Args.DockerAPIVersion))
	} else {
		logger.Flex("debug", err, "setting docker client with API version negotiation", false)
		cli, err = client.NewClientWithOpts(client.WithAPIVersionNegotiation())
	}
	return cli, err
}

// ExecContainerCommand execute command against a container
func ExecContainerCommand(containerID string, command []string) (string, error) {
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
