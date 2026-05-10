package s3

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Client interface {
	UploadImage(ctx context.Context, bucket string, key string, body io.Reader) (string, error)
	GeneratePresignedURL(ctx context.Context, bucket string, key string, expiration time.Duration) (string, error)
}

type s3ClientImpl struct {
	client *s3.Client
}

func NewS3Client(region, accessKey, secretKey string) (S3Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %w", err)
	}

	return &s3ClientImpl{
		client: s3.NewFromConfig(cfg),
	}, nil
}

func (s *s3ClientImpl) UploadImage(ctx context.Context, bucket string, key string, body io.Reader) (string, error) {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   body,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload image to S3: %w", err)
	}

	// Assuming the bucket is public, or we return a structured S3 URI.
	return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucket, key), nil
}

func (s *s3ClientImpl) GeneratePresignedURL(ctx context.Context, bucket string, key string, expiration time.Duration) (string, error) {
	presigner := s3.NewPresignClient(s.client)

	presignedReq, err := presigner.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(expiration))

	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return presignedReq.URL, nil
}
