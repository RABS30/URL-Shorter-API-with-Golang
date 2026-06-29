package handler

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"shorter-url/internal/service"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_RequestVerification_Pass(t *testing.T) {
	mockService := new(service.MockVerificationTokenService)
	h := NewVerificationTokenHandler(mockService)

	body := `{"email":"test@mail.com"}`
	req := httptest.NewRequest(http.MethodPost, "/verify", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()

	mockService.On("SendVerificationMail", mock.Anything, "test@mail.com").Return(nil)

	h.RequestVerification(rec, req, httprouter.Params{})

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "verification email has been sent")
	mockService.AssertExpectations(t)
}

func Test_RequestVerification_Fail_InvalidJSON(t *testing.T) {
	mockService := new(service.MockVerificationTokenService)
	h := NewVerificationTokenHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/verify", bytes.NewBufferString(`{"email":`))
	rec := httptest.NewRecorder()

	h.RequestVerification(rec, req, httprouter.Params{})

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "invalid request body")
}

func Test_RequestVerification_Fail_EmptyEmail(t *testing.T) {
	mockService := new(service.MockVerificationTokenService)
	h := NewVerificationTokenHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/verify", bytes.NewBufferString(`{"email":""}`))
	rec := httptest.NewRecorder()

	h.RequestVerification(rec, req, httprouter.Params{})

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "email is required")
}

func Test_RequestVerification_Fail_UserVerified(t *testing.T) {
	mockService := new(service.MockVerificationTokenService)
	h := NewVerificationTokenHandler(mockService)

	body := `{"email":"verified@mail.com"}`
	req := httptest.NewRequest(http.MethodPost, "/verify", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()

	mockService.On("SendVerificationMail", mock.Anything, "verified@mail.com").Return(errors.New("user is verified"))

	h.RequestVerification(rec, req, httprouter.Params{})

	assert.Equal(t, http.StatusConflict, rec.Code)
	assert.Contains(t, rec.Body.String(), "user is verified")
	mockService.AssertExpectations(t)
}

func Test_RequestVerification_Fail_ServiceError(t *testing.T) {
	mockService := new(service.MockVerificationTokenService)
	h := NewVerificationTokenHandler(mockService)

	body := `{"email":"test@mail.com"}`
	req := httptest.NewRequest(http.MethodPost, "/verify", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()

	mockService.On("SendVerificationMail", mock.Anything, "test@mail.com").Return(errors.New("db error"))

	h.RequestVerification(rec, req, httprouter.Params{})

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "failed to send verification email")
	mockService.AssertExpectations(t)
}

func Test_VerificationAccount_Pass(t *testing.T) {
	mockService := new(service.MockVerificationTokenService)
	h := NewVerificationTokenHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/verify-account?token=valid-token", nil)
	rec := httptest.NewRecorder()

	mockService.On("VerificationAccount", mock.Anything, "valid-token").Return(nil)

	h.VerificationAccount(rec, req, httprouter.Params{})

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "success")
	mockService.AssertExpectations(t)
}

func Test_VerificationAccount_Fail_TokenMissing(t *testing.T) {
	mockService := new(service.MockVerificationTokenService)
	h := NewVerificationTokenHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/verify-account", nil)
	rec := httptest.NewRecorder()

	h.VerificationAccount(rec, req, httprouter.Params{})

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "token not found")
}

func Test_VerificationAccount_Fail_TokenExpired(t *testing.T) {
	mockService := new(service.MockVerificationTokenService)
	h := NewVerificationTokenHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/verify-account?token=expired-token", nil)
	rec := httptest.NewRecorder()

	mockService.On("VerificationAccount", mock.Anything, "expired-token").Return(errors.New("token has expired"))

	h.VerificationAccount(rec, req, httprouter.Params{})

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "token has expired")
	mockService.AssertExpectations(t)
}

func Test_VerificationAccount_Fail_UserAlreadyVerified(t *testing.T) {
	mockService := new(service.MockVerificationTokenService)
	h := NewVerificationTokenHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/verify-account?token=verified-token", nil)
	rec := httptest.NewRecorder()

	mockService.On("VerificationAccount", mock.Anything, "verified-token").Return(errors.New("user is already verified"))

	h.VerificationAccount(rec, req, httprouter.Params{})

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "user is already verified")
	mockService.AssertExpectations(t)
}

func Test_VerificationAccount_Fail_ServiceError(t *testing.T) {
	mockService := new(service.MockVerificationTokenService)
	h := NewVerificationTokenHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/verify-account?token=some-token", nil)
	rec := httptest.NewRecorder()

	mockService.On("VerificationAccount", mock.Anything, "some-token").Return(errors.New("db error"))

	h.VerificationAccount(rec, req, httprouter.Params{})

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "failed to verified account")
	mockService.AssertExpectations(t)
}
