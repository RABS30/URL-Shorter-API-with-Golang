package repository

import (
	"context"
	"shorter-url/internal/domain"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock/v5"
	"github.com/stretchr/testify/assert"
)

func Test_Create_CreateEvent_Pass(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	assert.NoError(t, err)

	defer mockPool.Close()
	defer func() {
		err := mockPool.ExpectationsWereMet()
		assert.NoError(t, err)
	}()

	ctx := context.Background()

	inputData := &domain.ClickEvent{
		ShortUrlId: 1,
		IpAddress:  "192.168.1.1",
		UserAgent:  "Mozilla/5.0",
		Referer:    "https://google.com",
	}

	dateExample := time.Now()
	expectedData := &domain.ClickEvent{
		Id:         42,
		ShortUrlId: inputData.ShortUrlId,
		IpAddress:  inputData.IpAddress,
		UserAgent:  inputData.UserAgent,
		Referer:    inputData.Referer,
		ClickedAt:  dateExample,
	}

	query := `^INSERT INTO click_events\(short_url_id, ip_address, user_agent, referer\)VALUES\(\$1, \$2, \$3, \$4\) RETURNING id, ip_address, short_url_id, user_agent, referer, clicked_at$`
	mockRows := pgxmock.NewRows([]string{"id", "ip_address", "short_url_id", "user_agent", "referer", "clicked_at"}).
		AddRow(
			expectedData.Id,
			expectedData.IpAddress,
			expectedData.ShortUrlId,
			expectedData.UserAgent,
			expectedData.Referer,
			expectedData.ClickedAt,
		)

	mockPool.ExpectQuery(query).WithArgs(inputData.ShortUrlId, inputData.IpAddress, inputData.UserAgent, inputData.Referer).WillReturnRows(mockRows)

	repo := NewClickEventsRepository(mockPool)
	result, err := repo.Create(ctx, inputData)

	assert.NoError(t, err)
	assert.Equal(t, expectedData, result)
}

func Test_Delete_CreateEvent_Pass(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	assert.NoError(t, err)

	defer mockPool.Close()
	defer func() {
		err := mockPool.ExpectationsWereMet()
		assert.NoError(t, err)
	}()

	ctx := context.Background()
	id := int64(1)

	query := `^DELETE FROM click_events WHERE id = \$1$`

	mockPool.ExpectExec(query).WithArgs(id).WillReturnResult(pgxmock.NewResult("DELETE", 1))

	repo := NewClickEventsRepository(mockPool)
	err = repo.Delete(ctx, id)

	assert.NoError(t, err)

}

func Test_Delete_CreateEvent_Fail(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	assert.NoError(t, err)

	defer mockPool.Close()
	defer func() {
		err := mockPool.ExpectationsWereMet()
		assert.NoError(t, err)
	}()

	ctx := context.Background()
	id := int64(99)

	query := `^DELETE FROM click_events WHERE id = \$1$`

	mockPool.ExpectExec(query).WithArgs(id).WillReturnResult(pgxmock.NewResult("DELETE", 0))

	repo := NewClickEventsRepository(mockPool)
	err = repo.Delete(ctx, id)

	assert.Error(t, err)
	assert.ErrorContains(t, err, "there is no data deleted")
}

func Test_FindById_CreateEvent_Pass(t *testing.T) {
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

	expectedData := &domain.ClickEvent{
		Id:         id,
		IpAddress:  "127.0.0.1",
		ShortUrlId: 2,
		UserAgent:  "Mozilla",
		Referer:    "Google",
		ClickedAt:  now,
	}

	query := `^SELECT id, short_url_id, ip_address, user_agent, referer, clicked_at FROM click_events WHERE id = \$1$`

	mockRow := pgxmock.NewRows([]string{"id", "short_url_id", "ip_address", "user_agent", "referer", "clicked_at"}).AddRow(
		expectedData.Id,
		expectedData.ShortUrlId,
		expectedData.IpAddress,
		expectedData.UserAgent,
		expectedData.Referer,
		expectedData.ClickedAt,
	)

	mockPool.ExpectQuery(query).WithArgs(id).WillReturnRows(mockRow)

	repo := NewClickEventsRepository(mockPool)
	result, err := repo.FindById(ctx, id)

	assert.NoError(t, err)
	assert.Equal(t, expectedData, result)

}

