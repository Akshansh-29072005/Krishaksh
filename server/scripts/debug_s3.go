package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load("../.env")

	region := os.Getenv("AWS_REGION")
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	bucket := os.Getenv("S3_BUCKET")

	if region == "" || accessKey == "" || secretKey == "" || bucket == "" {
		log.Fatalf("Missing env vars: REGION=%s, KEY_ID=%s, BUCKET=%s", region, accessKey, bucket)
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
	)
	if err != nil {
		log.Fatalf("unable to load SDK config: %v", err)
	}

	s3Client := s3.NewFromConfig(cfg)

	// 1. Test Listing Buckets (Basic connectivity)
	_, err = s3Client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		log.Fatalf("Failed to list buckets (connectivity issue): %v", err)
	}
	fmt.Println("SUCCESS: Basic connectivity to AWS established.")

	// 2. Test Presigning
	presigner := s3.NewPresignClient(s3Client)

	key := fmt.Sprintf("test-debug-%d.jpg", time.Now().Unix())
	presignedReq, err := presigner.PresignPutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(15*time.Minute))

	if err != nil {
		log.Fatalf("failed to generate presigned URL: %v", err)
	}

	fmt.Printf("SUCCESS! Presigned URL: %s\n", presignedReq.URL)
}
