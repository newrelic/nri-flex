/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package discovery

import (
	"fmt"
	"strings"

	"github.com/newrelic/nri-flex/internal/load"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
	"github.com/sirupsen/logrus"
)

// ProcessNetworkStat x
type ProcessNetworkStat struct {
	Name string
	Data string
}

// Processes loops through tcp connections and returns the corresponding process and connection information
func Processes() {
	conns, err := net.Connections("tcp")
	if err != nil {
		load.Logrus.WithFields(logrus.Fields{
			"err": err,
		}).Error("discovery: processes unable to get tcp connections")
	} else {
		load.DiscoveredProcesses = map[string]string{}
		for _, conn := range conns {
			p, err := process.NewProcess(conn.Pid)
			if err == nil {
				running, _ := p.IsRunning()
				if running {
					name, err := p.Name()
					if err != nil {
						load.Logrus.WithFields(logrus.Fields{
							"err": err,
						}).Error("discovery: processes unable to get name")
					}
					cmd, err := p.Cmdline()
					if err != nil {
						load.Logrus.WithFields(logrus.Fields{
							"err": err,
						}).Error("discovery: processes unable to cmd line")
					}

					if checkBlacklistedProcess(name, cmd) {
						continue
					}

					load.DiscoveredProcesses[fmt.Sprintf("%v", conn.Pid)] =
						fmt.Sprintf(`{"name":"%v","cmd":"%v","localIP":"%v","localPort":"%v","remoteIP":"%v","remotePort":"%v"}`, name, cmd, conn.Laddr.IP, conn.Laddr.Port, conn.Raddr.IP, conn.Raddr.Port)
				}
			}
		}
	}
}

func checkBlacklistedProcess(name string, cmd string) bool {
	blacklistedProcesses := []string{
		"Chrome", "Visual Studio Code", "BlueJeans", "WhatsApp", "Insomnia", "Slack", "SpotifyWebHelper", "ZoomOpener",
		"Dashlane", "docker.for.mac", "svchost", "lsass", "wininit", "spoolsv", "[System Process]"}
	for _, blProcess := range blacklistedProcesses {
		if caseInsensitiveContains(name, blProcess) || caseInsensitiveContains(cmd, blProcess) {
			return true
		}
	}
	return false
}

func caseInsensitiveContains(s, substr string) bool {
	s, substr = strings.ToUpper(s), strings.ToUpper(substr)
	return strings.Contains(s, substr)
}
