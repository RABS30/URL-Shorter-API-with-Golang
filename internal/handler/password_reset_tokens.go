package handler

import (
	"encoding/json"
	"net/http"
	"shorter-url/internal/domain"
	"shorter-url/internal/helper"
	"shorter-url/internal/middleware"

	"github.com/julienschmidt/httprouter"
)

type inputEmail struct {
	Email string `json:"email"`
}

type newPassword struct {
	Password1 string `json:"password1"`
	Password2 string `json:"password2"`
}

type passwordResetTokensHandler struct {
	service domain.PasswordResetTokensService
}

func NewPasswordResetTokensHandler(service domain.PasswordResetTokensService) *passwordResetTokensHandler {
	return &passwordResetTokensHandler{
		service: service,
	}
}

func (h *passwordResetTokensHandler) ForgotPasswordHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var req inputEmail
	ctx := r.Context()

	request := json.NewDecoder(r.Body)
	request.DisallowUnknownFields()
	err := request.Decode(&req)
	if err != nil {
		helper.BadResponse(w, http.StatusBadRequest, "invalid request payload")

		if wrapper, ok := w.(*middleware.ResponseWriterWrapper); ok {
			wrapper.WriteError(err.Error())
		}
		return
	}

	err = h.service.RequestResetPassword(ctx, req.Email)
	if err != nil {
		helper.BadResponse(w, http.StatusInternalServerError, "unable to process request")

		if wrapper, ok := w.(*middleware.ResponseWriterWrapper); ok {
			wrapper.WriteError(err.Error())
		}
		return
	}

	helper.GoodResponse(w, http.StatusOK, "success", "")
}

func (h *passwordResetTokensHandler) ResetPasswordHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var req newPassword

	ctx := r.Context()
	token := r.URL.Query().Get("token")
	if token == "" {
		helper.BadResponse(w, http.StatusBadRequest, "token is required")

		if wrapper, ok := w.(*middleware.ResponseWriterWrapper); ok {
			wrapper.WriteError("token not found")
		}
		return
	}

	request := json.NewDecoder(r.Body)
	request.DisallowUnknownFields()
	err := request.Decode(&req)
	if err != nil {
		helper.BadResponse(w, http.StatusBadRequest, "invalid request payload")

		if wrapper, ok := w.(*middleware.ResponseWriterWrapper); ok {
			wrapper.WriteError(err.Error())
		}
		return
	}
	if req.Password1 != req.Password2 {
		helper.BadResponse(w, http.StatusBadRequest, "password1 and password2 do not match")

		if wrapper, ok := w.(*middleware.ResponseWriterWrapper); ok {
			wrapper.WriteError("password1 and password2 do not match")
		}

		return
	}

	err = h.service.ExecuteResetPassword(ctx, token, req.Password1, req.Password2)
	if err != nil {
		helper.BadResponse(w, http.StatusInternalServerError, "unable to process request")

		if wrapper, ok := w.(*middleware.ResponseWriterWrapper); ok {
			wrapper.WriteError(err.Error())
		}
		return
	}

	helper.GoodResponse(w, http.StatusOK, "success", "")
}
