package outputs

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/newrelic/nri-flex/internal/load"
	"github.com/newrelic/nri-flex/internal/logger"
)

// LambdaEnabled flag that lambda is enabled
var LambdaEnabled bool

// LambdaSuccess flag that lambda executed successfully
var LambdaSuccess bool

// LambdaCheck check if Flex is running within a Lambda
func LambdaCheck() {
	if os.Getenv("LAMBDA_TASK_ROOT") == "" {
		LambdaEnabled = false
	} else {
		LambdaEnabled = true
		if os.Getenv("INSIGHTS_URL") == "" || os.Getenv("INSIGHTS_API_KEY") == "" {
			fmt.Println("Missing INSIGHTS_URL and/or INSIGHTS_API_KEY")
			LambdaSuccess = false
			lambda.Start(HandleRequest)
		} else {
			load.Args.ConfigDir = "/var/task/pkg/flexConfigs/"
		}
	}
}

// LambdaFinish wrap up the lambda request
func LambdaFinish() {
	if os.Getenv("INSIGHTS_URL") != "" && os.Getenv("INSIGHTS_API_KEY") != "" {
		load.Args.InsightsURL = os.Getenv("INSIGHTS_URL")
		load.Args.InsightsAPIKey = os.Getenv("INSIGHTS_API_KEY")
		LambdaSuccess = true
	} else {
		LambdaSuccess = false
	}
	lambda.Start(HandleRequest)
}

// HandleRequest Handles lambda request
func HandleRequest(ctx context.Context, event map[string]interface{}) (string, error) {
	awsSource := ""
	if event["source"] != nil {
		switch source := event["source"].(type) {
		case []interface{}:
			if len(source) == 1 {
				awsSource = source[0].(string)
			}
		}
	}
	if awsSource != "" {
		logger.Flex("debug", nil, fmt.Sprintf("aws source detected %v", awsSource), false)
	}
	return fmt.Sprintf("Flex Finished - success: %t!", LambdaSuccess), nil
}
