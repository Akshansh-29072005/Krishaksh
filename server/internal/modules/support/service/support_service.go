package service

import (
	"context"
	"fmt"
	"time"

	"github.com/aarcsx/krishaksh-backend/internal/models"
	"github.com/aarcsx/krishaksh-backend/internal/modules/support/dto"
	"github.com/aarcsx/krishaksh-backend/internal/modules/support/repository"
	"github.com/google/uuid"
)

type SupportService interface {
	CreateTicket(ctx context.Context, userID uuid.UUID, req dto.CreateTicketRequest) (*models.SupportTicket, error)
	GetTickets(ctx context.Context, userID uuid.UUID) ([]*models.SupportTicket, error)
	GetTicketWithThread(ctx context.Context, userID uuid.UUID, ticketID uuid.UUID) (*models.SupportTicket, []*models.SupportMessage, error)
	SendMessage(ctx context.Context, userID uuid.UUID, role string, ticketID uuid.UUID, req dto.SendMessageRequest) (*models.SupportMessage, error)
}

type supportServiceImpl struct {
	repo repository.SupportRepository
}

func NewSupportService(repo repository.SupportRepository) SupportService {
	return &supportServiceImpl{repo: repo}
}

func (s *supportServiceImpl) CreateTicket(ctx context.Context, userID uuid.UUID, req dto.CreateTicketRequest) (*models.SupportTicket, error) {
	priority := req.Priority
	if priority == "" {
		priority = "medium"
	}

	ticket := &models.SupportTicket{
		ID:          uuid.New(),
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Status:      "open",
		Priority:    priority,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
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
