package gcs

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"
	"time"

	credentials "cloud.google.com/go/iam/credentials/apiv1"
	"cloud.google.com/go/iam/credentials/apiv1/credentialspb"
	"cloud.google.com/go/storage"
)

type GCSClient interface {
	UploadImage(ctx context.Context, bucket string, key string, body io.Reader) (string, error)
	GeneratePresignedURL(ctx context.Context, bucket string, key string, expiration time.Duration) (string, error)
	GeneratePresignedGetURL(ctx context.Context, bucket string, key string, expiration time.Duration) (string, error)
	ExtractKeyFromURL(rawURL string) (string, error)
}

type gcsClientImpl struct {
	client              *storage.Client
	serviceAccountEmail string
	iamClient           *credentials.IamCredentialsClient
}

func NewGCSClient(ctx context.Context) (GCSClient, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to create GCS client: %w", err)
	}

	serviceAccountEmail := os.Getenv("GCS_SERVICE_ACCOUNT_EMAIL")
	if serviceAccountEmail == "" {
		return nil, fmt.Errorf("GCS_SERVICE_ACCOUNT_EMAIL environment variable must be set")
	}

	iamClient, err := credentials.NewIamCredentialsClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to create IAM credentials client: %w", err)
	}

	return &gcsClientImpl{
		client:              client,
		serviceAccountEmail: serviceAccountEmail,
		iamClient:           iamClient,
	}, nil
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
		Scheme:         storage.SigningSchemeV4,
		Method:         "PUT",
		Expires:        time.Now().Add(expiration),
		GoogleAccessID: g.serviceAccountEmail,
		SignBytes: func(b []byte) ([]byte, error) {
			req := &credentialspb.SignBlobRequest{
				Name:    fmt.Sprintf("projects/-/serviceAccounts/%s", g.serviceAccountEmail),
				Payload: b,
			}

			resp, err := g.iamClient.SignBlob(ctx, req)
			if err != nil {
				return nil, err
			}

			return resp.SignedBlob, nil
		},
	}

	url, err := storage.SignedURL(bucket, key, opts)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned PUT URL for GCS: %w", err)
	}
	return url, nil
}

func (g *gcsClientImpl) GeneratePresignedGetURL(ctx context.Context, bucket string, key string, expiration time.Duration) (string, error) {
	opts := &storage.SignedURLOptions{
		Scheme:         storage.SigningSchemeV4,
		Method:         "GET",
		Expires:        time.Now().Add(expiration),
		GoogleAccessID: g.serviceAccountEmail,
		SignBytes: func(b []byte) ([]byte, error) {
			req := &credentialspb.SignBlobRequest{
				Name:    fmt.Sprintf("projects/-/serviceAccounts/%s", g.serviceAccountEmail),
				Payload: b,
			}

			resp, err := g.iamClient.SignBlob(ctx, req)
			if err != nil {
				return nil, err
			}

			return resp.SignedBlob, nil
		},
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
