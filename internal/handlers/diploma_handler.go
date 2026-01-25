package handlers

import (
	"BlockCertify/internal/services"
	"BlockCertify/internal/utils"
	apperrors "BlockCertify/pkg/errors"
	"log"
	"net/http"
	"path/filepath"
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

	// Parse multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.RespondError(w, "Failed to parse form", err.Error(), apperrors.ErrInvalidFile, http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("diploma")
	if err != nil {
		utils.RespondError(w, "No file uploaded", err.Error(), apperrors.ErrInvalidFile, http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Save file
	filename := filepath.Base(header.Filename)
	filePath, err := h.fileManager.SaveUploadedFile(file, filename)
	if err != nil {
		utils.RespondError(w, "Failed to save file", err.Error(), apperrors.ErrInvalidFile, http.StatusInternalServerError)
		return
	}
	defer func(fileManager *utils.FileManager, filepath string) {
		err := fileManager.DeleteFile(filepath)
		if err != nil {

		}
	}(h.fileManager, filePath)

	// Hash file
	log.Println("Hashing diploma...")
	diplomaHash, err := utils.HashFile(filePath)
	if err != nil {
		utils.RespondError(w, "Failed to hash file", err.Error(), apperrors.ErrHashingFailed, http.StatusInternalServerError)
		return
	}

	// Upload diploma
	response, err := h.service.Upload(filePath, diplomaHash)
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

	// Parse multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.RespondError(w, "Failed to parse form", err.Error(), apperrors.ErrInvalidFile, http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("diploma")
	if err != nil {
		utils.RespondError(w, "No file uploaded", err.Error(), apperrors.ErrInvalidFile, http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Save file
	filename := filepath.Base(header.Filename)
	filePath, err := h.fileManager.SaveUploadedFile(file, filename)
	if err != nil {
		utils.RespondError(w, "Failed to save file", err.Error(), apperrors.ErrInvalidFile, http.StatusInternalServerError)
		return
	}
	defer h.fileManager.DeleteFile(filePath)

	// Hash file
	diplomaHash, err := utils.HashFile(filePath)
	if err != nil {
		utils.RespondError(w, "Failed to hash file", err.Error(), apperrors.ErrHashingFailed, http.StatusInternalServerError)
		return
	}

	// Verify diploma
	response, err := h.service.Verify(diplomaHash)
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
