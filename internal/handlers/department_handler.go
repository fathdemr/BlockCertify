package handlers

import (
	"BlockCertify/internal/config"
	"BlockCertify/internal/services/DepartmentService"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
)

func GetDepartmentByID(c *gin.Context) {
	departmentService := DepartmentService.New(config.DB)

	// Accept faculty_id as query param (browser-friendly)
	// Falls back to JSON body for backward compatibility
	var facultyID uuid.UUID

	if qid := c.Query("faculty_id"); qid != "" {
		parsed, err := uuid.FromString(qid)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid faculty_id"})
			return
		}
		facultyID = parsed
	} else {
		var req struct {
			ID uuid.UUID `json:"id"`
		}
		if err := c.ShouldBindJSON(&req); err != nil || req.ID == uuid.Nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "faculty_id required"})
			return
		}
		facultyID = req.ID
	}

	response, err := departmentService.GetDepartmentByID(facultyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
