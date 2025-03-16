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

func (cw *CWClient) GetBucketMetrics(bucketName string) {
	metrics := []string{"NumberOfObjects", "BucketSizeBytes"}

	for _, metric := range metrics {
		input := &cloudwatch.GetMetricStatisticsInput{
			Namespace:  aws.String("AWS/S3"),
			MetricName: aws.String(metric),
			Dimensions: []types.Dimension{
				{
					Name:  aws.String("BucketName"),
					Value: aws.String(bucketName),
				},
				{
					Name:  aws.String("StorageType"),
					Value: aws.String("StandardStorage"),
				},
			},
			StartTime: aws.Time(time.Now().AddDate(0, 0, -1)),
			EndTime:   aws.Time(time.Now()),
			Period:    aws.Int32(86400),
			Statistics: []types.Statistic{
				types.StatisticAverage,
			},
		}

		resp, err := cw.Client.GetMetricStatistics(context.TODO(), input)
		if err != nil {
			log.Printf("Error fetching %s: %v", metric, err)
			continue
		}

		for _, dp := range resp.Datapoints {
			log.Printf("%s for bucket %s on %v: %f", metric, bucketName, *dp.Timestamp, *dp.Average)
		}
	}
}
