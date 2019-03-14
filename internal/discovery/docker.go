package discovery

import (
	"bufio"
	"context"

	"github.com/docker/docker/api/types"
)

func execContainerCommand(containerID string, command []string) (string, error) {
	cli, err := setDockerClient()
	if err != nil {
		return "", err
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
