package repository

import (
	"context"
	"errors"
	"fmt"
	"shorter-url/internal/database"
	"shorter-url/internal/domain"

	"github.com/jackc/pgx/v5"
)

type passwordResetTokensRepository struct {
	db database.PgxDatabase
}

func NewPasswordResetTokensRepository(db database.PgxDatabase) domain.PasswordResetTokensRepository {
	return &passwordResetTokensRepository{
		db: db,
	}
}

func (r *passwordResetTokensRepository) Create(ctx context.Context, passwordResetToken *domain.PasswordResetTokens) (*domain.PasswordResetTokens, error) {
	query := `INSERT INTO password_reset_tokens (user_id, token, expired_at) VALUES ($1, $2, $3) RETURNING id, user_id, token, expired_at, created_at`

	err := r.db.QueryRow(ctx, query, passwordResetToken.UserId, passwordResetToken.Token, passwordResetToken.ExpiredAt).
		Scan(&passwordResetToken.Id, &passwordResetToken.UserId, &passwordResetToken.Token, &passwordResetToken.ExpiredAt, &passwordResetToken.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("insert password reset token: %w", err)
	}

	return passwordResetToken, nil
}

func (r *passwordResetTokensRepository) Delete(ctx context.Context, id int64) error {
	query := "DELETE FROM password_reset_tokens WHERE id = $1"

	commandTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete password reset token by id: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("delete password reset token by id: %w", domain.ErrNotFound)
	}

	return nil
}

func (r *passwordResetTokensRepository) DeleteByUserId(ctx context.Context, userId int64) error {
	query := "DELETE FROM password_reset_tokens WHERE user_id = $1"

	_, err := r.db.Exec(ctx, query, userId)
	if err != nil {
		return fmt.Errorf("delete password reset token by user id: %w", err)
	}
	return nil
}

func (r *passwordResetTokensRepository) FindByToken(ctx context.Context, token string) (*domain.PasswordResetTokens, error) {
	query := "SELECT id, user_id, token, expired_at, created_at FROM password_reset_tokens WHERE token = $1"
	var passwordResetToken domain.PasswordResetTokens

	err := r.db.QueryRow(ctx, query, token).
		Scan(&passwordResetToken.Id, &passwordResetToken.UserId, &passwordResetToken.Token, &passwordResetToken.ExpiredAt, &passwordResetToken.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = domain.ErrNotFound
		}
		return nil, fmt.Errorf("query password reset tokens by token: %w", err)
	}

	return &passwordResetToken, nil
}
