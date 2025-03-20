package awsclient

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
)

type CostClient struct {
	Client *costexplorer.Client
}

func NewCostClient(cfg aws.Config) *CostClient {
	return &CostClient{Client: costexplorer.NewFromConfig(cfg)}
}

func (c *CostClient) GetDailyS3Cost() map[string]string {
	costData := make(map[string]string)
	end := time.Now()
	start := end.AddDate(0, 0, -1)

	input := &costexplorer.GetCostAndUsageInput{
		TimePeriod: &types.DateInterval{
			Start: aws.String(start.Format("2006-01-02")),
			End:   aws.String(end.Format("2006-01-02")),
		},
		Granularity: types.GranularityDaily,
		Metrics:     []string{"UnblendedCost"},
		Filter: &types.Expression{
			Dimensions: &types.DimensionValues{
				Key:    types.DimensionService,
				Values: []string{"Amazon Simple Storage Service"},
			},
		},
	}

	resp, err := c.Client.GetCostAndUsage(context.TODO(), input)
	if err != nil {
		log.Fatal(err)
	}

	for _, result := range resp.ResultsByTime {
		costData[*result.TimePeriod.Start] = *result.Total["UnblendedCost"].Amount
	}
	return costData
}
