package repository

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/stretchr/testify/require"
)

// snapshotScanner supplies valid defaults for the shared projection and overrides
// only the billing values under test. It also checks the SQL/scan column counts.
type snapshotScanner struct {
	columns string
	values  map[string]any
}

func (s snapshotScanner) Scan(dest ...any) error {
	columns := strings.Split(s.columns, ",")
	if len(columns) != len(dest) {
		return fmt.Errorf("projection/scan mismatch: %d columns, %d destinations", len(columns), len(dest))
	}
	for i, col := range columns {
		target := reflect.ValueOf(dest[i]).Elem()
		if value, ok := s.values[strings.TrimSpace(col)]; ok {
			target.Set(reflect.ValueOf(value))
		} else {
			target.SetZero()
		}
	}
	return nil
}

func TestUsageDisplayRateSnapshotPreservesActualAndHistoricalNull(t *testing.T) {
	for _, display := range []sql.NullFloat64{{}, {Float64: 1.25, Valid: true}} {
		got, err := scanUsageLog(snapshotScanner{columns: usageLogSelectColumns, values: map[string]any{
			"rate_multiplier":         3.5,
			"display_rate_multiplier": display,
		}})
		require.NoError(t, err)
		require.Equal(t, 3.5, got.RateMultiplier)
		if display.Valid {
			require.Equal(t, display.Float64, *got.DisplayRateMultiplier)
		} else {
			require.Nil(t, got.DisplayRateMultiplier, "historic rows must retain the fallback signal")
		}
	}
}

func TestUsageDisplayRateSnapshotAllInsertPaths(t *testing.T) {
	public := 1.25
	for _, display := range []*float64{nil, &public} {
		log := &service.UsageLog{UserID: 1, APIKeyID: 2, AccountID: 3, RequestID: "snapshot", Model: "gpt-5", RateMultiplier: 3.5, DisplayRateMultiplier: display, CreatedAt: time.Now().UTC()}
		prepared := prepareUsageLogInsert(log)
		require.Len(t, prepared.args, len(usageLogInsertArgTypes))
		require.Equal(t, 3.5, prepared.rateMultiplier)
		require.Equal(t, display, prepared.args[len(prepared.args)-1])

		key := usageLogBatchKey(log.RequestID, log.APIKeyID)
		query, args := buildUsageLogBatchInsertQuery([]string{key}, map[string]usageLogInsertPrepared{key: prepared})
		require.Equal(t, 3, strings.Count(query, "display_rate_multiplier"))
		require.Equal(t, display, args[len(args)-1])
		query, args = buildUsageLogBestEffortInsertQuery([]usageLogInsertPrepared{prepared})
		require.Equal(t, 3, strings.Count(query, "display_rate_multiplier"))
		require.Equal(t, display, args[len(args)-1])

		db, mock := newSQLMock(t)
		repo := &usageLogRepository{sql: db}
		mock.ExpectQuery("(?s)INSERT INTO usage_logs.*display_rate_multiplier.*\\$63").
			WithArgs(anySliceToDriverValues(prepared.args)...).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(int64(1), log.CreatedAt))
		_, err := repo.Create(context.Background(), log)
		require.NoError(t, err)
		mock.ExpectExec("(?s)INSERT INTO usage_logs.*display_rate_multiplier.*\\$63").
			WithArgs(anySliceToDriverValues(prepared.args)...).
			WillReturnResult(sqlmock.NewResult(0, 1))
		require.NoError(t, execUsageLogInsertNoResult(context.Background(), db, prepared))
		require.NoError(t, mock.ExpectationsWereMet())
	}
}

func TestGroupEntityMappingPreservesHiddenRate(t *testing.T) {
	hidden := 3.5
	for _, value := range []*float64{nil, &hidden} {
		group := groupEntityToService(&dbent.Group{RateMultiplier: 1.25, HiddenRateMultiplier: value})
		require.Equal(t, 1.25, group.RateMultiplier)
		require.Equal(t, value, group.HiddenRateMultiplier)
	}
}

func TestBatchImageJobDisplaySnapshotCreateAndScan(t *testing.T) {
	for _, display := range []sql.NullFloat64{{}, {Float64: 0.625, Valid: true}} {
		job, err := scanBatchImageJob(snapshotScanner{columns: batchImageJobColumns, values: map[string]any{
			"group_rate_multiplier":   3.5,
			"display_rate_multiplier": display,
		}})
		require.NoError(t, err)
		require.Equal(t, 3.5, job.GroupRateMultiplier)
		if display.Valid {
			require.Equal(t, display.Float64, *job.DisplayRateMultiplier)
		} else {
			require.Nil(t, job.DisplayRateMultiplier)
		}

		db, mock := newSQLMock(t)
		args := make([]driver.Value, 38)
		for i := range args {
			args[i] = sqlmock.AnyArg()
		}
		args[len(args)-1] = nil
		if display.Valid {
			args[len(args)-1] = display.Float64
		}
		mock.ExpectQuery("(?s)INSERT INTO batch_image_jobs.*display_rate_multiplier.*\\$38").
			WithArgs(args...).WillReturnError(sql.ErrNoRows)
		_, err = createBatchImageJobWithSQL(context.Background(), db, service.CreateBatchImageJobParams{DisplayRateMultiplier: job.DisplayRateMultiplier})
		require.ErrorIs(t, err, sql.ErrNoRows)
		require.NoError(t, mock.ExpectationsWereMet())
	}
}
