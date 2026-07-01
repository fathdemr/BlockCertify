package handlers

import (
	"BlockCertify/internal/config"
	"BlockCertify/internal/dto"
	"BlockCertify/internal/helper"
	apperrors "BlockCertify/internal/pkg/errors"
	"BlockCertify/internal/services/DiplomaService"
	"BlockCertify/internal/utils"
	"errors"
	"io"
	"log"
	"log/slog"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Upload handles the full diploma issuance flow in a single request.
// The backend uses its own platform keys (Arweave + Polygon) — no MetaMask required.
func Upload(c *gin.Context) {

	fileManager := utils.NewFileManager("/tmp/uploads")
	diplomaService := DiplomaService.New(config.DB)
	diplomaService.UseBlockchainService(config.Blockchain)
	diplomaService.UseArweaveService(config.Arweave)

	contentType := c.GetHeader("Content-Type")
	if !strings.HasPrefix(contentType, "multipart/form-data") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "multipart/form-data required"})
		return
	}

	reader, err := c.Request.MultipartReader()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid multipart request", "details": err.Error()})
		return
	}

	var (
		filePath    string
		diplomaHash string
		reqMeta     = dto.DiplomaMetadataRequest{}
	)

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read multipart data", "details": err.Error()})
			return
		}

		switch part.FormName() {
		case "diploma":
			if !strings.HasSuffix(strings.ToLower(part.FileName()), ".pdf") {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Only PDF files are allowed"})
				return
			}
			filePath, err = fileManager.SaveUploadedFile(part, filepath.Base(part.FileName()))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file", "details": err.Error()})
				return
			}
			defer fileManager.DeleteFile(filePath)
			diplomaHash, err = utils.HashFile(filePath)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash file", "details": err.Error()})
				return
			}
		case "firstName":
			reqMeta.FirstName = readPartValue(part)
		case "lastName":
			reqMeta.LastName = readPartValue(part)
		case "email":
			reqMeta.Email = readPartValue(part)
		case "university":
			reqMeta.University = readPartValue(part)
		case "faculty":
			reqMeta.Faculty = readPartValue(part)
		case "department":
			reqMeta.Department = readPartValue(part)
		case "graduationYear":
			reqMeta.GraduationYear = helper.AtoiSafe(readPartValue(part))
		case "studentNumber":
			reqMeta.StudentNumber = readPartValue(part)
		case "nationality":
			reqMeta.Nationality = readPartValue(part)
		}
	}

	if err := validateUploadMetadata(reqMeta); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid metadata", "details": err.Error()})
		return
	}

	response, err := diplomaService.Upload(filePath, diplomaHash, reqMeta)
	if err != nil {
		appErr, ok := err.(*apperrors.AppError)
		if ok {
			errDetails := ""
			if appErr.Err != nil {
				errDetails = appErr.Err.Error()
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": appErr.Message, "details": errDetails, "code": appErr.Code})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, response)
}

// PrepareUpload handles the first phase of diploma issuance.
// It receives the PDF + metadata, uploads the file to Arweave, and returns
// the diploma hash and Arweave tx ID for the frontend to sign on Polygon via MetaMask.
func PrepareUpload(c *gin.Context) {

	fileManager := utils.NewFileManager("/tmp/uploads")
	diplomaService := DiplomaService.New(config.DB)
	diplomaService.UseBlockchainService(config.Blockchain)
	diplomaService.UseArweaveService(config.Arweave)

	contentType := c.GetHeader("Content-Type")
	if !strings.HasPrefix(contentType, "multipart/form-data") {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid content type",
			"details": "multipart/form-data required",
		})
		return
	}

	reader, err := c.Request.MultipartReader()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid multipart request",
			"details": err.Error(),
		})
		return
	}

	var (
		filePath    string
		diplomaHash string
		reqMeta     = dto.DiplomaMetadataRequest{}
	)

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Failed to read multipart data",
				"details": err.Error(),
			})
			return
		}

		switch part.FormName() {

		case "diploma":
			if !strings.HasSuffix(strings.ToLower(part.FileName()), ".pdf") {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   "Invalid file type",
					"details": "Only PDF files are allowed",
				})
				return
			}

			filePath, err = fileManager.SaveUploadedFile(part, filepath.Base(part.FileName()))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Failed to save file",
					"details": err.Error(),
				})
				return
			}
			defer fileManager.DeleteFile(filePath)

			log.Println("Hashing diploma...")
			diplomaHash, err = utils.HashFile(filePath)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Failed to hash file",
					"details": err.Error(),
				})
				return
			}

		case "firstName":
			reqMeta.FirstName = readPartValue(part)
		case "lastName":
			reqMeta.LastName = readPartValue(part)
		case "email":
			reqMeta.Email = readPartValue(part)
		case "university":
			reqMeta.University = readPartValue(part)
		case "faculty":
			reqMeta.Faculty = readPartValue(part)
		case "department":
			reqMeta.Department = readPartValue(part)
		case "graduationYear":
			reqMeta.GraduationYear = helper.AtoiSafe(readPartValue(part))
		case "studentNumber":
			reqMeta.StudentNumber = readPartValue(part)
		case "nationality":
			reqMeta.Nationality = readPartValue(part)
		}
	}

	if err := validateUploadMetadata(reqMeta); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid metadata",
			"details": err.Error(),
		})
		return
	}

	response, err := diplomaService.PrepareUpload(filePath, diplomaHash, reqMeta)
	if err != nil {
		appErr, ok := err.(*apperrors.AppError)
		if ok {
			errDetails := ""
			if appErr.Err != nil {
				errDetails = appErr.Err.Error()
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   appErr.Message,
				"details": errDetails,
				"code":    appErr.Code,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, response)
}

// ConfirmUpload handles the second phase of diploma issuance.
// It receives the Polygon tx hash (signed by MetaMask) and saves the diploma record to the DB.
func ConfirmUpload(c *gin.Context) {

	diplomaService := DiplomaService.New(config.DB)

	var req dto.ConfirmUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	response, err := diplomaService.ConfirmUpload(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

func Verify(c *gin.Context) {

	diplomaService := DiplomaService.New(config.DB)

	var req dto.VerifyDiplomaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	response, err := diplomaService.Verify(req)
	if err != nil {
		c.JSON(http.StatusOK, dto.VerifyResponse{Verified: false, DiplomaID: req.DiplomaID})
		return
	}

	c.JSON(http.StatusOK, response)
}

func GetDiplomaById(c *gin.Context) {

	diplomaService := DiplomaService.New(config.DB)

	publicID := strings.TrimSpace(c.Param("diplomaId"))

	if publicID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Diploma Id",
		})
	}

	slog.Info("stream diploma request", "publicID", publicID)

	arweaveUrl := diplomaService.GetArweaveUrlByDiplomaID(publicID)

	resp, err := http.Get(arweaveUrl)
	if err != nil {
		slog.Error("Failed to fetch diploma", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to fetch diploma",
		})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("arweave returned non-200 status code", "statusCode", resp.StatusCode)
		c.JSON(resp.StatusCode, gin.H{
			"error": "file not found",
		})
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "inline; filename=diploma.pdf")

	_, err = io.Copy(c.Writer, resp.Body)
	if err != nil {
		slog.Error("Failed to write response", "err", err)
	}
}

