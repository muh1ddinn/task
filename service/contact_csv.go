package service

import (
	"context"
	"task/pkg/logger"
	"task/storage"
)

type exportService struct {
	storage storage.IStorage
	logger  logger.ILogger
}

func NewExportService(storage storage.IStorage, logger logger.ILogger) exportService {
	return exportService{
		storage: storage,
		logger:  logger,
	}
}

func (s exportService) ExportToCSV(ctx context.Context, tableName, outputFile string) error {
	err := s.storage.Contactcsv().ExportToCSV(ctx, tableName, outputFile)
	if err != nil {
		s.logger.Error("ERROR in service layer while ExportToCSV", logger.Error(err))
		return err
	}

	return nil
}

