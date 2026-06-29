package service

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockVerificationTokenService struct {
	mock.Mock
}

func (m *MockVerificationTokenService) SendVerificationMail(ctx context.Context, email string) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}

func (m *MockVerificationTokenService) VerificationAccount(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}
