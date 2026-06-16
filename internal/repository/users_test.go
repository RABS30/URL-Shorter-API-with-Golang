package repository

import (
	"context"
	"fmt"
	"log"
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

func TestUpdateUsersPass(t *testing.T) {
	mockpool, err := pgxmock.NewPool()
	assert.NoError(t, err)

	defer mockpool.Close()
	defer func() {
		err := mockpool.ExpectationsWereMet()
		assert.NoError(t, err)
	}()

	ctx := context.Background()

	inputData := &domain.User{
		Id:           9,
		Email:        "",
		PasswordHash: "iniadalahpasswordhash",
		IsVerified:   true,
		Status:       "active",
	}
	createdAt := time.Now()
	query := `^UPDATE users SET password_hash = \$1, is_verified = \$2, status = \$3 WHERE id = \$4 RETURNING id, email, password_hash, is_verified, status, created_at$`

	mockRow := pgxmock.NewRows([]string{"1", "2", "3", "4", "5", "6"}).AddRow(
		inputData.Id,
		inputData.Email,
		inputData.PasswordHash,
		inputData.IsVerified,
		inputData.Status,
		createdAt,
	)
	mockpool.ExpectQuery(query).WithArgs(inputData.PasswordHash, inputData.IsVerified, inputData.Status, inputData.Id).WillReturnRows(mockRow)

	repo := NewUserRepository(mockpool)
	result, err := repo.Update(ctx, inputData)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	log.Printf("result: %+v", result)
}

func TestUpdateUsersFail(t *testing.T) {
	mockpool, err := pgxmock.NewPool()
	assert.NoError(t, err)

	defer mockpool.Close()
	defer func() {
		err := mockpool.ExpectationsWereMet()
		assert.NoError(t, err)
	}()

	ctx := context.Background()

	inputData := &domain.User{
		Id:           1,
		Email:        "",
		PasswordHash: "passwordhashhash",
		IsVerified:   false,
		Status:       "inactive",
		CreatedAt:    time.Now(),
	}

	query := `^UPDATE users SET password_hash = \$1, is_verified = \$2, status = \$3 WHERE id = \$4 RETURNING id, email, password_hash, is_verified, status, created_at$`

	mockpool.ExpectQuery(query).WithArgs(inputData.PasswordHash, inputData.IsVerified, inputData.Status, inputData.Id).WillReturnError(fmt.Errorf("User ID Not Found"))

	repo := NewUserRepository(mockpool)
	result, err := repo.Update(ctx, inputData)

	assert.Nil(t, result)
	assert.ErrorContains(t, err, "User ID Not Found")
}

func TestDeteleUsersPass(t *testing.T) {
	mockpool, err := pgxmock.NewPool()
	assert.NoError(t, err)

	defer mockpool.Close()
	defer func() {
		err := mockpool.ExpectationsWereMet()
		assert.NoError(t, err)
	}()

	ctx := context.Background()

	inputData := &domain.User{
		Id: 1,
	}

	query := `^DELETE FROM users WHERE id = \$1$`

	mockpool.ExpectExec(query).WithArgs(inputData.Id).WillReturnResult(pgxmock.NewResult("DELETE", 1))

	repo := NewUserRepository(mockpool)
	err = repo.Delete(ctx, inputData.Id)

	assert.NoError(t, err)
}

func TestDeleteUsersFail(t *testing.T) {
	mockpool, err := pgxmock.NewPool()
	assert.NoError(t, err)

	defer mockpool.Close()
	defer func() {
		mockpool.ExpectationsWereMet()
	}()

	ctx := context.Background()

	inputData := &domain.User{
		Id: 99,
	}

	query := `^DELETE FROM users WHERE id = \$1$`

	mockpool.ExpectExec(query).WithArgs(inputData.Id).WillReturnResult(pgxmock.NewResult("DELETE", 0))
	repo := NewUserRepository(mockpool)

	err = repo.Delete(ctx, inputData.Id)

	assert.ErrorContains(t, err, "there is no data deleted")
}

func TestFindByIdPass(t *testing.T) {
	mockpool, err := pgxmock.NewPool()
	assert.NoError(t, err)

	defer mockpool.Close()
	defer func() {
		mockpool.ExpectationsWereMet()
	}()

	ctx := context.Background()

	inputData := &domain.User{
		Id:           1,
		Email:        "rokubi27@gmail.com",
		PasswordHash: "this_password_is_hashing",
		IsVerified:   true,
		Status:       "active",
		CreatedAt:    time.Now(),
	}

	query := `^SELECT id, email, password_hash, is_verified, status, created_at FROM users WHERE id = \$1$`

	mockRow := mockpool.NewRows([]string{"1", "2", "3", "4", "5", "6"}).AddRow(
		inputData.Id,
		inputData.Email,
		inputData.PasswordHash,
		inputData.IsVerified,
		inputData.Status,
		inputData.CreatedAt,
	)

	mockpool.ExpectQuery(query).WithArgs(inputData.Id).WillReturnRows(mockRow)

	repo := NewUserRepository(mockpool)
	result, err := repo.FindById(ctx, inputData.Id)

	assert.Equal(t, inputData, result)

	assert.Nil(t, err)

}

func TestFindByIdFail(t *testing.T) {
	mockpool, err := pgxmock.NewPool()
	assert.NoError(t, err)

	defer mockpool.Close()
	defer func() {
		err := mockpool.ExpectationsWereMet()
		assert.NoError(t, err)
	}()

	ctx := context.Background()

	inputData := &domain.User{
		Id: 99,
	}

	query := `^SELECT id, email, password_hash, is_verified, status, created_at FROM users WHERE id = \$1$`
	mockpool.ExpectQuery(query).WithArgs(inputData.Id).WillReturnError(fmt.Errorf("user not found"))

	repo := NewUserRepository(mockpool)
	result, err := repo.FindById(ctx, inputData.Id)

	assert.Nil(t, result)
	assert.ErrorContains(t, err, "user not found")
}
func TestFindByEmailPass(t *testing.T) {
	mockpool, err := pgxmock.NewPool()
	assert.NoError(t, err)

	defer mockpool.Close()
	defer func() {
		mockpool.ExpectationsWereMet()
	}()

	ctx := context.Background()

	inputData := &domain.User{
		Id:           1,
		Email:        "rokubi27@gmail.com",
		PasswordHash: "this_password_is_hashing",
		IsVerified:   true,
		Status:       "active",
		CreatedAt:    time.Now(),
	}

	query := `^SELECT id, email, password_hash, is_verified, status, created_at FROM users WHERE email = \$1$`

	mockRow := mockpool.NewRows([]string{"1", "2", "3", "4", "5", "6"}).AddRow(
		inputData.Id,
		inputData.Email,
		inputData.PasswordHash,
		inputData.IsVerified,
		inputData.Status,
		inputData.CreatedAt,
	)

	mockpool.ExpectQuery(query).WithArgs(inputData.Email).WillReturnRows(mockRow)

	repo := NewUserRepository(mockpool)
	result, err := repo.FindByEmail(ctx, inputData.Email)

	assert.Nil(t, err)
	assert.Equal(t, inputData, result)

}

func TestFindByEmailFail(t *testing.T) {
	mockpool, err := pgxmock.NewPool()
	assert.NoError(t, err)

	defer mockpool.Close()
	defer func() {
		err := mockpool.ExpectationsWereMet()
		assert.NoError(t, err)
	}()

	ctx := context.Background()

	inputData := &domain.User{
		Email: "rokubi27@gmail.com",
	}

	query := `^SELECT id, email, password_hash, is_verified, status, created_at FROM users WHERE email = \$1$`
	mockpool.ExpectQuery(query).WithArgs(inputData.Email).WillReturnError(fmt.Errorf("email not found"))

	repo := NewUserRepository(mockpool)
	result, err := repo.FindByEmail(ctx, inputData.Email)

	assert.Nil(t, result)
	assert.ErrorContains(t, err, "email not found")
}
