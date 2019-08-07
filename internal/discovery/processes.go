package discovery

import (
	"fmt"
	"strings"

	"github.com/newrelic/nri-flex/internal/load"
	"github.com/newrelic/nri-flex/internal/logger"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
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
		logger.Flex("error", err, "unable to get tcp connections", false)
	} else {
		load.DiscoveredProcesses = map[string]string{}
		for _, conn := range conns {
			p, err := process.NewProcess(conn.Pid)
			if err == nil {
				running, _ := p.IsRunning()
				if running {
					name, err := p.Name()
					logger.Flex("error", err, "", false)
					cmd, err := p.Cmdline()
					logger.Flex("error", err, "", false)

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
