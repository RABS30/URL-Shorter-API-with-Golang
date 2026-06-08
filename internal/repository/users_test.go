package repository

import (
	"context"
	"fmt"
	"shorter-url/internal/domain"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock/v5"
	"github.com/stretchr/testify/assert"
)

func TestCreateUsersPass(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mockPool.Close()
	defer func() {
		err := mockPool.ExpectationsWereMet()
		assert.NoError(t, err)
	}()

	ctx := context.Background()

	inputData := &domain.User{
		Email:        "contoh@gmail.com",
		PasswordHash: "iniadalahhashrandom",
	}
	id := int64(9)
	isVerified := false
	status := "active"
	createdAt := time.Now()

	query := `^INSERT INTO users\(email, password_hash\)VALUES\(\$1, \$2\) RETURNING id, email, password_hash, is_verified, status, created_at$`

	mockRow := pgxmock.NewRows([]string{"1", "2", "3", "4", "5", "6"}).AddRow(
		id,
		inputData.Email,
		inputData.PasswordHash,
		isVerified,
		status,
		createdAt,
	)

	mockPool.ExpectQuery(query).WithArgs(inputData.Email, inputData.PasswordHash).WillReturnRows(mockRow)

	repo := NewUserRepository(mockPool)

	result, err := repo.Create(ctx, inputData)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, id, result.Id)
	assert.Equal(t, inputData.Email, result.Email)
	assert.Equal(t, inputData.PasswordHash, result.PasswordHash)
	assert.Equal(t, isVerified, result.IsVerified)
	assert.Equal(t, status, result.Status)

}

func TestCreateDuplicateUsers(t *testing.T) {
	mockpool, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mockpool.Close()

	defer func() {
		err := mockpool.ExpectationsWereMet()
		assert.NoError(t, err)
	}()

	ctx := context.Background()

	inputData := &domain.User{
		Email:        "contoh@gmail.com",
		PasswordHash: "iniadalahpasswordhash",
	}

	query := `^INSERT INTO users\(email, password_hash\)VALUES\(\$1, \$2\) RETURNING id, email, password_hash, is_verified, status, created_at$`

	mockpool.ExpectQuery(query).WithArgs(inputData.Email, inputData.PasswordHash).WillReturnError(fmt.Errorf("ERROR: Duplicate key value violates unique constraint (SQLSTATE 23505)"))

	repo := NewUserRepository(mockpool)
	result, err := repo.Create(ctx, inputData)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "something wrong when create new data")
	assert.Contains(t, err.Error(), "23505")

}
