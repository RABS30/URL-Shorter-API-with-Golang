package repository

import (
	"context"
	"fmt"
	"shorter-url/internal/domain"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v5"
	"github.com/stretchr/testify/assert"
)

func Test_Create_ShortUrls_Pass(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	assert.NoError(t, err)

	defer mockPool.Close()
	defer func() {
		err := mockPool.ExpectationsWereMet()
		assert.NoError(t, err)
	}()

	ctx := context.Background()
	expiredAt := time.Now()

	inputData := &domain.ShortUrl{
		UserId:      1,
		OriginalUrl: "https://example.com/very/long/url",
		ShortCode:   "abcd",
		ExpiredAt:   expiredAt,
	}

	expectedData := &domain.ShortUrl{
		Id:          10,
		UserId:      inputData.UserId,
		OriginalUrl: inputData.OriginalUrl,
		ShortCode:   inputData.ShortCode,
		ExpiredAt:   inputData.ExpiredAt,
	}

	queryPattern := `^INSERT INTO short_urls \(user_id, original_url, short_code, expired_at\) VALUES \(\$1, \$2, \$3, \$4\) RETURNING id$`

	mockRow := pgxmock.NewRows([]string{"id"}).AddRow(expectedData.Id)

	mockPool.ExpectQuery(queryPattern).WithArgs(inputData.UserId, inputData.OriginalUrl, inputData.ShortCode, inputData.ExpiredAt).WillReturnRows(mockRow)

	repo := NewShortUrlRepository(mockPool)
	result, err := repo.Create(ctx, inputData)

	assert.NoError(t, err)
	assert.Equal(t, expectedData, result)
}

func Test_Create_ShortUrls_Fail(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	assert.NoError(t, err)

	defer mockPool.Close()
	defer func() {
		err := mockPool.ExpectationsWereMet()
		assert.NoError(t, err)
	}()

	ctx := context.Background()
	expiredAt := time.Now().Add(24 * time.Hour)

	inputData := &domain.ShortUrl{
		UserId:      1,
		OriginalUrl: "https://example.com",
		ShortCode:   "exmpl",
		ExpiredAt:   expiredAt,
	}

	queryPattern := `^INSERT INTO short_urls \(user_id, original_url, short_code, expired_at\) VALUES \(\$1, \$2, \$3, \$4\) RETURNING id$`

	mockPool.ExpectQuery(queryPattern).WithArgs(inputData.UserId, inputData.OriginalUrl, inputData.ShortCode, inputData.ExpiredAt).WillReturnError(fmt.Errorf("db error"))

	repo := NewShortUrlRepository(mockPool)
	result, err := repo.Create(ctx, inputData)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.ErrorContains(t, err, "db error")
}

func Test_Delete_ShortUrls_Pass(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	assert.NoError(t, err)

	defer mockPool.Close()
	defer func() {
		err := mockPool.ExpectationsWereMet()
		assert.NoError(t, err)
	}()

	ctx := context.Background()
	id := int64(1)

	query := `^DELETE FROM short_urls WHERE id = \$1$`

	mockPool.ExpectExec(query).WithArgs(id).WillReturnResult(pgxmock.NewResult("DELETE", 1))

	repo := NewShortUrlRepository(mockPool)
	err = repo.Delete(ctx, id)

	assert.NoError(t, err)
}

func Test_Delete_ShortUrls_Fail(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	assert.NoError(t, err)

	defer mockPool.Close()
	defer func() {
		err := mockPool.ExpectationsWereMet()
		assert.NoError(t, err)
	}()

	ctx := context.Background()
	id := int64(99)

	query := `^DELETE FROM short_urls WHERE id = \$1$`

	mockPool.ExpectExec(query).WithArgs(id).WillReturnResult(pgxmock.NewResult("DELETE", 0))

	repo := NewShortUrlRepository(mockPool)
	err = repo.Delete(ctx, id)

	assert.Error(t, err)
	assert.ErrorContains(t, err, "resource not found")
}

