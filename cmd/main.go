package main

import (
	"context"
	"log"

	"project_security_one/awsclient"
	"project_security_one/config"

	awsConfig "github.com/aws/aws-sdk-go-v2/config" // Alias to avoid conflict
)

func main() {
	// Load configurations
	config.LoadEnv()

	// Load AWS Config
	awsCfg, err := awsConfig.LoadDefaultConfig(context.TODO(),
		awsConfig.WithRegion(config.GetEnv("AWS_REGION", "ap-south-1")),
	)
	if err != nil {
		log.Fatal("Unable to load AWS SDK config, ", err)
	}

	// Initialize AWS Clients
	s3Client := awsclient.NewS3Client(awsCfg)
	s3Client.ListBuckets()

	costClient := awsclient.NewCostClient(awsCfg)
	costClient.GetDailyS3Cost()

	cwClient := awsclient.NewCWClient(awsCfg)
	bucketName := config.GetEnv("S3_BUCKET_NAME", "default-bucket-name")

	cwClient.GetBucketMetrics(bucketName)

	// Initialize API Router
	//router := routes.SetupRouter()

	// Start API Server
	port := config.GetEnv("PORT", "8080")
	log.Println("API Server running on port", port)
	//log.Fatal(http.ListenAndServe(":"+port, router))
}
