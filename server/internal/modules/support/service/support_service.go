package service

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/aarcsx/krisho-backend/internal/models"
	"github.com/aarcsx/krisho-backend/internal/modules/support/dto"
	"github.com/aarcsx/krisho-backend/internal/modules/support/repository"
	"github.com/aarcsx/krisho-backend/pkg/s3"
	"github.com/google/uuid"
)

type SupportService interface {
	CreateTicket(ctx context.Context, userID uuid.UUID, req dto.CreateTicketRequest) (*models.SupportTicket, error)
	GetTickets(ctx context.Context, userID uuid.UUID) ([]*models.SupportTicket, error)
	GetTicketWithThread(ctx context.Context, userID uuid.UUID, ticketID uuid.UUID) (*models.SupportTicket, []*models.SupportMessage, error)
	SendMessage(ctx context.Context, userID uuid.UUID, role string, ticketID uuid.UUID, req dto.SendMessageRequest) (*models.SupportMessage, error)
	RequestCallback(ctx context.Context, userID uuid.UUID, ticketID uuid.UUID) (*models.SupportTicket, error)
	UploadVoiceAttachment(ctx context.Context, userID uuid.UUID, ticketID uuid.UUID, file io.Reader, filename string, size int64) (*models.SupportAttachment, error)
}

type supportServiceImpl struct {
	repo     repository.SupportRepository
	s3Client s3.S3Client
	bucket   string
}

func NewSupportService(repo repository.SupportRepository, s3Client s3.S3Client, bucket string) SupportService {
	return &supportServiceImpl{repo: repo, s3Client: s3Client, bucket: bucket}
}

func (s *supportServiceImpl) CreateTicket(ctx context.Context, userID uuid.UUID, req dto.CreateTicketRequest) (*models.SupportTicket, error) {
	priority := req.Priority
	if priority == "" {
		priority = "medium"
	}

	ticket := &models.SupportTicket{
		ID:                uuid.New(),
		UserID:            userID,
		Title:             req.Title,
		Description:       req.Description,
		Status:            "open",
		Priority:          priority,
		CallbackRequested: req.RequestCallback,
		CallbackStatus: func() string {
			if req.RequestCallback {
				return "pending"
			}
			return "none"
		}(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.CreateTicket(ctx, ticket); err != nil {
		return nil, fmt.Errorf("failed to create ticket: %w", err)
	}
	return ticket, nil
}

func (s *supportServiceImpl) GetTickets(ctx context.Context, userID uuid.UUID) ([]*models.SupportTicket, error) {
	return s.repo.GetTicketsByUser(ctx, userID)
}

func (s *supportServiceImpl) GetTicketWithThread(ctx context.Context, userID uuid.UUID, ticketID uuid.UUID) (*models.SupportTicket, []*models.SupportMessage, error) {
	ticket, err := s.repo.GetTicketByID(ctx, ticketID, userID)
	if err != nil {
		return nil, nil, fmt.Errorf("ticket not found: %w", err)
	}
	messages, err := s.repo.GetMessages(ctx, ticketID)
	if err != nil {
		return ticket, nil, nil
	}
	return ticket, messages, nil
}

func (s *supportServiceImpl) RequestCallback(ctx context.Context, userID uuid.UUID, ticketID uuid.UUID) (*models.SupportTicket, error) {
	if err := s.repo.RequestCallback(ctx, ticketID, userID); err != nil {
		return nil, fmt.Errorf("failed to request callback: %w", err)
	}
	ticket, err := s.repo.GetTicketByID(ctx, ticketID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to load ticket after callback request: %w", err)
	}
	return ticket, nil
}

func (s *supportServiceImpl) UploadVoiceAttachment(ctx context.Context, userID uuid.UUID, ticketID uuid.UUID, file io.Reader, filename string, size int64) (*models.SupportAttachment, error) {
	// Verify ticket exists and belongs to user
	_, err := s.repo.GetTicketByID(ctx, ticketID, userID)
	if err != nil {
		return nil, fmt.Errorf("ticket not found: %w", err)
	}

	// Generate unique key for S3
	key := fmt.Sprintf("support/voice/%s/%s", ticketID, uuid.New().String()+filepath.Ext(filename))

	// Upload to S3
	url, err := s.s3Client.UploadImage(ctx, s.bucket, key, file)
	if err != nil {
		return nil, fmt.Errorf("failed to upload voice file: %w", err)
	}

	// Create attachment record
	attachment := &models.SupportAttachment{
		ID:            uuid.New(),
		TicketID:      ticketID,
		UploaderID:    userID,
		FileURL:       url,
		FileType:      "voice",
		FileSizeBytes: size,
		CreatedAt:     time.Now(),
	}

	if err := s.repo.CreateAttachment(ctx, attachment); err != nil {
		return nil, fmt.Errorf("failed to save attachment: %w", err)
	}

	return attachment, nil
}

func (s *supportServiceImpl) SendMessage(ctx context.Context, userID uuid.UUID, role string, ticketID uuid.UUID, req dto.SendMessageRequest) (*models.SupportMessage, error) {
	msg := &models.SupportMessage{
		ID:         uuid.New(),
		TicketID:   ticketID,
		SenderID:   userID,
		SenderRole: role,
		Body:       req.Body,
		CreatedAt:  time.Now(),
	}
	if err := s.repo.AddMessage(ctx, msg); err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}
	// When a farmer replies, auto-transition ticket from awaiting_user back to under_review
	if role == "farmer" {
		_ = s.repo.UpdateTicketStatus(ctx, ticketID, "under_review")
	}
	return msg, nil
}
