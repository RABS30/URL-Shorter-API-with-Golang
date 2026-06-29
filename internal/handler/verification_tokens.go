package handler

import (
	"encoding/json"
	"net/http"
	"shorter-url/internal/domain"
	"shorter-url/internal/helper"
	"shorter-url/internal/middleware"

	"github.com/julienschmidt/httprouter"
)

type verificationTokenHandler struct {
	verificationService domain.VerificationTokenService
}

func NewVerificationTokenHandler(service domain.VerificationTokenService) *verificationTokenHandler {
	return &verificationTokenHandler{
		verificationService: service,
	}
}

func (h *verificationTokenHandler) RequestVerification(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.BadResponse(w, http.StatusBadRequest, "invalid request body")

		if wrapper, ok := w.(*middleware.ResponseWriterWrapper); ok {
			wrapper.WriteError(err.Error())
		}

		return
	}

	if req.Email == "" {
		helper.BadResponse(w, http.StatusBadRequest, "email is required")

		if wrapper, ok := w.(*middleware.ResponseWriterWrapper); ok {
			wrapper.WriteError("email is required")
		}
		return
	}

	err := h.verificationService.SendVerificationMail(r.Context(), req.Email)
	if err != nil {
		if err.Error() == "user is verified" {
			helper.BadResponse(w, http.StatusConflict, err.Error())
		} else {
			helper.BadResponse(w, http.StatusInternalServerError, "failed to send verification email")
		}
		if wrapper, ok := w.(*middleware.ResponseWriterWrapper); ok {
			wrapper.WriteError(err.Error())
		}

		return
	}

	helper.GoodResponse(w, http.StatusOK, "verification email has been sent", "")
}

func (h *verificationTokenHandler) VerificationAccount(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	token := r.URL.Query().Get("token")
	if token == "" {
		helper.BadResponse(w, http.StatusBadRequest, "token not found")

		if wrapper, ok := w.(*middleware.ResponseWriterWrapper); ok {
			wrapper.WriteError("token not found")
		}
		return
	}

	ctx := r.Context()

	err := h.verificationService.VerificationAccount(ctx, token)
	if err != nil {
		if err.Error() == "token has expired" || err.Error() == "user is already verified" {
			helper.BadResponse(w, http.StatusBadRequest, err.Error())
		} else {
			helper.BadResponse(w, http.StatusInternalServerError, "failed to verified account")
		}

		if wrapper, ok := w.(*middleware.ResponseWriterWrapper); ok {
			wrapper.WriteError(err.Error())
		}

		return
	}

	helper.GoodResponse(w, http.StatusOK, "success", "")
}
