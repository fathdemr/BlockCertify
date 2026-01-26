package handlers

import (
	"BlockCertify/internal/dto"
	"BlockCertify/internal/services"
	"BlockCertify/internal/utils"
	"encoding/json"
	"net/http"

	apperrors "BlockCertify/pkg/errors"
)

type UserHandler struct {
	service services.UserService
}

func NewUserHandler(service services.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		utils.RespondError(w, "Method not allowed", "", "", http.StatusMethodNotAllowed)
		return
	}

	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, "Invalid request body", err.Error(), apperrors.ErrInvalidRequest, http.StatusBadRequest)
		return
	}

	response, err := h.service.Login(req)
	if err != nil {
		appErr, ok := err.(*apperrors.AppError)
		if ok {
			details := ""
			if appErr.Err != nil {
				details = appErr.Err.Error()
			}
			utils.RespondError(w, appErr.Message, details, appErr.Code, http.StatusBadRequest)
			return
		}
		utils.RespondError(w, "Login failed", err.Error(), "", http.StatusInternalServerError)
		return
	}
	utils.RespondJSON(w, response, http.StatusOK)
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		utils.RespondError(w, "Method not allowed", "", "", http.StatusMethodNotAllowed)
	}

	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, "Invalid request body", err.Error(), apperrors.ErrInvalidRequest, http.StatusBadRequest)
	}

	if err := h.service.Register(req); err != nil {
		appErr, ok := err.(*apperrors.AppError)
		if ok {
			details := ""
			if appErr.Err != nil {
				details = appErr.Err.Error()
			}
			utils.RespondError(w, appErr.Message, details, appErr.Code, http.StatusBadRequest)
			return
		}
		utils.RespondError(w, "Registration failed", err.Error(), "", http.StatusInternalServerError)
		return
	}
	utils.RespondJSON(w, map[string]string{
		"message": "User regsitered succesfully",
	}, http.StatusCreated)

}