func Test_FindByShortCode_CreateEvent_Pass(t *testing.T) {
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

	expectedData := domain.ClickEvent{
		Id:         1,
		ShortUrlId: 2,
		IpAddress:  "127.0.0.1",
		UserAgent:  "Mozilla",
		Referer:    "Google",
		ClickedAt:  now,
	}

	query := `^SELECT ce.id, ce.short_url_id, ce.ip_address, ce.user_agent, ce.referer, ce.clicked_at FROM click_events ce JOIN short_urls su ON ce.short_url_id = su.id WHERE su.short_code = \$1$`

	mockRow := pgxmock.NewRows([]string{"id", "short_url_id", "ip_address", "user_agent", "referer", "clicked_at"}).AddRow(
		expectedData.Id,
		expectedData.ShortUrlId,
		expectedData.IpAddress,
		expectedData.UserAgent,
		expectedData.Referer,
		expectedData.ClickedAt,
	)

	mockPool.ExpectQuery(query).WithArgs(shortCode).WillReturnRows(mockRow)

	repo := NewClickEventsRepository(mockPool)
	result, err := repo.FindByShortCode(ctx, shortCode)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, expectedData, result[0])
}

func Test_FilterByDat_CreateEvent_Pass(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	assert.NoError(t, err)

	defer mockPool.Close()
	defer func() {
		err := mockPool.ExpectationsWereMet()
		assert.NoError(t, err)
	}()

	ctx := context.Background()
	date := time.Now()

	expectedData := domain.ClickEvent{
		Id:         1,
		IpAddress:  "127.0.0.1",
		ShortUrlId: 2,
		UserAgent:  "Mozilla",
		Referer:    "Google",
		ClickedAt:  date,
	}

	query := `^SELECT id, ip_address, short_url_id, user_agent, referer, clicked_at FROM click_events WHERE DATE\(clicked_at\) = \$1$`

	mockRow := pgxmock.NewRows([]string{"id", "ip_address", "short_url_id", "user_agent", "referer", "clicked_at"}).AddRow(
		expectedData.Id,
		expectedData.IpAddress,
		expectedData.ShortUrlId,
		expectedData.UserAgent,
		expectedData.Referer,
		expectedData.ClickedAt,
	)

	mockPool.ExpectQuery(query).WithArgs(date.Format("2006-01-02")).WillReturnRows(mockRow)

	repo := NewClickEventsRepository(mockPool)
	result, err := repo.FilterByDate(ctx, date)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, expectedData, result[0])
}

func Test_FindAll_CreateEvent_Pass(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	assert.NoError(t, err)

	defer mockPool.Close()
	defer func() {
		err := mockPool.ExpectationsWereMet()
		assert.NoError(t, err)
	}()

	ctx := context.Background()
	now := time.Now()

	expectedData := domain.ClickEvent{
		Id:         1,
		IpAddress:  "127.0.0.1",
		ShortUrlId: 2,
		UserAgent:  "Mozilla",
		Referer:    "Google",
		ClickedAt:  now,
	}

	query := `^SELECT id, short_url_id, ip_address, user_agent, referer, clicked_at FROM click_events$`

	mockRow := pgxmock.NewRows([]string{"id", "short_url_id", "ip_address", "user_agent", "referer", "clicked_at"}).AddRow(
		expectedData.Id,
		expectedData.ShortUrlId,
		expectedData.IpAddress,
		expectedData.UserAgent,
		expectedData.Referer,
		expectedData.ClickedAt,
	)

	mockPool.ExpectQuery(query).WillReturnRows(mockRow)

	repo := NewClickEventsRepository(mockPool)
	result, err := repo.FindAll(ctx)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, expectedData, result[0])
}
