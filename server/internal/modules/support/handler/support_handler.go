package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/aarcsx/krisho-backend/internal/core/response"
	"github.com/aarcsx/krisho-backend/internal/modules/support/dto"
	"github.com/aarcsx/krisho-backend/internal/modules/support/service"
	"github.com/aarcsx/krisho-backend/internal/observability"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SupportHandler struct {
	service service.SupportService
}

func NewSupportHandler(svc service.SupportService) *SupportHandler {
	return &SupportHandler{service: svc}
}

func (h *SupportHandler) CreateTicket(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uuid.UUID)

	var req dto.CreateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	ticket, err := h.service.CreateTicket(c.Request.Context(), userID, req)
	if err != nil {
		observability.InitLogger().Error("support_create_ticket_failed", "error", err.Error(), "user_id", userID.String(), "path", c.FullPath())
		response.Error(c, http.StatusInternalServerError, "Failed to create ticket")
		return
	}
	response.Success(c, http.StatusCreated, "Ticket created successfully", ticket)
}

func (h *SupportHandler) GetTickets(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uuid.UUID)

	tickets, err := h.service.GetTickets(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch tickets")
		return
	}
	response.Success(c, http.StatusOK, "Tickets fetched", tickets)
}

func (h *SupportHandler) GetTicket(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uuid.UUID)

	ticketID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ticket ID")
		return
	}

	ticket, messages, err := h.service.GetTicketWithThread(c.Request.Context(), userID, ticketID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Ticket not found")
		return
	}
	response.Success(c, http.StatusOK, "Ticket fetched", gin.H{
		"ticket":   ticket,
		"messages": messages,
	})
}

func (h *SupportHandler) SendMessage(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uuid.UUID)

	ticketID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ticket ID")
		return
	}

	var req dto.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request")
		return
	}

	msg, err := h.service.SendMessage(c.Request.Context(), userID, "farmer", ticketID, req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to send message")
		return
	}
	response.Success(c, http.StatusCreated, "Message sent", msg)
}

func (h *SupportHandler) RequestCallback(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uuid.UUID)

	ticketID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ticket ID")
		return
	}

	ticket, err := h.service.RequestCallback(c.Request.Context(), userID, ticketID)
	if err != nil {
		if errors.Is(err, service.ErrPhoneNumberRequired) {
			response.Error(c, http.StatusBadRequest, "Phone number required to request callback")
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to request callback")
		return
	}
	response.Success(c, http.StatusOK, "Callback requested successfully", ticket)
}

func (h *SupportHandler) UploadVoice(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(uuid.UUID)

	ticketID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid ticket ID")
		return
	}

	file, header, err := c.Request.FormFile("voice")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Failed to get voice file")
		return
	}
	defer file.Close()

	// Validate file type
	if !strings.HasPrefix(header.Header.Get("Content-Type"), "audio/") {
		response.Error(c, http.StatusBadRequest, "Invalid file type, must be audio")
		return
	}

	// Validate file size (max 10MB)
	if header.Size > 10<<20 {
		response.Error(c, http.StatusBadRequest, "File too large, max 10MB")
		return
	}

	attachment, err := h.service.UploadVoiceAttachment(c.Request.Context(), userID, ticketID, file, header.Filename, header.Size)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to upload voice file")
		return
	}
	response.Success(c, http.StatusCreated, "Voice file uploaded successfully", attachment)
}
