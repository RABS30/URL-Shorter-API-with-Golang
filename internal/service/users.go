package service

import (
	"context"
	"errors"
	"fmt"
	"shorter-url/internal/domain"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo      domain.UserRepository
	JwtSecret []byte
}

func NewUserService(repo domain.UserRepository, JwtSecret []byte) domain.AuthService {
	return &userService{
		repo:      repo,
		JwtSecret: JwtSecret,
	}
}

func (s *userService) Register(ctx context.Context, email string, password string) (*domain.User, error) {
	existingUser, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			// PERBAIKAN: Bersihkan spasi ganda sebelum %w rabs
			return nil, fmt.Errorf("failed to check existing email: %w", err)
		}
	}
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		// PERBAIKAN: Ubah koma menjadi titik dua dan ganti kata encrypt menjadi hash rabs
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	newUser := &domain.User{
		Email:        email,
		PasswordHash: string(hashedPassword),
	}

	result, err := s.repo.Create(ctx, newUser)
	if err != nil {
		// PERBAIKAN: Format lowercase dan ubah koma menjadi titik dua rabs
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return result, nil
}

func (s *userService) Login(ctx context.Context, email string, password string) (string, error) {
	invalidError := errors.New("invalid email or password")

	existingUser, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", invalidError
		}
		return "", fmt.Errorf("failed to find user by email: %w", err)
	}
	if existingUser == nil {
		return "", invalidError
	}

	err = bcrypt.CompareHashAndPassword([]byte(existingUser.PasswordHash), []byte(password))
	if err != nil {
		return "", invalidError
	}

	claims := jwt.MapClaims{
		"user_id":     existingUser.Id,
		"email":       existingUser.Email,
		"is_verified": existingUser.IsVerified,
		"exp":         time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.JwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to generate jwt token: %w", err)
	}

	return tokenString, nil
}
