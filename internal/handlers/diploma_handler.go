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
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"log/slog"

	"github.com/gin-gonic/gin"
)

// maxUploadBytes caps diploma upload requests; every uploaded byte is paid for
// on Arweave, so oversized files are rejected before any processing.
const maxUploadBytes = 25 << 20 // 25 MB

// parsedDiplomaUpload is the result of reading a multipart diploma request.
type parsedDiplomaUpload struct {
	FilePath string
	Hash     string
	Meta     dto.DiplomaMetadataRequest
	Cleanup  func()
}

// parseDiplomaUpload reads the multipart body (PDF + metadata fields), saves the
// file to a temp location and hashes it. On failure it writes the HTTP error
// response itself and returns ok=false. Callers must defer result.Cleanup().
func parseDiplomaUpload(c *gin.Context) (parsedDiplomaUpload, bool) {

	result := parsedDiplomaUpload{Cleanup: func() {}}
	fileManager := utils.NewFileManager("/tmp/uploads")

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxUploadBytes)

	contentType := c.GetHeader("Content-Type")
	if !strings.HasPrefix(contentType, "multipart/form-data") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "multipart/form-data required"})
		return result, false
	}

	reader, err := c.Request.MultipartReader()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid multipart request", "details": err.Error()})
		return result, false
	}

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			status := http.StatusBadRequest
			if strings.Contains(err.Error(), "request body too large") {
				status = http.StatusRequestEntityTooLarge
			}
			c.JSON(status, gin.H{"error": "Failed to read multipart data", "details": err.Error()})
			return result, false
		}

		switch part.FormName() {
		case "diploma":
			if !strings.HasSuffix(strings.ToLower(part.FileName()), ".pdf") {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Only PDF files are allowed"})
				return result, false
			}
			filePath, err := fileManager.SaveUploadedFile(part, filepath.Base(part.FileName()))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file", "details": err.Error()})
				return result, false
			}
			result.FilePath = filePath
			result.Cleanup = func() { fileManager.DeleteFile(filePath) }

			result.Hash, err = utils.HashFile(filePath)
			if err != nil {
				result.Cleanup()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash file", "details": err.Error()})
				return result, false
			}
		case "firstName":
			result.Meta.FirstName = readPartValue(part)
		case "lastName":
			result.Meta.LastName = readPartValue(part)
		case "email":
			result.Meta.Email = readPartValue(part)
		case "university":
			result.Meta.University = readPartValue(part)
		case "faculty":
			result.Meta.Faculty = readPartValue(part)
		case "department":
			result.Meta.Department = readPartValue(part)
		case "graduationYear":
			result.Meta.GraduationYear = helper.AtoiSafe(readPartValue(part))
		case "studentNumber":
			result.Meta.StudentNumber = readPartValue(part)
		case "nationality":
			result.Meta.Nationality = readPartValue(part)
		}
	}

	if result.FilePath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "diploma file is required"})
		return result, false
	}

	if err := validateUploadMetadata(result.Meta); err != nil {
		result.Cleanup()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid metadata", "details": err.Error()})
		return result, false
	}

	return result, true
}

// writeServiceError maps service errors onto a JSON error response.
func writeServiceError(c *gin.Context, err error) {
	var appErr *apperrors.AppError
	if errors.As(err, &appErr) {
		errDetails := ""
		if appErr.Err != nil {
			errDetails = appErr.Err.Error()
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": appErr.Message, "details": errDetails, "code": appErr.Code})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}

// Upload handles the full diploma issuance flow in a single request.
// The backend uses its own platform keys (Arweave + Polygon) — no MetaMask required.
func Upload(c *gin.Context) {

	diplomaService := DiplomaService.New(config.DB)
	diplomaService.UseBlockchainService(config.Blockchain)
	diplomaService.UseArweaveService(config.Arweave)

	parsed, ok := parseDiplomaUpload(c)
	if !ok {
		return
	}
	defer parsed.Cleanup()

	response, err := diplomaService.Upload(parsed.FilePath, parsed.Hash, parsed.Meta)
	if err != nil {
		writeServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// PrepareUpload handles the first phase of diploma issuance.
// It receives the PDF + metadata, uploads the file to Arweave, and returns
// the diploma hash and Arweave tx ID for the frontend to sign on Polygon via MetaMask.
func PrepareUpload(c *gin.Context) {

	diplomaService := DiplomaService.New(config.DB)
	diplomaService.UseBlockchainService(config.Blockchain)
	diplomaService.UseArweaveService(config.Arweave)

	parsed, ok := parseDiplomaUpload(c)
	if !ok {
		return
	}
	defer parsed.Cleanup()

	response, err := diplomaService.PrepareUpload(parsed.FilePath, parsed.Hash, parsed.Meta)
	if err != nil {
		writeServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// ConfirmUpload handles the second phase of diploma issuance.
// It receives the Polygon tx hash (signed by MetaMask), verifies the diploma
// on-chain, and saves the record to the DB.
func ConfirmUpload(c *gin.Context) {

	diplomaService := DiplomaService.New(config.DB)
	diplomaService.UseBlockchainService(config.Blockchain)

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
		writeServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// Verify checks a diploma against the DB and the Polygon contract.
// verified=false with HTTP 200 means the diploma could not be validated;
// HTTP 502 means verification infrastructure is unavailable — the two cases
// must not be conflated.
func Verify(c *gin.Context) {

	diplomaService := DiplomaService.New(config.DB)
	diplomaService.UseBlockchainService(config.Blockchain)

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
		slog.Error("diploma verification failed", "diplomaID", req.DiplomaID, "err", err)
		c.JSON(http.StatusBadGateway, gin.H{
			"error":     "Verification is temporarily unavailable",
			"diplomaID": req.DiplomaID,
		})
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
		return
	}

	slog.Info("stream diploma request", "publicID", publicID)

	arweaveUrl := diplomaService.GetArweaveUrlByDiplomaID(publicID)
	if arweaveUrl == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Diploma not found",
		})
		return
	}

	resp, err := http.Get(arweaveUrl)
	if err != nil {
		slog.Error("Failed to fetch diploma", "err", err)
		c.JSON(http.StatusBadGateway, gin.H{
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
		return
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
