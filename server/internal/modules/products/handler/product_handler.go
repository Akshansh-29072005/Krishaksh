package handler

import (
	"net/http"

	"github.com/aarcsx/krishaksh-backend/internal/core/response"
	"github.com/aarcsx/krishaksh-backend/internal/models"
	productRepo "github.com/aarcsx/krishaksh-backend/internal/modules/products/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProductHandler struct {
	repo productRepo.ProductRepository
}

func NewProductHandler(repo productRepo.ProductRepository) *ProductHandler {
	return &ProductHandler{repo: repo}
}

func (h *ProductHandler) GetAll(c *gin.Context) {
	cropType := c.Query("crop_type")
	products, err := h.repo.GetAll(c.Request.Context(), cropType, false)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch products")
		return
	}
	if products == nil {
		products = []*models.Product{}
	}
	response.Success(c, http.StatusOK, "Products fetched", products)
}

func (h *ProductHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid product ID")
		return
	}

	product, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Product not found")
		return
	}

	response.Success(c, http.StatusOK, "Product details", product)
}
