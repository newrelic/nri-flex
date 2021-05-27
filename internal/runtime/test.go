package runtime

import (
	"fmt"
	"github.com/newrelic/nri-flex/internal/config"
	"github.com/newrelic/nri-flex/internal/discovery"
	"github.com/newrelic/nri-flex/internal/load"
)

type Test struct {
}

// The Test runtime is always available
func (i *Test) isAvailable() bool {
	return true
}

// Run the Test  runtime
func (i *Test) loadConfigs(configs *[]load.Config) error {
	if load.Args.EncryptPass != "" {
		err := logEncryptPass()
		if err != nil {
			return err
		}
	}

	_, err := config.SyncGitConfigs("")
	if err != nil {
		log.WithError(err).Warn("Test.loadConfigs: failed to sync git configs")
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
			return fmt.Errorf(" Test.loadConfigs: failed to load configurations files")
		}
	}

	// should we stop if this fails??
	if load.Args.ContainerDiscovery || load.Args.Fargate {
		discovery.Run(configs)
	}
	return nil
}

func (i *Test) SetConfigDir(s string) {
}
