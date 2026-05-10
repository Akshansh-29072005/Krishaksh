package service

import (
	"context"
	"fmt"
	"time"

	"github.com/aarcsx/krishaksh-backend/internal/models"
	"github.com/aarcsx/krishaksh-backend/internal/modules/scans/dto"
	"github.com/aarcsx/krishaksh-backend/internal/modules/scans/repository"
	"github.com/aarcsx/krishaksh-backend/pkg/queue"
	"github.com/aarcsx/krishaksh-backend/pkg/s3"
	"github.com/google/uuid"
)

type ScanService interface {
	GetPresignedUploadURL(ctx context.Context, userID uuid.UUID, contentType string) (string, string, error)
	ProcessUpload(ctx context.Context, userID uuid.UUID, req dto.CreateScanRequest) (*models.Scan, error)
	GetScanHistory(ctx context.Context, userID uuid.UUID) ([]*models.Scan, error)
	GetScanDetails(ctx context.Context, userID uuid.UUID, scanID uuid.UUID) (*models.Scan, error)
}

type scanServiceImpl struct {
	repo    repository.ScanRepository
	s3      s3.S3Client
	qClient queue.QueueClient
	bucket  string
}

func NewScanService(repo repository.ScanRepository, s3Client s3.S3Client, q queue.QueueClient, bucket string) ScanService {
	return &scanServiceImpl{
		repo:    repo,
		s3:      s3Client,
		qClient: q,
		bucket:  bucket,
	}
}

func (s *scanServiceImpl) GetPresignedUploadURL(ctx context.Context, userID uuid.UUID, contentType string) (string, string, error) {
	// Generate unique image key preventing overlaps
	ext := ".jpg"
	if contentType == "image/png" {
		ext = ".png"
	}
	key := fmt.Sprintf("scans/%s/%s%s", userID.String(), uuid.New().String(), ext)

	// URL expires in exactly 15 minutes, ensuring high security parameters
	url, err := s.s3.GeneratePresignedURL(ctx, s.bucket, key, 15*time.Minute)
	return url, key, err
}

func (s *scanServiceImpl) ProcessUpload(ctx context.Context, userID uuid.UUID, req dto.CreateScanRequest) (*models.Scan, error) {
	// 1. In reality we should ideally hit S3 HeadObject first to verify the image size/integrity via 'req.ImageKey'

	imageURL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", s.bucket, req.ImageKey)

	// 2. Insert record into PGX
	scan := &models.Scan{
		ID:               uuid.New(),
		UserID:           userID,
		ImageURL:         imageURL,
		CropType:         req.CropType,
		PredictionStatus: "QUEUED",
		CreatedAt:        time.Now(),
	}

	if err := s.repo.CreateScan(ctx, scan); err != nil {
		return nil, fmt.Errorf("failed to save scan to db: %w", err)
	}

	// 3. Emit asynchronous task to Asynq allowing ultra-low latency response latency
	err := s.qClient.EnqueueScanTask(scan.ID.String(), scan.ImageURL, scan.CropType)
	if err != nil {
		// Log error, but don't fail standard response
		fmt.Printf("Critical: Failed to enqueue AI task: %v\n", err)
	}

	return scan, nil
}

func (s *scanServiceImpl) GetScanHistory(ctx context.Context, userID uuid.UUID) ([]*models.Scan, error) {
	return s.repo.GetUserScans(ctx, userID)
}

func (s *scanServiceImpl) GetScanDetails(ctx context.Context, userID uuid.UUID, scanID uuid.UUID) (*models.Scan, error) {
	return s.repo.GetScanByID(ctx, scanID, userID)
}
