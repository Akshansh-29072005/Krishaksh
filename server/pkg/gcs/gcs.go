package gcs

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"cloud.google.com/go/storage"
)

type GCSClient interface {
	UploadImage(ctx context.Context, bucket string, key string, body io.Reader) (string, error)
	GeneratePresignedURL(ctx context.Context, bucket string, key string, expiration time.Duration) (string, error)
	GeneratePresignedGetURL(ctx context.Context, bucket string, key string, expiration time.Duration) (string, error)
	ExtractKeyFromURL(rawURL string) (string, error)
}

type gcsClientImpl struct {
	client *storage.Client
}

func NewGCSClient(ctx context.Context) (GCSClient, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to create GCS client: %w", err)
	}
	return &gcsClientImpl{client: client}, nil
}

func (g *gcsClientImpl) UploadImage(ctx context.Context, bucket string, key string, body io.Reader) (string, error) {
	wc := g.client.Bucket(bucket).Object(key).NewWriter(ctx)
	if _, err := io.Copy(wc, body); err != nil {
		return "", fmt.Errorf("failed to write object to GCS: %w", err)
	}
	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("failed to finalize GCS object: %w", err)
	}

	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucket, key), nil
}

func (g *gcsClientImpl) GeneratePresignedURL(ctx context.Context, bucket string, key string, expiration time.Duration) (string, error) {
	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  "PUT",
		Expires: time.Now().Add(expiration),
	}
	url, err := g.client.Bucket(bucket).SignedURL(key, opts)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned PUT URL for GCS: %w", err)
	}
	return url, nil
}

func (g *gcsClientImpl) GeneratePresignedGetURL(ctx context.Context, bucket string, key string, expiration time.Duration) (string, error) {
	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  "GET",
		Expires: time.Now().Add(expiration),
	}
	url, err := g.client.Bucket(bucket).SignedURL(key, opts)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned GET URL for GCS: %w", err)
	}
	return url, nil
}

func (g *gcsClientImpl) ExtractKeyFromURL(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	// Handle both https://storage.googleapis.com/bucket/key and https://bucket.storage.googleapis.com/key formats
	if strings.HasPrefix(u.Host, "storage.googleapis.com") {
		return strings.TrimPrefix(u.Path, "/"), nil
	} else {
		// Bucket is subdomain: bucket.storage.googleapis.com/key → key is path after /
		return strings.TrimPrefix(u.Path, "/"), nil
	}
}