func GetDiplomaRecords(c *gin.Context) {

	email, err := extractEmailFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	diplomaService := DiplomaService.New(config.DB)

	ch := diplomaService.GetDiplomasByEmail(email)

	var records []map[string]interface{}
	for diploma := range ch {
		records = append(records, map[string]interface{}{
			"id":              diploma.ID,
			"diploma_no":      diploma.PublicID,
			"tx_hash":         diploma.PolygonTxID,
			"arweave_tx":      diploma.ArweaveTxID,
			"status":          "approved",
			"created_at":      diploma.CreatedAt,
			"first_name":      diploma.MetaData.FirstName,
			"last_name":       diploma.MetaData.LastName,
			"email":           diploma.MetaData.Email,
			"student_no":      diploma.MetaData.StudentNumber,
			"university":      diploma.MetaData.University,
			"faculty":         diploma.MetaData.Faculty,
			"department":      diploma.MetaData.Department,
			"graduation_year": diploma.MetaData.GraduationYear,
		})
	}

	if records == nil {
		records = []map[string]interface{}{}
	}

	c.JSON(http.StatusOK, records)
}

func validateUploadMetadata(meta dto.DiplomaMetadataRequest) error {
	if strings.TrimSpace(meta.FirstName) == "" {
		return errors.New("firstName is required")
	}
	if strings.TrimSpace(meta.LastName) == "" {
		return errors.New("lastName is required")
	}
	if strings.TrimSpace(meta.Email) == "" {
		return errors.New("email is required")
	}
	if strings.TrimSpace(meta.University) == "" {
		return errors.New("university is required")
	}
	if strings.TrimSpace(meta.Department) == "" {
		return errors.New("department is required")
	}
	if meta.GraduationYear < 1950 || meta.GraduationYear > time.Now().Year()+1 {
		return errors.New("invalid graduation year")
	}
	return nil
}

func readPartValue(part *multipart.Part) string {
	b, _ := io.ReadAll(part)
	return strings.TrimSpace(string(b))
}