func Test_FindById_ShortUrls_Pass(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	assert.NoError(t, err)

	defer mockPool.Close()
	defer func() {
		err := mockPool.ExpectationsWereMet()
		assert.NoError(t, err)
	}()

	ctx := context.Background()

	id := int64(1)
	now := time.Now()

	expectedData := &domain.ShortUrl{
		Id:          id,
		UserId:      2,
		ShortCode:   "abcd",
		OriginalUrl: "https://google.com",
		ExpiredAt:   now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	query := `^SELECT id, user_id, short_code, original_url, expired_at, created_at, updated_at FROM short_urls WHERE id = \$1$`

	mockRow := pgxmock.NewRows([]string{"id", "user_id", "short_code", "original_url", "expired_at", "created_at", "updated_at"}).AddRow(
		expectedData.Id,
		expectedData.UserId,
		expectedData.ShortCode,
		expectedData.OriginalUrl,
		expectedData.ExpiredAt,
		expectedData.CreatedAt,
		expectedData.UpdatedAt,
	)

	mockPool.ExpectQuery(query).WithArgs(id).WillReturnRows(mockRow)

	repo := NewShortUrlRepository(mockPool)
	result, err := repo.FindById(ctx, id)

	assert.NoError(t, err)
	assert.Equal(t, expectedData, result)
}

func Test_FindById_ShortUrls_Fail(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	assert.NoError(t, err)

	defer mockPool.Close()
	defer func() {
		err := mockPool.ExpectationsWereMet()
		assert.NoError(t, err)
	}()

	ctx := context.Background()

	id := int64(5)

	query := `^SELECT id, user_id, short_code, original_url, expired_at, created_at, updated_at FROM short_urls WHERE id = \$1$`

	mockPool.ExpectQuery(query).WithArgs(id).WillReturnError(pgx.ErrNoRows)

	repo := NewShortUrlRepository(mockPool)
	result, err := repo.FindById(ctx, id)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "resource not found")
}

func Test_FindByUserId_ShortUrls_Pass(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	assert.NoError(t, err)

	defer mockPool.Close()
	defer func() {
		err := mockPool.ExpectationsWereMet()
		assert.NoError(t, err)
	}()

	ctx := context.Background()

	userId := int64(2)
	now := time.Now()

	expectedData := domain.ShortUrl{
		Id:          1,
		UserId:      userId,
		ShortCode:   "abcd",
		OriginalUrl: "https://google.com",
		ExpiredAt:   now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	query := `^SELECT id, user_id, short_code, original_url, expired_at, created_at, updated_at FROM short_urls WHERE user_id = \$1$`

	mockRow := pgxmock.NewRows([]string{"id", "user_id", "short_code", "original_url", "expired_at", "created_at", "updated_at"}).AddRow(
		expectedData.Id,
		expectedData.UserId,
		expectedData.ShortCode,
		expectedData.OriginalUrl,
		expectedData.ExpiredAt,
		expectedData.CreatedAt,
		expectedData.UpdatedAt,
	)

	mockPool.ExpectQuery(query).WithArgs(userId).WillReturnRows(mockRow)

	repo := NewShortUrlRepository(mockPool)
	result, err := repo.FindByUserId(ctx, userId)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, expectedData, result[0])
}

func Test_FindByShortCode_ShortUrls_Pass(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	assert.NoError(t, err)

	defer mockPool.Close()
	defer func() {
		err := mockPool.ExpectationsWereMet()
		assert.NoError(t, err)
	}()

	ctx := context.Background()

	shortCode := "abcd"
	now := time.Now()

	expectedData := &domain.ShortUrl{
		Id:          1,
		UserId:      2,
		ShortCode:   shortCode,
		OriginalUrl: "https://google.com",
		ExpiredAt:   now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	query := `^SELECT id, user_id, short_code, original_url, expired_at, created_at, updated_at FROM short_urls WHERE short_code = \$1$`

	mockRow := pgxmock.NewRows([]string{"id", "user_id", "short_code", "original_url", "expired_at", "created_at", "updated_at"}).AddRow(
		expectedData.Id,
		expectedData.UserId,
		expectedData.ShortCode,
		expectedData.OriginalUrl,
		expectedData.ExpiredAt,
		expectedData.CreatedAt,
		expectedData.UpdatedAt,
	)

	mockPool.ExpectQuery(query).WithArgs(shortCode).WillReturnRows(mockRow)

	repo := NewShortUrlRepository(mockPool)
	result, err := repo.FindByShortCode(ctx, shortCode)

	assert.NoError(t, err)
	assert.Equal(t, expectedData, result)

}
