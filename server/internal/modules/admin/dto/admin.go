package dto

import "github.com/google/uuid"

type ChangeRoleRequest struct {
	RoleID int `json:"role_id" binding:"required"`
}
type SuspendUserRequest struct {
	Suspended bool `json:"suspended"`
}
type UpdateTicketRequest struct {
	Status     string     `json:"status" binding:"required"`
	AssignedTo *uuid.UUID `json:"assigned_to"`
}
type CampaignStatusRequest struct {
	IsActive bool `json:"is_active"`
}
