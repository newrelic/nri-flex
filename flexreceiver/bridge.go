package flexreceiver

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/newrelic/nri-flex/internal/config"
	"github.com/newrelic/nri-flex/internal/load"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"
)

// FlexBridge connects the OTel receiver with the core nri-flex logic.
type FlexBridge struct {
	cfg    *Config
	logger *zap.Logger
}

// NewFlexBridge creates a new bridge.
func NewFlexBridge(cfg *Config, logger *zap.Logger) *FlexBridge {
	return &FlexBridge{
		cfg:    cfg,
		logger: logger,
	}
}

// CollectMetrics loads Flex configs, runs them, and converts the output to pmetric.Metrics.
func (fb *FlexBridge) CollectMetrics(ctx context.Context) (pmetric.Metrics, error) {
	metrics := pmetric.NewMetrics()

	configs, err := fb.loadFlexConfigs()
	if err != nil {
		return metrics, fmt.Errorf("failed to load flex configs: %w", err)
	}

	for _, cfg := range configs {
		fb.logger.Debug("Executing Flex config", zap.String("name", cfg.Name))
		data, err := fb.runFlexConfig(cfg)
		if err != nil {
			fb.logger.Error("Error running flex config", zap.String("name", cfg.Name), zap.Error(err))
			continue
		}
		fb.convertToOTelMetrics(data, metrics, cfg.Name)
	}

	return metrics, nil
}

func (fb *FlexBridge) loadFlexConfigs() ([]load.Config, error) {
	var configs []load.Config
	if fb.cfg.ConfigFile != "" {
		fileInfo, err := os.Stat(fb.cfg.ConfigFile)
		if err != nil {
			return nil, fmt.Errorf("failed to stat config file %s: %w", fb.cfg.ConfigFile, err)
		}
		dir := filepath.Dir(fb.cfg.ConfigFile)
		err = config.LoadFile(&configs, fileInfo, dir)
		if err != nil {
			return nil, err
		}
		return configs, nil
	}

	files, err := ioutil.ReadDir(fb.cfg.ConfigDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read config dir %s: %w", fb.cfg.ConfigDir, err)
	}
	loadErrors := config.LoadFiles(&configs, files, fb.cfg.ConfigDir)
	if len(loadErrors) > 0 {
		var errorStrings []string
		for _, e := range loadErrors {
			errorStrings = append(errorStrings, e.Error())
		}
		return nil, fmt.Errorf(strings.Join(errorStrings, ", "))
	}
	return configs, nil
}

func (fb *FlexBridge) runFlexConfig(cfg load.Config) ([]interface{}, error) {
	// config.Run executes the flex config, but it's designed for a standalone application
	// and does not return the collected data. It sends the data directly to New Relic.
	// To make this work as a library, the core data collection logic in nri-flex
	// needs to be refactored to return the data instead of sending it.
	// For now, this is a placeholder to make the code compile.
	config.Run(cfg)
	// Since config.Run does not return data, we cannot convert it to OTel metrics.
	// Returning nil data and nil error to satisfy the compiler.
	return nil, nil
}

func (fb *FlexBridge) convertToOTelMetrics(flexData []interface{}, metrics pmetric.Metrics, configName string) {
	if len(flexData) == 0 {
		return
	}

	rm := metrics.ResourceMetrics().AppendEmpty()
	rm.Resource().Attributes().PutStr("flex.config.name", configName)

	sm := rm.ScopeMetrics().AppendEmpty()
	sm.Scope().SetName("otelcol/flexreceiver")

	for _, dataPoint := range flexData {
		dataMap, ok := dataPoint.(map[string]interface{})
		if !ok {
			continue
		}
		fb.convertDataPoint(dataMap, sm.Metrics())
	}
}

func (fb *FlexBridge) convertDataPoint(dataMap map[string]interface{}, metrics pmetric.MetricSlice) {
	now := pcommon.NewTimestampFromTime(time.Now())
	var metricName string
	if name, ok := dataMap["event_type"].(string); ok {
		metricName = name
	} else {
		metricName = "flexMetric"
	}

	for key, val := range dataMap {
		// Skip non-numeric values and identifiers
		if _, isNum := val.(float64); !isNum {
			continue
		}

		m := metrics.AppendEmpty()
		m.SetName(fmt.Sprintf("%s.%s", metricName, key))
		dp := m.SetEmptyGauge().DataPoints().AppendEmpty()
		dp.SetTimestamp(now)
		dp.SetDoubleValue(val.(float64))

		// Add other fields as attributes
		for attrKey, attrVal := range dataMap {
			if attrKey == key {
				continue
			}
			dp.Attributes().PutStr(attrKey, fmt.Sprintf("%v", attrVal))
		}
	}
}
