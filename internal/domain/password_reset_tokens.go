package domain

import (
	"context"
	"time"
)

type PasswordResetTokens struct {
	Id        int64     `json:"id"`
	UserId    int64     `json:"user_id"`
	Token     string    `json:"token"`
	ExpiredAt time.Time `json:"expired_at"`
	CreatedAt time.Time `json:"created_at"`
}

type PasswordResetTokensRepository interface {
	Create(ctx context.Context, passwordResetToken *PasswordResetTokens) (*PasswordResetTokens, error)
	Delete(ctx context.Context, id int64) error
	DeleteByUserId(ctx context.Context, userId int64) error
	FindByToken(ctx context.Context, token string) (*PasswordResetTokens, error)
}

type PasswordResetTokensService interface {
	RequestResetPassword(ctx context.Context, email string) error
	ExecuteResetPassword(ctx context.Context, token string, password1 string, password2 string) error
}
