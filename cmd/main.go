package main

import (
	"context"
	"encoding/json"
	"log"

	"project_security_one/awsclient"
	"project_security_one/config"
	"project_security_one/internal/services"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
)

func main() {
	config.LoadEnv()

	awsCfg, err := awsConfig.LoadDefaultConfig(context.TODO(),
		awsConfig.WithRegion(config.GetEnv("AWS_REGION", "default-region")),
	)
	if err != nil {
		log.Fatal("Unable to load AWS SDK config, ", err)
	}

	// Init Clients
	s3Client := awsclient.NewS3Client(awsCfg)
	costClient := awsclient.NewCostClient(awsCfg)
	cwCfg, err := awsConfig.LoadDefaultConfig(context.TODO(),
		awsConfig.WithRegion("ap-south-1"),
	)
	if err != nil {
		log.Fatal("Unable to load AWS SDK config for CloudWatch, ", err)
	}
	cwClient := awsclient.NewCWClient(cwCfg)
	// Fetch data
	buckets := s3Client.ListBuckets()
	cost := costClient.GetDailyS3Cost()
	bucketName := config.GetEnv("S3_BUCKET_NAME", "default-bucket-name")
	cwClient.ListAvailableMetrics(bucketName)
	metrics := cwClient.GetBucketMetrics(bucketName)

	// Combine into final payload
	payload := map[string]interface{}{
		"Buckets": buckets,
		"Cost":    cost,
		"Metrics": metrics,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Fatal("Error marshalling JSON: ", err)
	}
	log.Printf("Final Buckets: %+v", buckets)
	log.Printf("Final Cost: %+v", cost)
	log.Printf("Final Metrics: %+v", metrics)

	// Push using service
	endpoint := config.GetEnv("PUSH_ENDPOINT", "default-endpoint")
	services.PushData(endpoint, jsonData)

	log.Println("✅ All data pushed successfully!")
}
