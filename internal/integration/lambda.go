package integration

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/newrelic/nri-flex/internal/load"
	"github.com/newrelic/nri-flex/internal/logger"
	logrus "github.com/sirupsen/logrus"
)

// LambdaCheck check if Flex is running within a Lambda and insights url and api key has been supplied
func LambdaCheck() bool {
	if os.Getenv("LAMBDA_TASK_ROOT") == "" {
		return false
	}

	load.LambdaName = os.Getenv("AWS_LAMBDA_FUNCTION_NAME")
	load.Logrus.SetFormatter(&logrus.JSONFormatter{})

	if load.Args.InsightsURL == "" || load.Args.InsightsAPIKey == "" {
		logger.Flex("error", fmt.Errorf("missing insights URL and/or API key"), "", false)
		return false
	}
	return true
}

// Lambda handles lambda invocation
func Lambda() {
	lambda.Start(HandleRequest)
}

// HandleRequest Handles incoming lambda request
func HandleRequest(ctx context.Context, event interface{}) (string, error) {
	logger.Flex("info", nil, "running as lambda", false)
	SetDefaults()

	if event != nil {
		load.IngestData = event
	}

	RunFlex("lambda")

	return fmt.Sprintf("Flex Lambda Complete"), nil
}
