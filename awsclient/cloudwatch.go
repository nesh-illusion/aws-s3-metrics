package awsclient

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
)

type CWClient struct {
	Client *cloudwatch.Client
}

func NewCWClient(cfg aws.Config) *CWClient {
	return &CWClient{Client: cloudwatch.NewFromConfig(cfg)}
}

func (cw *CWClient) GetBucketMetrics(bucketName string) map[string]interface{} {
	metricsData := make(map[string]interface{})

	metrics := []string{"BucketSizeBytes", "NumberOfObjects"}
	queryIDs := []string{"m1", "m2"}

	queries := []types.MetricDataQuery{}

	for idx, metric := range metrics {
		dimensions := []types.Dimension{
			{
				Name:  aws.String("BucketName"),
				Value: aws.String(bucketName),
			},
		}

		// Add StorageType only for BucketSizeBytes
		if metric == "BucketSizeBytes" {
			dimensions = append(dimensions, types.Dimension{
				Name:  aws.String("StorageType"),
				Value: aws.String("StandardStorage"),
			})
		}

		if metric == "NumberOfObjects" {
			dimensions = append(dimensions, types.Dimension{
				Name:  aws.String("StorageType"),
				Value: aws.String("AllStorageTypes"),
			})
		}

		queries = append(queries, types.MetricDataQuery{
			Id: aws.String(queryIDs[idx]),
			MetricStat: &types.MetricStat{
				Metric: &types.Metric{
					Namespace:  aws.String("AWS/S3"),
					MetricName: aws.String(metric),
					Dimensions: dimensions,
				},
				Period: aws.Int32(86400), // Daily data points
				Stat:   aws.String("Average"),
			},
			Label: aws.String(metric), // Makes result parsing easier
		})
	}

	input := &cloudwatch.GetMetricDataInput{
		StartTime:         aws.Time(time.Now().AddDate(0, 0, -3)), // Make ot to -1 for 1 day before moving to prod
		EndTime:           aws.Time(time.Now()),
		MetricDataQueries: queries,
	}

	resp, err := cw.Client.GetMetricData(context.TODO(), input)
	if err != nil {
		log.Printf("Error fetching metrics: %v", err)
		return metricsData
	}

	// Process results
	for _, result := range resp.MetricDataResults {
		if len(result.Values) > 0 {
			// Use latest datapoint
			latestValue := result.Values[len(result.Values)-1]
			log.Printf("Metric %s: %f", *result.Label, latestValue)
			metricsData[*result.Label] = latestValue
		} else {
			log.Printf("No datapoints found for metric %s", *result.Label)
		}
	}

	log.Printf("Final Metrics: %+v", metricsData)
	return metricsData
}

func (cw *CWClient) ListAvailableMetrics(bucketName string) {
	input := &cloudwatch.ListMetricsInput{
		Namespace: aws.String("AWS/S3"),
		Dimensions: []types.DimensionFilter{
			{
				Name:  aws.String("BucketName"),
				Value: aws.String(bucketName),
			},
		},
	}

	resp, err := cw.Client.ListMetrics(context.TODO(), input)
	if err != nil {
		log.Printf("Error listing metrics: %v", err)
		return
	}

	if len(resp.Metrics) == 0 {
		log.Printf("No metrics found for bucket: %s", bucketName)
	} else {
		for _, metric := range resp.Metrics {
			log.Printf("Metric: %s, Dimensions: %+v", *metric.MetricName, metric.Dimensions)
		}
	}

	// Optional: Log NextToken if there are more metrics
	if resp.NextToken != nil {
		log.Printf("NextToken: %s", *resp.NextToken)
	}
}
