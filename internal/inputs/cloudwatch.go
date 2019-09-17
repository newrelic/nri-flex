package inputs

import (
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/newrelic/nri-flex/internal/formatter"
	"github.com/newrelic/nri-flex/internal/load"
	"github.com/newrelic/nri-flex/internal/logger"
)

var metricIDAWSMetric map[string]load.CloudwatchMetric
var namespaceMetricNameAWSStats map[string]string
var awsGetMetricDataOutput *cloudwatch.GetMetricDataOutput

// RunCloudwatch Input
func RunCloudwatch(dataStore *[]interface{}, config *load.Config, apiNo int) {
	metricIDAWSMetric = make(map[string]load.CloudwatchMetric)
	namespaceMetricNameAWSStats = make(map[string]string)

	cwConfig := config.APIs[apiNo]
	var sess *session.Session
	sharedConfigFiles := []string{}

	if cwConfig.Cloudwatch.CredentialFile != "" {
		sharedConfigFiles = append(sharedConfigFiles, cwConfig.Cloudwatch.CredentialFile)
	}
	if cwConfig.Cloudwatch.ConfigFile != "" {
		sharedConfigFiles = append(sharedConfigFiles, cwConfig.Cloudwatch.ConfigFile)
	}

	if len(sharedConfigFiles) > 0 {
		logger.Flex("debug", nil, "aws cloudwatch using custom credentials and/or config", false)
		sess = session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
			SharedConfigFiles: sharedConfigFiles,
		}))
	} else {
		logger.Flex("debug", nil, "aws cloudwatch using default credentials", false)
		sess = session.Must(session.NewSession(&aws.Config{
			Region: aws.String(cwConfig.Cloudwatch.Region),
		}))
	}

	cwClient := cloudwatch.New(sess)

	awsGetMetricDataInput, err := getMetricList(cwClient, &cwConfig)
	if err != nil {
		logger.Flex("error", err, "cloudwatch unable to build metric list for querying", false)
	} else {
		// fmt.Println(awsGetMetricDataInput)

		awsMetricDataQueriesAll := awsGetMetricDataInput.MetricDataQueries
		numberOfMetrics := len(awsMetricDataQueriesAll)
		// stay under AWS API limit per request
		metricsPerRequest := 99
		if cwConfig.Cloudwatch.MetricsPerRequest > 0 {
			metricsPerRequest = cwConfig.Cloudwatch.MetricsPerRequest
		}
		fmt.Println(metricsPerRequest, numberOfMetrics)
		for j := 0; j <= numberOfMetrics/metricsPerRequest && numberOfMetrics != 0; j++ {
			startInd := j * metricsPerRequest
			endInd := (j + 1) * metricsPerRequest

			if endInd > numberOfMetrics {
				endInd = numberOfMetrics
			}

			awsMetricDataQueriesSmall := awsMetricDataQueriesAll[startInd:endInd]
			awsGetMetricDataInput.MetricDataQueries = awsMetricDataQueriesSmall
			awsGetMetricDataOutput, err = getAWSMetricData(awsGetMetricDataInput, cwClient)

			for _, metricOutput := range awsGetMetricDataOutput.MetricDataResults {
				// only take completed metric results with values
				if *metricOutput.StatusCode == "Complete" && len(metricOutput.Values) != 0 {

				}
			}

		}
	}
}

func getMetricList(svc *cloudwatch.CloudWatch, cwConfig *load.API) (*cloudwatch.GetMetricDataInput, error) {
	var mlqs []*cloudwatch.MetricDataQuery

	for j, cwMetric := range cwConfig.Cloudwatch.Metrics {
		res, err := getAWSMetricList(cwMetric.MetricName, cwMetric.Namespace, svc)
		if err != nil {
			logger.Flex("error", err, "cloudwatch failed to get metric list", false)
		} else {
			for i, awsMetric := range res.Metrics {
				if len(awsMetric.Dimensions) > 0 {
					Namespace := cwMetric.Namespace
					MetricName := cwMetric.MetricName
					Dimensions := []*cloudwatch.Dimension{}

					for _, dimension := range awsMetric.Dimensions {
						if len(cwMetric.Filter) == 0 {
							cwMetric.Filter = append(cwMetric.Filter, map[string]string{".*": ".*"})
						}

						for _, filter := range cwMetric.Filter {
							for regex1, regex2 := range filter {
								if formatter.KvFinder("regex", *dimension.Name, regex1) && formatter.KvFinder("regex", *dimension.Value, regex2) {
									var Dimension *cloudwatch.Dimension
									Dimension = new(cloudwatch.Dimension)
									Dimension.Name = dimension.Name
									Dimension.Value = dimension.Value
									Dimensions = append(Dimensions, Dimension)
								}
							}
						}
					}

					if len(Dimensions) != 0 {
						ml := load.CloudwatchMetric{
							Namespace:  cwMetric.Namespace,
							MetricName: cwMetric.MetricName,
							Dimensions: Dimensions,
						}

						Metric := cloudwatch.Metric{Dimensions: Dimensions,
							MetricName: &MetricName,
							Namespace:  &Namespace,
						}

						Period := int64(60)

						// set statistics default to atleast average
						if len(cwMetric.Statistics) == 0 {
							cwMetric.Statistics = append(cwMetric.Statistics, "Average")
						}

						for _, stat := range cwMetric.Statistics {

							MetricStat := cloudwatch.MetricStat{Metric: &Metric,
								Period: &Period,
								Stat:   &stat,
							}
							ID := "id" + (strconv.Itoa(j)) + (strconv.Itoa(i)) + stat
							t := true

							metricIDAWSMetric[ID] = ml
							namespaceMetricNameAWSStats[ID] = stat

							mlq := cloudwatch.MetricDataQuery{
								Id:         &ID,
								MetricStat: &MetricStat,
								ReturnData: &t,
							}

							mlqs = append(mlqs, &mlq)
						}
					}

				}
			}

		}
	}

	endTime := time.Now().UTC()
	startTime := endTime.Add(time.Duration(-cwConfig.Cloudwatch.NativeInterval) * time.Minute)
	metricListRequest := cloudwatch.GetMetricDataInput{
		StartTime:         &startTime,
		EndTime:           &endTime,
		MetricDataQueries: mlqs,
	}
	return &metricListRequest, nil
}

func getAWSMetricList(metricname string, namespace string, svc *cloudwatch.CloudWatch) (*cloudwatch.ListMetricsOutput, error) {
	var nextToken *string
	var err error
	var result *cloudwatch.ListMetricsOutput
	var input cloudwatch.ListMetricsInput
	input.MetricName = aws.String(metricname)
	input.Namespace = aws.String(namespace)
	var metrics []*cloudwatch.Metric

	for {
		load.StatusCounterIncrement("cloudwatchApiCalls")
		result, err = svc.ListMetrics(&input)
		nextToken = result.NextToken
		metrics = append(metrics, result.Metrics...)
		if nextToken == nil {
			break
		} else {
			input.NextToken = nextToken
		}
	}

	result.Metrics = metrics
	return result, err
}

func getAWSMetricData(mlrs *cloudwatch.GetMetricDataInput, svc *cloudwatch.CloudWatch) (*cloudwatch.GetMetricDataOutput, error) {
	load.StatusCounterIncrement("cloudwatchApiCalls")
	result, err := svc.GetMetricData(mlrs)

	return result, err
}
