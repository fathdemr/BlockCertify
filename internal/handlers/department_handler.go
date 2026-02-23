package handlers

import (
	"BlockCertify/internal/dto"
	"BlockCertify/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DepartmentHandler struct {
	service services.DepartmentService
}

func NewDepartmentHandler(service services.DepartmentService) *DepartmentHandler {
	return &DepartmentHandler{service: service}
}

func (h *DepartmentHandler) GetDepartmentByID(c *gin.Context) {

	var req dto.DepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	response, err := h.service.GetDepartmentByID(req.FacultyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
