package repository

import (
	"context"
	"shorter-url/internal/domain"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock/v5"
	"github.com/stretchr/testify/assert"
)

func TestCreateClickEventPass(t *testing.T) {

}

// 1. SKENARIO CREATE (SUKSES)
func TestClickEventsRepository_Create_Success(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mockPool.Close()

	ctx := context.Background()

	inputData := &domain.ClickEvent{
		ShortUrlId: 1,
		IpAddress:  "192.168.1.1",
		UserAgent:  "Mozilla/5.0",
		Referer:    "https://google.com",
	}

	expectedId := int64(42)

	// 1. Kunci ekspektasi waktu menjadi STRING dengan format yang pas sesuai actual result
	expectedClickedAt := time.Now() // Format yang sama dengan yang dihasilkan oleh database

	queryPattern := `^INSERT INTO click_events\(short_url_id, ip_address, user_agent, referer\)VALUES\(\$1, \$2, \$3, \$4\) RETURNING id, ip_address, short_url_id, user_agent, referer, clicked_at$`

	// 2. Kirim expectedClickedAt sebagai string ke dalam row mock
	mockRows := pgxmock.NewRows([]string{"1", "2", "3", "4", "5", "6"}).
		AddRow(
			expectedId,
			inputData.IpAddress,
			inputData.ShortUrlId,
			inputData.UserAgent,
			inputData.Referer,
			expectedClickedAt,
		)

	mockPool.ExpectQuery(queryPattern).
		WithArgs(inputData.ShortUrlId, inputData.IpAddress, inputData.UserAgent, inputData.Referer).
		WillReturnRows(mockRows)

	repo := NewClickEventsRepository(mockPool)
	result, err := repo.Create(ctx, inputData)

	// 3. Validasi Akhir (Sekarang string vs string akan bernilai TRUE)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedId, result.Id)
	assert.Equal(t, expectedClickedAt, result.ClickedAt)

	err = mockPool.ExpectationsWereMet()
	assert.NoError(t, err)
}
