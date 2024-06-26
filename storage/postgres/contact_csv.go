package postgres

import (
	"context"
	"fmt"
	"os"
	"task/pkg/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

type exportStorage struct {
	pool   *pgxpool.Pool
	logger logger.ILogger
}

func NewExportStorage(pool *pgxpool.Pool, log logger.ILogger) exportStorage {
	return exportStorage{
		pool:   pool,
		logger: log,
	}
}

func (s *exportStorage) ExportToCSV(ctx context.Context, tableName, outputFile string) error {
	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	exportQuery := fmt.Sprintf(`COPY (SELECT * FROM %s) TO STDOUT WITH CSV HEADER`, tableName)
	_, err = conn.Conn().PgConn().CopyTo(ctx, file, exportQuery)
	if err != nil {
		return err
	}

	return nil
}
