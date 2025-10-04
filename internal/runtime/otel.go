package runtime

import (
	"github.com/newrelic/nri-flex/internal/load"
)

type OTel struct {
}

func (o *OTel) isAvailable() bool {
	return true
}

func (o *OTel) loadConfigs(configs *[]load.Config) error {
	return nil
}

func (o *OTel) SetConfigDir(dir string) {
}
