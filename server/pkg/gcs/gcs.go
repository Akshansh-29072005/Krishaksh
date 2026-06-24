package gcs

import (
	"context"
	"fmt"
	"io"
	"net/http"
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
	// Get service account email from metadata server (Cloud Run)
	serviceAccountEmail, err := getServiceAccountEmail(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get service account email: %w", err)
	}

	signer, err := g.client.Bucket(bucket).Signer(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get bucket signer: %w", err)
	}

	opts := &storage.SignedURLOptions{
		Scheme:         storage.SigningSchemeV4,
		Method:         "PUT",
		Expires:        time.Now().Add(expiration),
		GoogleAccessID: serviceAccountEmail,
		SignBytes:      signer.SignBytes,
	}

	url, err := storage.SignedURL(bucket, key, opts)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned PUT URL for GCS: %w", err)
	}
	return url, nil
}

func (g *gcsClientImpl) GeneratePresignedGetURL(ctx context.Context, bucket string, key string, expiration time.Duration) (string, error) {
	serviceAccountEmail, err := getServiceAccountEmail(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get service account email: %w", err)
	}

	signer, err := g.client.Bucket(bucket).Signer(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get bucket signer: %w", err)
	}

	opts := &storage.SignedURLOptions{
		Scheme:         storage.SigningSchemeV4,
		Method:         "GET",
		Expires:        time.Now().Add(expiration),
		GoogleAccessID: serviceAccountEmail,
		SignBytes:      signer.SignBytes,
	}

	url, err := storage.SignedURL(bucket, key, opts)
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

	if strings.HasPrefix(u.Host, "storage.googleapis.com") {
		return strings.TrimPrefix(u.Path, "/"), nil
	} else {
		return strings.TrimPrefix(u.Path, "/"), nil
	}
}

// Helper function to get service account email from Cloud Run metadata server
func getServiceAccountEmail(ctx context.Context) (string, error) {
	// Use Google Cloud metadata server to get default service account email
	client := &http.Client{Timeout: 2 * time.Second}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://metadata.google.internal/computeMetadata/v1/instance/service-accounts/default/email", nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Metadata-Flavor", "Google")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("metadata server returned status: %d", resp.StatusCode)
	}

	emailBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(emailBytes)), nil
}
