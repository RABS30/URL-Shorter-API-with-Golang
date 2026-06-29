package repository

import (
	"context"
	"shorter-url/internal/domain"

	"github.com/stretchr/testify/mock"
)

type MockVerificationTokenRepository struct {
	mock.Mock
}

func (m *MockVerificationTokenRepository) Create(ctx context.Context, verificationToken *domain.VerificationToken) (*domain.VerificationToken, error) {
	args := m.Called(ctx, verificationToken)
	var res *domain.VerificationToken
	if args.Get(0) != nil {
		res = args.Get(0).(*domain.VerificationToken)
	}
	return res, args.Error(1)
}

func (m *MockVerificationTokenRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockVerificationTokenRepository) FindByToken(ctx context.Context, token string) (*domain.VerificationToken, error) {
	args := m.Called(ctx, token)
	var res *domain.VerificationToken
	if args.Get(0) != nil {
		res = args.Get(0).(*domain.VerificationToken)
	}
	return res, args.Error(1)
}

func (m *MockVerificationTokenRepository) DeleteByUserId(ctx context.Context, userId int64) error {
	args := m.Called(ctx, userId)
	return args.Error(0)
}
