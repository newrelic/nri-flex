/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package runtime

import (
	"fmt"
	"github.com/newrelic/nri-flex/internal/config"
	"github.com/newrelic/nri-flex/internal/discovery"
	"github.com/newrelic/nri-flex/internal/load"
	"runtime"
)

type Default struct {
}

// The Default, server-based, runtime is always available
func (i *Default) isAvailable() bool {
	return true
}

// Run Flex on a server
func (i *Default) loadConfigs(configs *[]load.Config) error {
	if load.Args.EncryptPass != "" {
		err := logEncryptPass()
		if err != nil {
			return err
		}
	}

	_, err := config.SyncGitConfigs("")
	if err != nil {
		log.WithError(err).Warn("Default.loadConfigs: failed to sync git configs")
	}

	var errors []error
	if load.Args.ConfigFile != "" {
		err = addSingleConfigFile(load.Args.ConfigFile, configs)
		if err != nil {
			return err
		}
	} else {
		errors = addConfigsFromPath(load.Args.ConfigDir, configs)
		if len(errors) > 0 {
			return fmt.Errorf(" Default.loadConfigs: failed to load configurations files")
		}
	}

	// should we stop if this fails??
	if load.Args.ContainerDiscovery || load.Args.Fargate {
		discovery.Run(configs)
	}
	if load.ContainerID == "" {
		switch runtime.GOOS {
		case "windows":
			if load.Args.DiscoverProcessWin {
				discovery.Processes()
			}
		case "linux":
			if load.Args.DiscoverProcessLinux {
				discovery.Processes()
			}
		}
	}

	return nil
}

func (i *Default) SetConfigDir(s string) {
}
