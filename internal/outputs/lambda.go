package outputs

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/newrelic/nri-flex/internal/load"
)

type MyEvent struct {
	Name string `json:"name"`
}

var LambdaEnabled bool
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
func HandleRequest(ctx context.Context, name MyEvent) (string, error) {
	return fmt.Sprintf("Flex Finished - success: %t!", LambdaSuccess), nil
}
