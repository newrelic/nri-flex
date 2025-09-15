package simplereceiver

import (
	"context"
	"time"

	"github.com/newrelic/nri-flex/internal/load"
	"github.com/newrelic/nri-flex/internal/runtime"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/receiver"
	"go.uber.org/zap"
)

// Configurable settings for the receiver
type Config struct {
	CollectionInterval time.Duration `mapstructure:"collection_interval"`
	ConfigDir          string        `mapstructure:"config_dir"`
}

func createDefaultConfig() component.Config {
	return &Config{
		CollectionInterval: 30 * time.Second,
		ConfigDir:          "./flexConfigs",
	}
}

const (
	typeStr   = "simpleflex"
	stability = component.StabilityLevelBeta
)

// NewFactory creates a factory for the receiver.
func NewFactory() receiver.Factory {
	return receiver.NewFactory(
		component.MustNewType(typeStr),
		createDefaultConfig,
		receiver.WithMetrics(createMetricsReceiver, stability),
	)
}

type simpleReceiver struct {
	config   *Config
	settings receiver.Settings
	consumer consumer.Metrics
	cancel   context.CancelFunc
	logger   *zap.Logger
	runtime  *runtime.OTel
}

func createMetricsReceiver(
	_ context.Context,
	settings receiver.Settings,
	cfg component.Config,
	consumer consumer.Metrics,
) (receiver.Metrics, error) {
	config := cfg.(*Config)

	// Initialize our runtime instance for nri-flex
	runtimeInstance := &runtime.OTel{}

	return &simpleReceiver{
		config:   config,
		settings: settings,
		consumer: consumer,
		logger:   settings.Logger,
		runtime:  runtimeInstance,
	}, nil
}

func (r *simpleReceiver) Start(ctx context.Context, _ component.Host) error {
	r.logger.Info("Starting simple flex receiver")
	var runCtx context.Context
	runCtx, r.cancel = context.WithCancel(ctx)

	// Set the config directory for nri-flex
	r.runtime.SetConfigDir(r.config.ConfigDir)

	go r.runCollectionLoop(runCtx)
	return nil
}

func (r *simpleReceiver) Shutdown(ctx context.Context) error {
	r.logger.Info("Shutting down simple flex receiver")
	if r.cancel != nil {
		r.cancel()
	}
	return nil
}

func (r *simpleReceiver) runCollectionLoop(ctx context.Context) {
	ticker := time.NewTicker(r.config.CollectionInterval)
	defer ticker.Stop()

	r.collectAndSend(ctx) // Initial collection

	for {
		select {
		case <-ticker.C:
			r.collectAndSend(ctx)
		case <-ctx.Done():
			return
		}
	}
}

func (r *simpleReceiver) collectAndSend(ctx context.Context) {
	// Create a simple demo metric
	metrics := pmetric.NewMetrics()

	// Create a sample configs list - in a real implementation we would load this from disk
	var configs []load.Config

	// Manually create a config for demonstration purposes
	sampleConfig := load.Config{
		Name:     "sample-flex-config",
		FileName: "sample.yml",
	}
	configs = append(configs, sampleConfig)

	// If configs were loaded, execute them
	if len(configs) > 0 {
		r.logger.Info("Running flex with loaded configurations", zap.Int("configCount", len(configs)))
		for _, cfg := range configs {
			// Run the flex config and collect metrics
			r.runFlexConfig(cfg, metrics)
		}
	} else {
		// Just add a sample metric if no configs were loaded
		r.logger.Info("No flex configurations loaded, using sample metric")
		metric := metrics.ResourceMetrics().AppendEmpty().ScopeMetrics().AppendEmpty().Metrics().AppendEmpty()
		metric.SetName("flexdemo.sample")
		metric.SetDescription("A sample metric from the simple flex receiver")

		dp := metric.SetEmptyGauge().DataPoints().AppendEmpty()
		dp.SetTimestamp(pcommon.NewTimestampFromTime(time.Now()))
		dp.SetDoubleValue(100.0)
	}

	r.logger.Info("Collected metrics from simple flex receiver")

	if err := r.consumer.ConsumeMetrics(ctx, metrics); err != nil {
		r.logger.Error("Failed to consume metrics", zap.Error(err))
	}
}

// runFlexConfig executes a specific nri-flex configuration and adds metrics to the metrics object
func (r *simpleReceiver) runFlexConfig(cfg load.Config, metrics pmetric.Metrics) {
	r.logger.Info("Running flex config", zap.String("name", cfg.Name))

	// Here we would process the config using nri-flex's core functionality
	// This is a simplified version that would need to be expanded to fully utilize nri-flex

	// As a placeholder, just add a metric with the config name
	rm := metrics.ResourceMetrics().AppendEmpty()
	sm := rm.ScopeMetrics().AppendEmpty()
	metric := sm.Metrics().AppendEmpty()

	metric.SetName("flex.config.run")
	metric.SetDescription("Execution of a flex configuration")

	dp := metric.SetEmptyGauge().DataPoints().AppendEmpty()
	dp.SetTimestamp(pcommon.NewTimestampFromTime(time.Now()))
	dp.SetDoubleValue(1.0)

	// Add attributes about the config
	dp.Attributes().PutStr("config.name", cfg.Name)
	if cfg.FileName != "" {
		dp.Attributes().PutStr("config.file", cfg.FileName)
	}
}
