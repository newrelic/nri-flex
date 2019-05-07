package inputs

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net"
	"net/textproto"
	"time"

	"github.com/newrelic/nri-flex/internal/load"
	"github.com/newrelic/nri-flex/internal/logger"
)

// NetDialWithTimeout performs network dial without timeout
func NetDialWithTimeout(command load.Command, dataSample *map[string]interface{}, api load.API, processType *string) {
	ctx := context.Background()
	// Create a channel for signal handling
	c := make(chan struct{})
	// Define a cancellation after default dial timeout in the context
	timeout := load.DefaultDialTimeout // ms
	if api.Timeout > 0 {
		timeout = api.Timeout
	}
	if command.Timeout > 0 {
		timeout = command.Timeout
	}

	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Millisecond)
	defer cancel()

	addr := command.Dial
	netw := "tcp"
	if command.Network != "" {
		netw = command.Network
	}

	var dialError error
	var data string
	// Run dial via a goroutine
	go func() {
		logger.Flex("debug", nil, fmt.Sprintf("dialling %v : %v", addr, netw), false)
		dialConn, err := net.DialTimeout(netw, addr, time.Duration(timeout)*time.Millisecond)
		if err != nil {
			dialError = err
		} else {
			defer dialConn.Close()
			if command.Run != "" {
				fmt.Fprintf(dialConn, command.Run)
				reader := bufio.NewReader(dialConn)
				tp := textproto.NewReader(reader)
				for {
					line, _ := tp.ReadLine()
					data += line + "\n"
				}
			}
		}
		c <- struct{}{}
	}()

	// Listen for signals
	select {
	case <-ctx.Done():
		if command.Run == "" {
			load.StoreAppend(map[string]interface{}{"portStatus": "closed", "addr": command.Dial, "netw": netw, "err": ctx.Err().Error()})
		} else if command.Run != "" && data != "" {
			processOutput(data, dataSample, command, api, processType)
		}
		if data == "" {
			logger.Flex("error", errors.New("dial: "+ctx.Err().Error()), "", false)
		} else {
			logger.Flex("debug", errors.New("dial: "+ctx.Err().Error()), "", false)
		}
	case <-c:
		if command.Run == "" && dialError == nil {
			load.StoreAppend(map[string]interface{}{"portStatus": "open", "addr": command.Dial, "netw": netw})
		} else if command.Run == "" && dialError != nil {
			load.StoreAppend(map[string]interface{}{"portStatus": "closed", "addr": command.Dial, "netw": netw, "err": dialError.Error()})
		} else if command.Run != "" && dialError == nil && data != "" {
			processOutput(data, dataSample, command, api, processType)
		}
		logger.Flex("debug", nil, "dial finished", false)
	}
}
