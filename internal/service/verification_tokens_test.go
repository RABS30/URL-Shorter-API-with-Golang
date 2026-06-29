package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"shorter-url/internal/domain"
	"shorter-url/internal/helper"
	"shorter-url/internal/repository"
	"shorter-url/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_SendVerificationMail_Pass(t *testing.T) {
	mockRepo := new(repository.MockVerificationTokenRepository)
	mockUsers := new(repository.MockUserRepository)
	mockEmail := new(helper.MockEmailSender)

	ctx := context.Background()
	email := "test@mail.com"
	user := &domain.User{Id: 1, Email: email, IsVerified: false}

	mockUsers.On("FindByEmail", ctx, email).Return(user, nil)
	mockRepo.On("Create", ctx, mock.Anything).Return(&domain.VerificationToken{}, nil)
	mockEmail.On("SendEmailWithHTML", ctx, email, mock.Anything, "verification_account_mail.html").Return(nil)

	s := service.NewVerificationTokenService(mockRepo, mockUsers, mockEmail, "http://localhost:8080")
	err := s.SendVerificationMail(ctx, email)

	assert.NoError(t, err)
	mockUsers.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	mockEmail.AssertExpectations(t)
}

func Test_SendVerificationMail_Fail_UserNotFound(t *testing.T) {
	mockRepo := new(repository.MockVerificationTokenRepository)
	mockUsers := new(repository.MockUserRepository)
	mockEmail := new(helper.MockEmailSender)

	ctx := context.Background()
	email := "unknown@mail.com"

	mockUsers.On("FindByEmail", ctx, email).Return(nil, errors.New("user not found"))

	s := service.NewVerificationTokenService(mockRepo, mockUsers, mockEmail, "http://localhost:8080")
	err := s.SendVerificationMail(ctx, email)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to find user by email")
	mockUsers.AssertExpectations(t)
}

func Test_SendVerificationMail_Fail_UserVerified(t *testing.T) {
	mockRepo := new(repository.MockVerificationTokenRepository)
	mockUsers := new(repository.MockUserRepository)
	mockEmail := new(helper.MockEmailSender)

	ctx := context.Background()
	email := "verified@mail.com"
	user := &domain.User{Id: 1, Email: email, IsVerified: true}

	mockUsers.On("FindByEmail", ctx, email).Return(user, nil)

	s := service.NewVerificationTokenService(mockRepo, mockUsers, mockEmail, "http://localhost:8080")
	err := s.SendVerificationMail(ctx, email)

	assert.Error(t, err)
	assert.Equal(t, "user is verified", err.Error())
	mockUsers.AssertExpectations(t)
}

