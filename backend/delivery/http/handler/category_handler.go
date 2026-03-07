package handler

import (
	"net/http"
	"taskmanager/usecase"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	categoryUseCase usecase.CategoryUseCase
}

func NewCategoryHandler(categoryUseCase usecase.CategoryUseCase) *CategoryHandler {
	return &CategoryHandler{categoryUseCase: categoryUseCase}
}

type createCategoryReq struct {
	Name string `json:"name" binding:"required"`
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	userID := c.GetInt("user_id")

	var req createCategoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category, err := h.categoryUseCase.CreateCategory(c.Request.Context(), userID, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, category)
}

func (h *CategoryHandler) GetCategories(c *gin.Context) {
	userID := c.GetInt("user_id")

	categories, err := h.categoryUseCase.GetCategories(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}
