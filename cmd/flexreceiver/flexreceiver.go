package main

import (
	"log"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/otelcol"

	"github.com/newrelic/nri-flex/flexreceiver"
)

func main() {
	info := component.BuildInfo{
		Command:     "flexreceiver",
		Description: "New Relic Flex OpenTelemetry Receiver",
		Version:     "1.0.0", // Replace with dynamic versioning if needed
	}

	factories, err := otelcol.MakeFactoryMap(
		flexreceiver.NewFactory(),
	)

	if err != nil {
		log.Fatalf("failed to build factories: %v", err)
	}

	app := otelcol.NewCommand(
		otelcol.CollectorSettings{
			BuildInfo: info,
			Factories: func() (otelcol.Factories, error) {
				return otelcol.Factories{
					Receivers: factories,
				}, nil
			},
		},
	)

	if err := app.Execute(); err != nil {
		log.Fatalf("collector server run finished with error: %v", err)
	}
}
