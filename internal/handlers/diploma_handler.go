package handlers

import (
	"BlockCertify/internal/dto"
	"BlockCertify/internal/helper"
	"BlockCertify/internal/services"
	"BlockCertify/internal/utils"
	apperrors "BlockCertify/pkg/errors"
	"encoding/json"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

type DiplomaHandler struct {
	service     services.DiplomaService
	fileManager *utils.FileManager
}

func NewDiplomaHandler(service services.DiplomaService) *DiplomaHandler {
	return &DiplomaHandler{
		service:     service,
		fileManager: utils.NewFileManager("uploads"),
	}
}

func (h *DiplomaHandler) Upload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondError(w, "Method not allowed", "", "", http.StatusMethodNotAllowed)
		return
	}

	contentType := r.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "multipart/form-data") {
		utils.RespondError(w, "Invalid Content type", "multipart/form-data required", apperrors.ErrInvalidRequest, http.StatusBadRequest)
		return
	}

	reader, err := r.MultipartReader()
	if err != nil {
		utils.RespondError(w, "Invalid multipart request", err.Error(), apperrors.ErrInvalidFile, http.StatusBadRequest)
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
			utils.RespondError(w, "Failed to read multipart data", err.Error(), apperrors.ErrInvalidFile, http.StatusBadRequest)
			return
		}

		switch part.FormName() {

		case "diploma":
			if !strings.HasSuffix(strings.ToLower(part.FileName()), ".pdf") {
				utils.RespondError(w, "Invalid file type", "Only PDF files are allowed", apperrors.ErrInvalidFile, http.StatusBadRequest)
				return
			}

			filePath, err = h.fileManager.SaveUploadedFile(part, filepath.Base(part.FileName()))
			if err != nil {
				utils.RespondError(w, "Failed to save file", err.Error(), apperrors.ErrInvalidFile, http.StatusInternalServerError)
				return
			}
			defer h.fileManager.DeleteFile(filePath)

			log.Println("Hashing diploma...")
			diplomaHash, err = utils.HashFile(filePath)
			if err != nil {
				utils.RespondError(w, "Failed to hash file", err.Error(), apperrors.ErrHashingFailed, http.StatusInternalServerError)
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
		utils.RespondError(w, "Invalid metadata", err.Error(), apperrors.ErrInvalidRequest, http.StatusBadRequest)
		return
	}

	// Upload diploma
	response, err := h.service.Upload(filePath, diplomaHash, reqMeta)
	if err != nil {
		appErr, ok := err.(*apperrors.AppError)
		if ok {
			errDetails := ""
			if appErr.Err != nil {
				errDetails = appErr.Err.Error()
			}
			utils.RespondError(w, appErr.Message, errDetails, appErr.Code, http.StatusInternalServerError)
		} else {
			utils.RespondError(w, "Failed to process diploma", err.Error(), "", http.StatusInternalServerError)
		}
		return
	}

	utils.RespondJSON(w, response, http.StatusOK)
}

func (h *DiplomaHandler) Verify(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondError(w, "Method not allowed", "", "", http.StatusMethodNotAllowed)
		return
	}

	var req dto.VerifyDiplomaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, "Invalid request", err.Error(), apperrors.ErrInvalidRequest, http.StatusBadRequest)
		return
	}

	// Verify diploma
	response, err := h.service.Verify(req)
	if err != nil {
		appErr, ok := err.(*apperrors.AppError)
		if ok {
			utils.RespondError(w, appErr.Message, appErr.Err.Error(), appErr.Code, http.StatusInternalServerError)
		} else {
			utils.RespondError(w, "Failed to verify diploma", err.Error(), "", http.StatusInternalServerError)
		}
		return
	}

	utils.RespondJSON(w, response, http.StatusOK)
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
