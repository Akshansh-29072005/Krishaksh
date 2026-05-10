package repository

import (
	"context"
	"time"

	"github.com/aarcsx/krisho-backend/internal/database"
	"github.com/aarcsx/krisho-backend/internal/models"
	"github.com/google/uuid"
)

type SupportRepository interface {
	CreateTicket(ctx context.Context, ticket *models.SupportTicket) error
	GetTicketsByUser(ctx context.Context, userID uuid.UUID) ([]*models.SupportTicket, error)
	GetTicketByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*models.SupportTicket, error)
	UpdateTicketStatus(ctx context.Context, id uuid.UUID, status string) error
	AddMessage(ctx context.Context, msg *models.SupportMessage) error
	GetMessages(ctx context.Context, ticketID uuid.UUID) ([]*models.SupportMessage, error)
}

type supportRepoImpl struct {
	db *database.DB
}

func NewSupportRepository(db *database.DB) SupportRepository {
	return &supportRepoImpl{db: db}
}

func (r *supportRepoImpl) CreateTicket(ctx context.Context, ticket *models.SupportTicket) error {
	q := `INSERT INTO support_tickets (id, user_id, title, description, status, priority, created_at, updated_at)
	      VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.Pool.Exec(ctx, q, ticket.ID, ticket.UserID, ticket.Title, ticket.Description,
		ticket.Status, ticket.Priority, ticket.CreatedAt, ticket.UpdatedAt)
	return err
}

func (r *supportRepoImpl) GetTicketsByUser(ctx context.Context, userID uuid.UUID) ([]*models.SupportTicket, error) {
	q := `SELECT id, user_id, title, description, status, priority, assigned_to, created_at, updated_at, resolved_at
	      FROM support_tickets WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.Pool.Query(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tickets []*models.SupportTicket
	for rows.Next() {
		t := &models.SupportTicket{}
		if err := rows.Scan(&t.ID, &t.UserID, &t.Title, &t.Description, &t.Status, &t.Priority,
			&t.AssignedTo, &t.CreatedAt, &t.UpdatedAt, &t.ResolvedAt); err != nil {
			return nil, err
		}
		tickets = append(tickets, t)
	}
	return tickets, nil
}

func (r *supportRepoImpl) GetTicketByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*models.SupportTicket, error) {
	q := `SELECT id, user_id, title, description, status, priority, assigned_to, created_at, updated_at, resolved_at
	      FROM support_tickets WHERE id = $1 AND user_id = $2`
	t := &models.SupportTicket{}
	err := r.db.Pool.QueryRow(ctx, q, id, userID).Scan(&t.ID, &t.UserID, &t.Title, &t.Description, &t.Status,
		&t.Priority, &t.AssignedTo, &t.CreatedAt, &t.UpdatedAt, &t.ResolvedAt)
	return t, err
}

func (r *supportRepoImpl) UpdateTicketStatus(ctx context.Context, id uuid.UUID, status string) error {
	var resolvedAt *time.Time
	if status == "resolved" || status == "closed" {
		now := time.Now()
		resolvedAt = &now
	}
	q := `UPDATE support_tickets SET status = $1, resolved_at = $2, updated_at = NOW() WHERE id = $3`
	_, err := r.db.Pool.Exec(ctx, q, status, resolvedAt, id)
	return err
}

func (r *supportRepoImpl) AddMessage(ctx context.Context, msg *models.SupportMessage) error {
	q := `INSERT INTO support_messages (id, ticket_id, sender_id, sender_role, body, created_at)
	      VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.Pool.Exec(ctx, q, msg.ID, msg.TicketID, msg.SenderID, msg.SenderRole, msg.Body, msg.CreatedAt)
	return err
}

func (r *supportRepoImpl) GetMessages(ctx context.Context, ticketID uuid.UUID) ([]*models.SupportMessage, error) {
	q := `SELECT id, ticket_id, sender_id, sender_role, body, created_at
	      FROM support_messages WHERE ticket_id = $1 ORDER BY created_at ASC`
	rows, err := r.db.Pool.Query(ctx, q, ticketID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []*models.SupportMessage
	for rows.Next() {
		m := &models.SupportMessage{}
		if err := rows.Scan(&m.ID, &m.TicketID, &m.SenderID, &m.SenderRole, &m.Body, &m.CreatedAt); err != nil {
			return nil, err
		}
		msgs = append(msgs, m)
	}
	return msgs, nil
}
