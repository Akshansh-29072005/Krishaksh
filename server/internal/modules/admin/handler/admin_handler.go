package handler

import (
	"net/http"
	"strconv"

	"github.com/aarcsx/krisho-backend/internal/core/response"
	"github.com/aarcsx/krisho-backend/internal/models"
	"github.com/aarcsx/krisho-backend/internal/modules/admin/dto"
	"github.com/aarcsx/krisho-backend/internal/modules/admin/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AdminHandler struct{ service service.AdminService }

func NewAdminHandler(s service.AdminService) *AdminHandler { return &AdminHandler{service: s} }

func (h *AdminHandler) ListUsers(c *gin.Context) {
	l, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	u, err := h.service.ListUsers(c.Request.Context(), l)
	if err != nil {
		response.Error(c, 500, "failed")
		return
	}
	response.Success(c, 200, "users", u)
}
func (h *AdminHandler) ChangeRole(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, 400, "invalid")
		return
	}
	var req dto.ChangeRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "invalid")
		return
	}
	if err := h.service.ChangeRole(c.Request.Context(), id, req.RoleID); err != nil {
		response.Error(c, 500, "failed")
		return
	}
	response.Success(c, 200, "updated", nil)
}
func (h *AdminHandler) SuspendUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, 400, "invalid")
		return
	}
	var req dto.SuspendUserRequest
	_ = c.ShouldBindJSON(&req)
	if err := h.service.SuspendUser(c.Request.Context(), id, req.Suspended); err != nil {
		response.Error(c, 500, "failed")
		return
	}
	response.Success(c, 200, "updated", nil)
}

func (h *AdminHandler) CreateDisease(c *gin.Context) {
	var d models.Disease
	if err := c.ShouldBindJSON(&d); err != nil {
		response.Error(c, 400, "invalid")
		return
	}
	d.ID = uuid.New()
	if err := h.service.CreateDisease(c.Request.Context(), &d); err != nil {
		response.Error(c, 500, "failed")
		return
	}
	response.Success(c, 201, "created", d)
}
func (h *AdminHandler) UpdateDisease(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, 400, "invalid")
		return
	}
	var d models.Disease
	if err := c.ShouldBindJSON(&d); err != nil {
		response.Error(c, 400, "invalid")
		return
	}
	d.ID = id
	if err := h.service.UpdateDisease(c.Request.Context(), &d); err != nil {
		response.Error(c, 500, "failed")
		return
	}
	response.Success(c, 200, "updated", nil)
}
func (h *AdminHandler) DeleteDisease(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, 400, "invalid")
		return
	}
	if err := h.service.DeleteDisease(c.Request.Context(), id); err != nil {
		response.Error(c, 500, "failed")
		return
	}
	response.Success(c, 200, "deleted", nil)
}

func (h *AdminHandler) CreateProduct(c *gin.Context) {
	var p models.Product
	if err := c.ShouldBindJSON(&p); err != nil {
		response.Error(c, 400, "invalid")
		return
	}
	p.ID = uuid.New()
	if err := h.service.CreateProduct(c.Request.Context(), &p); err != nil {
		response.Error(c, 500, "failed")
		return
	}
	response.Success(c, 201, "created", p)
}
func (h *AdminHandler) UpdateProduct(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, 400, "invalid")
		return
	}
	var p models.Product
	if err := c.ShouldBindJSON(&p); err != nil {
		response.Error(c, 400, "invalid")
		return
	}
	p.ID = id
	if err := h.service.UpdateProduct(c.Request.Context(), &p); err != nil {
		response.Error(c, 500, "failed")
		return
	}
	response.Success(c, 200, "updated", nil)
}
func (h *AdminHandler) DeleteProduct(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, 400, "invalid")
		return
	}
	if err := h.service.DeleteProduct(c.Request.Context(), id); err != nil {
		response.Error(c, 500, "failed")
		return
	}
	response.Success(c, 200, "deleted", nil)
}

func (h *AdminHandler) ListSupport(c *gin.Context) {
	l, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	t, err := h.service.ListSupport(c.Request.Context(), l)
	if err != nil {
		response.Error(c, 500, "failed")
		return
	}
	response.Success(c, 200, "tickets", t)
}
func (h *AdminHandler) UpdateSupport(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, 400, "invalid")
		return
	}
	var req dto.UpdateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "invalid")
		return
	}
	if err := h.service.UpdateSupport(c.Request.Context(), id, req.Status, req.AssignedTo); err != nil {
		response.Error(c, 500, "failed")
		return
	}
	response.Success(c, 200, "updated", nil)
}

func (h *AdminHandler) ListCampaigns(c *gin.Context) {
	l, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	x, err := h.service.ListCampaigns(c.Request.Context(), l)
	if err != nil {
		response.Error(c, 500, "failed")
		return
	}
	response.Success(c, 200, "campaigns", x)
}
func (h *AdminHandler) UpdateCampaign(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, 400, "invalid")
		return
	}
	var req dto.CampaignStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "invalid")
		return
	}
	if err := h.service.UpdateCampaign(c.Request.Context(), id, req.IsActive); err != nil {
		response.Error(c, 500, "failed")
		return
	}
	response.Success(c, 200, "updated", nil)
}

var _ = http.StatusOK
