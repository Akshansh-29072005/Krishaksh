package handler

import (
	"net/http"
	"strconv"

	"github.com/aarcsx/krisho-backend/internal/core/response"
	"github.com/aarcsx/krisho-backend/internal/modules/company/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CompanyHandler struct{ service service.CompanyService }

func NewCompanyHandler(s service.CompanyService) *CompanyHandler { return &CompanyHandler{service: s} }

func (h *CompanyHandler) Dashboard(c *gin.Context) {
	companyID, err := uuid.Parse(c.Query("company_id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "company_id required")
		return
	}
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))
	out, err := h.service.Dashboard(c.Request.Context(), companyID, days)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed")
		return
	}
	response.Success(c, http.StatusOK, "company dashboard", gin.H{
		"campaign_performance":       out,
		"product_performance":        out,
		"recommendation_performance": out,
		"dashboard_metrics":          out,
	})
}
