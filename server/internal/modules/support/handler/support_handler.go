package handler

import (
	"net/http"

	"github.com/aarcsx/krishaksh-backend/internal/core/response"
	"github.com/aarcsx/krishaksh-backend/internal/modules/support/dto"
	"github.com/aarcsx/krishaksh-backend/internal/modules/support/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SupportHandler struct {
	service service.SupportService
}

func NewSupportHandler(s service.SupportService) *SupportHandler {
	return &SupportHandler{service: s}
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
