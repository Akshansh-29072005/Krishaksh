package handler

import (
	"net/http"

	"github.com/aarcsx/krisho-backend/internal/core/response"
	"github.com/aarcsx/krisho-backend/internal/models"
	diseaseRepo "github.com/aarcsx/krisho-backend/internal/modules/diseases/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DiseaseHandler struct {
	repo diseaseRepo.DiseaseRepository
}

func NewDiseaseHandler(repo diseaseRepo.DiseaseRepository) *DiseaseHandler {
	return &DiseaseHandler{repo: repo}
}

func (h *DiseaseHandler) GetAll(c *gin.Context) {
	diseases, err := h.repo.GetAll(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch diseases")
		return
	}
	if diseases == nil {
		diseases = []*models.Disease{}
	}
	response.Success(c, http.StatusOK, "Diseases fetched", diseases)
}

func (h *DiseaseHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid disease ID")
		return
	}

	disease, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Disease not found")
		return
	}

	treatments, _ := h.repo.GetTreatments(c.Request.Context(), id)

	response.Success(c, http.StatusOK, "Disease details", gin.H{
		"disease":    disease,
		"treatments": treatments,
	})
}