func Test_SendVerificationMail_Fail_CreateTokenError(t *testing.T) {
	mockRepo := new(repository.MockVerificationTokenRepository)
	mockUsers := new(repository.MockUserRepository)
	mockEmail := new(helper.MockEmailSender)

	ctx := context.Background()
	email := "test@mail.com"
	user := &domain.User{Id: 1, Email: email, IsVerified: false}

	mockUsers.On("FindByEmail", ctx, email).Return(user, nil)
	mockRepo.On("Create", ctx, mock.Anything).Return(nil, errors.New("db error"))

	s := service.NewVerificationTokenService(mockRepo, mockUsers, mockEmail, "http://localhost:8080")
	err := s.SendVerificationMail(ctx, email)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create verification token")
	mockUsers.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func Test_SendVerificationMail_Fail_SendEmailError(t *testing.T) {
	mockRepo := new(repository.MockVerificationTokenRepository)
	mockUsers := new(repository.MockUserRepository)
	mockEmail := new(helper.MockEmailSender)

	ctx := context.Background()
	email := "test@mail.com"
	user := &domain.User{Id: 1, Email: email, IsVerified: false}

	mockUsers.On("FindByEmail", ctx, email).Return(user, nil)
	mockRepo.On("Create", ctx, mock.Anything).Return(&domain.VerificationToken{}, nil)
	mockEmail.On("SendEmailWithHTML", ctx, email, mock.Anything, "verification_account_mail.html").Return(errors.New("smtp gateway timeout"))

	s := service.NewVerificationTokenService(mockRepo, mockUsers, mockEmail, "http://localhost:8080")
	err := s.SendVerificationMail(ctx, email)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to send verification email")
	mockUsers.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	mockEmail.AssertExpectations(t)
}

func Test_VerificationAccount_Pass(t *testing.T) {
	mockRepo := new(repository.MockVerificationTokenRepository)
	mockUsers := new(repository.MockUserRepository)
	mockEmail := new(helper.MockEmailSender)

	ctx := context.Background()
	token := "valid-token"
	verificationToken := &domain.VerificationToken{
		UserId:    1,
		Token:     token,
		ExpiredAt: time.Now().Add(time.Minute * 10),
	}
	user := &domain.User{Id: 1, IsVerified: false}

	mockRepo.On("FindByToken", ctx, token).Return(verificationToken, nil)
	mockUsers.On("FindById", ctx, int64(1)).Return(user, nil)
	mockUsers.On("UpdateVerified", ctx, int64(1), true).Return(nil)
	mockRepo.On("DeleteByUserId", ctx, int64(1)).Return(nil)

	s := service.NewVerificationTokenService(mockRepo, mockUsers, mockEmail, "http://localhost:8080")
	err := s.VerificationAccount(ctx, token)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockUsers.AssertExpectations(t)
}

func Test_VerificationAccount_Fail_TokenNotFound(t *testing.T) {
	mockRepo := new(repository.MockVerificationTokenRepository)
	mockUsers := new(repository.MockUserRepository)
	mockEmail := new(helper.MockEmailSender)

	ctx := context.Background()
	token := "invalid-token"

	mockRepo.On("FindByToken", ctx, token).Return(nil, errors.New("token not found"))

	s := service.NewVerificationTokenService(mockRepo, mockUsers, mockEmail, "http://localhost:8080")
	err := s.VerificationAccount(ctx, token)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to find verification token")
	mockRepo.AssertExpectations(t)
}

func Test_VerificationAccount_Fail_TokenExpired(t *testing.T) {
	mockRepo := new(repository.MockVerificationTokenRepository)
	mockUsers := new(repository.MockUserRepository)
	mockEmail := new(helper.MockEmailSender)

	ctx := context.Background()
	token := "expired-token"
	verificationToken := &domain.VerificationToken{
		UserId:    1,
		Token:     token,
		ExpiredAt: time.Now().Add(-time.Minute * 10),
	}

	mockRepo.On("FindByToken", ctx, token).Return(verificationToken, nil)

	s := service.NewVerificationTokenService(mockRepo, mockUsers, mockEmail, "http://localhost:8080")
	err := s.VerificationAccount(ctx, token)

	assert.Error(t, err)
	assert.Equal(t, "token has expired", err.Error())
	mockRepo.AssertExpectations(t)
}

func Test_VerificationAccount_Fail_UserNotFound(t *testing.T) {
	mockRepo := new(repository.MockVerificationTokenRepository)
	mockUsers := new(repository.MockUserRepository)
	mockEmail := new(helper.MockEmailSender)

	ctx := context.Background()
	token := "valid-token"
	verificationToken := &domain.VerificationToken{
		UserId:    1,
		Token:     token,
		ExpiredAt: time.Now().Add(time.Minute * 10),
	}

	mockRepo.On("FindByToken", ctx, token).Return(verificationToken, nil)
	mockUsers.On("FindById", ctx, int64(1)).Return(nil, errors.New("user not found"))

	s := service.NewVerificationTokenService(mockRepo, mockUsers, mockEmail, "http://localhost:8080")
	err := s.VerificationAccount(ctx, token)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to find user by ID")
	mockRepo.AssertExpectations(t)
	mockUsers.AssertExpectations(t)
}

func Test_VerificationAccount_Fail_UserAlreadyVerified(t *testing.T) {
	mockRepo := new(repository.MockVerificationTokenRepository)
	mockUsers := new(repository.MockUserRepository)
	mockEmail := new(helper.MockEmailSender)

	ctx := context.Background()
	token := "valid-token"
	verificationToken := &domain.VerificationToken{
		UserId:    1,
		Token:     token,
		ExpiredAt: time.Now().Add(time.Minute * 10),
	}
	user := &domain.User{Id: 1, IsVerified: true}

	mockRepo.On("FindByToken", ctx, token).Return(verificationToken, nil)
	mockUsers.On("FindById", ctx, int64(1)).Return(user, nil)

	s := service.NewVerificationTokenService(mockRepo, mockUsers, mockEmail, "http://localhost:8080")
	err := s.VerificationAccount(ctx, token)

	assert.Error(t, err)
	assert.Equal(t, "user is already verified", err.Error())
	mockRepo.AssertExpectations(t)
	mockUsers.AssertExpectations(t)
}

func Test_VerificationAccount_Fail_UpdateVerifiedError(t *testing.T) {
	mockRepo := new(repository.MockVerificationTokenRepository)
	mockUsers := new(repository.MockUserRepository)
	mockEmail := new(helper.MockEmailSender)

	ctx := context.Background()
	token := "valid-token"
	verificationToken := &domain.VerificationToken{
		UserId:    1,
		Token:     token,
		ExpiredAt: time.Now().Add(time.Minute * 10),
	}
	user := &domain.User{Id: 1, IsVerified: false}

	mockRepo.On("FindByToken", ctx, token).Return(verificationToken, nil)
	mockUsers.On("FindById", ctx, int64(1)).Return(user, nil)
	mockUsers.On("UpdateVerified", ctx, int64(1), true).Return(errors.New("db error"))

	s := service.NewVerificationTokenService(mockRepo, mockUsers, mockEmail, "http://localhost:8080")
	err := s.VerificationAccount(ctx, token)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update user verified status")
	mockRepo.AssertExpectations(t)
	mockUsers.AssertExpectations(t)
}

func Test_VerificationAccount_Fail_DeleteTokenError(t *testing.T) {
	mockRepo := new(repository.MockVerificationTokenRepository)
	mockUsers := new(repository.MockUserRepository)
	mockEmail := new(helper.MockEmailSender)

	ctx := context.Background()
	token := "valid-token"
	verificationToken := &domain.VerificationToken{
		UserId:    1,
		Token:     token,
		ExpiredAt: time.Now().Add(time.Minute * 10),
	}
	user := &domain.User{Id: 1, IsVerified: false}

	mockRepo.On("FindByToken", ctx, token).Return(verificationToken, nil)
	mockUsers.On("FindById", ctx, int64(1)).Return(user, nil)
	mockUsers.On("UpdateVerified", ctx, int64(1), true).Return(nil)
	mockRepo.On("DeleteByUserId", ctx, int64(1)).Return(errors.New("db error"))

	s := service.NewVerificationTokenService(mockRepo, mockUsers, mockEmail, "http://localhost:8080")
	err := s.VerificationAccount(ctx, token)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to delete verification token after use")
	mockRepo.AssertExpectations(t)
	mockUsers.AssertExpectations(t)
}
