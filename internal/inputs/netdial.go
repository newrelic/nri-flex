/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package inputs

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"net/textproto"
	"time"

	"github.com/newrelic/nri-flex/internal/load"
)

// NetDialWithTimeout performs network dial without timeout
func NetDialWithTimeout(dataStore *[]interface{}, command load.Command, dataSample *map[string]interface{}, api load.API, processType *string) {

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
	load.Logrus.Debugf("commands: dialling %v : %v", addr, netw)
	dialConn, err := net.DialTimeout(netw, addr, time.Duration(timeout)*time.Millisecond)
	if err == nil {
		defer dialConn.Close()
	}
	go func(dialConn net.Conn, err error) {
		if err != nil {
			dialError = err
		} else {
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
	}(dialConn, err)

	// Listen for signals
	select {
	case <-ctx.Done():
		if command.Run == "" {
			*dataStore = append(*dataStore, map[string]interface{}{"portStatus": "closed", "addr": command.Dial, "netw": netw, "err": ctx.Err().Error()})
			// load.StoreAppend(map[string]interface{}{"portStatus": "closed", "addr": command.Dial, "netw": netw, "err": ctx.Err().Error()})
		} else if command.Run != "" && data != "" {
			processOutput(dataStore, data, dataSample, command, api, processType)
		}
		if data == "" {
			load.Logrus.Error("commands: dial " + ctx.Err().Error())
		} else {
			load.Logrus.Debug("commands: dial " + ctx.Err().Error())
		}
	case <-c:
		if command.Run == "" && dialError == nil {
			*dataStore = append(*dataStore, map[string]interface{}{"portStatus": "open", "addr": command.Dial, "netw": netw})
			// load.StoreAppend(map[string]interface{}{"portStatus": "open", "addr": command.Dial, "netw": netw})
		} else if command.Run == "" && dialError != nil {
			*dataStore = append(*dataStore, map[string]interface{}{"portStatus": "closed", "addr": command.Dial, "netw": netw, "err": dialError.Error()})
			// load.StoreAppend(map[string]interface{}{"portStatus": "closed", "addr": command.Dial, "netw": netw, "err": dialError.Error()})
		} else if command.Run != "" && dialError == nil && data != "" {
			processOutput(dataStore, data, dataSample, command, api, processType)
		}
		load.Logrus.Debugf("commands: finished dial %v : %v", command.Dial, netw)
	}
}
