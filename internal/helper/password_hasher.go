package helper

import (
	"context"
	"fmt"
	"shorter-url/internal/domain"

	"golang.org/x/crypto/bcrypt"
)

type bcryptHasher struct{}

func NewBcryptHasher() domain.PasswordHasher {
	return &bcryptHasher{}
}

func (b *bcryptHasher) Hash(ctx context.Context, password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}

	return string(hashedPassword), nil
}

func (b *bcryptHasher) Compare(ctx context.Context, password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
