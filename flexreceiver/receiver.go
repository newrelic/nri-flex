package flexreceiver

import (
	"context"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
	"go.uber.org/zap"
)

// flexReceiver is the struct that implements the OTel receiver interface.
type flexReceiver struct {
	cfg      *Config
	settings receiver.Settings
	consumer consumer.Metrics
	cancel   context.CancelFunc
	logger   *zap.Logger
	bridge   *FlexBridge
}

// newFlexReceiver creates a new instance of the flexReceiver.
func newFlexReceiver(cfg *Config, settings receiver.Settings, consumer consumer.Metrics) *flexReceiver {
	return &flexReceiver{
		cfg:      cfg,
		settings: settings,
		consumer: consumer,
		logger:   settings.Logger,
	}
}

// Start is called when the receiver is started.
func (r *flexReceiver) Start(ctx context.Context, host component.Host) error {
	r.logger.Info("Starting nri-flex receiver")
	var runCtx context.Context
	runCtx, r.cancel = context.WithCancel(ctx)

	r.bridge = NewFlexBridge(r.cfg, r.logger)

	go r.runCollectionLoop(runCtx)
	return nil
}

// Shutdown is called when the receiver is stopped.
func (r *flexReceiver) Shutdown(ctx context.Context) error {
	r.logger.Info("Shutting down nri-flex receiver")
	if r.cancel != nil {
		r.cancel()
	}
	return nil
}

// runCollectionLoop runs the metric collection at the configured interval.
func (r *flexReceiver) runCollectionLoop(ctx context.Context) {
	ticker := time.NewTicker(r.cfg.CollectionInterval)
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

// collectAndSend triggers the Flex data collection and sends the converted metrics to the consumer.
func (r *flexReceiver) collectAndSend(ctx context.Context) {
	metrics, err := r.bridge.CollectMetrics(ctx)
	if err != nil {
		r.logger.Error("Failed to collect metrics from Flex", zap.Error(err))
		return
	}

	if metrics.MetricCount() == 0 {
		r.logger.Debug("No metrics were collected by Flex")
		return
	}

	if err := r.consumer.ConsumeMetrics(ctx, metrics); err != nil {
		r.logger.Error("Failed to consume metrics", zap.Error(err))
	}
}
