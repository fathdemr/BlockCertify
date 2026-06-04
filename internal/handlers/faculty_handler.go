package handlers

import (
	"BlockCertify/internal/config"
	"BlockCertify/internal/services/FacultyService"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
)

func GetFaculties(c *gin.Context) {
	facultyService := FacultyService.New(config.DB)

	// Accept university_id as query param (browser-friendly)
	// Falls back to JSON body for backward compatibility
	var universityID uuid.UUID

	if qid := c.Query("university_id"); qid != "" {
		parsed, err := uuid.FromString(qid)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid university_id"})
			return
		}
		universityID = parsed
	} else {
		var req struct {
			ID uuid.UUID `json:"id"`
		}
		if err := c.ShouldBindJSON(&req); err != nil || req.ID == uuid.Nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "university_id required"})
			return
		}
		universityID = req.ID
	}

	response, err := facultyService.GetAllFacultiesByID(universityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
