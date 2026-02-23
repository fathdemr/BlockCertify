package handlers

import (
	"BlockCertify/internal/dto"
	"BlockCertify/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FacultyHandler struct {
	service services.FacultyService
}

func NewFacultyHandler(service services.FacultyService) *FacultyHandler {
	return &FacultyHandler{
		service: service,
	}
}

func (h *FacultyHandler) GetFaculties(c *gin.Context) {

	var req dto.FacultiesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.GetFaculties(req.UniversityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
