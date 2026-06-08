package repository

import (
	"context"
	"fmt"
	"shorter-url/internal/domain"

	"github.com/pashagolub/pgxmock/v5"
)

type userRepository struct {
	db pgxmock.PgxPoolIface
}

func NewUserRepository(db pgxmock.PgxPoolIface) domain.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (u *userRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	query := `INSERT INTO users(email, password_hash)VALUES($1, $2) RETURNING id, email, password_hash, is_verified, status, created_at`

	err := u.db.QueryRow(ctx, query, user.Email, user.PasswordHash).Scan(&user.Id, &user.Email, &user.PasswordHash, &user.IsVerified, &user.Status, &user.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("something wrong when create new data: %w", err)
	}

	return user, nil
}

func (u *userRepository) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	query := `
		UPDATE users 
		SET email = $1, password_hash = $2, is_verified = $3, status = $4
		WHERE id = $5
		RETURNING id, email, password_hash, is_verified, status, created_at`

	err := u.db.QueryRow(ctx, query, user.Email, user.PasswordHash, user.IsVerified, user.Status, user.Id).
		Scan(&user.Id, &user.Email, &user.PasswordHash, &user.IsVerified, &user.Status, &user.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("something wrong when update data: %w", err)
	}

	return user, nil
}

func (u *userRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = $1`

	commandTag, err := u.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("something wrong when delete data : %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("there is no data deleted, user with ID %d not found", id)
	}
	return nil
}

func (u *userRepository) FindById(ctx context.Context, id int64) (*domain.User, error) {
	query := `SELECT id, email, password_hash, is_verified, status, created_at FROM users WHERE id = $1`
	user := &domain.User{}

	err := u.db.QueryRow(ctx, query, id).Scan(&user.Id, &user.Email, &user.PasswordHash, &user.IsVerified, &user.Status, &user.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("something error when find user by id : %w", err)
	}
	return user, nil
}

func (u *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, email, password_hash, is_verified, status, created_at FROM users WHERE email = $1`
	user := &domain.User{}

	err := u.db.QueryRow(ctx, query, email).Scan(&user.Id, &user.Email, &user.PasswordHash, &user.IsVerified, &user.Status, &user.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("something error when find user by email : %w", err)
	}
	return user, nil
}
