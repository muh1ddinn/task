package service

import (
	"context"
	"fmt"
	"task/api/model"
	"task/pkg/logger"
	"task/storage"
)

type categoriesser struct {
	storage storage.IStorage
	logger  logger.ILogger
}

func NewCategories(storage storage.IStorage, logger logger.ILogger) categoriesser {

	return categoriesser{
		storage: storage,
		logger:  logger,
	}
}

func (u categoriesser) Create(ctx context.Context, categories model.Categories) (model.Getcategoriest, error) {

	msg, err := u.storage.Categories().Createcat(ctx, categories)
	if err != nil {
		u.logger.Error("ERROR in service layer while create :Create", logger.Error(err))
		return model.Getcategoriest{}, err
	}

	return msg, nil
}

func (s categoriesser) Patchcat(ctx context.Context, categories model.Patchcategories) (model.Getcategoriest, error) {

	msg, err := s.storage.Categories().Patchcat(ctx, categories)
	if err != nil {
		fmt.Println(err)
		s.logger.Error("error in service layer while getting categories:Patch", logger.Error(err))
		return model.Getcategoriest{}, err
	}
	return msg, nil
}

func (s categoriesser) GetAll(ctx context.Context, categories model.GetAllCategoriestRequest) ([]model.GetAllcategoriesResponse, error) {
	us, err := s.storage.Categories().GetAllcat(ctx, categories)
	if err != nil {
		s.logger.Error("error in service layer while getting categories:GetAll ", logger.Error(err))
		return nil, err
	}
	return []model.GetAllcategoriesResponse{us}, nil
}

func (s categoriesser) GetByID(ctx context.Context, id string) (model.Getcategoriest, error) {
	users, err := s.storage.Categories().GetByIDcat(ctx, id)
	if err != nil {
		s.logger.Error("failed to get categories by ID:GetByID", logger.Error(err))
		return users, err
	}

	return users, nil
}

func (s categoriesser) Deletesoft(ctx context.Context, id string) (string, error) {
	usersid, err := s.storage.Categories().SoftDeletecat(ctx, id)
	if err != nil {
		s.logger.Error("error in service layer while getting categories:DELETEsoft", logger.Error(err))
		return usersid, err
	}
	return usersid, nil
}

func (s categoriesser) Delet(ctx context.Context, id string) (string, error) {
	usersid, err := s.storage.Categories().Deletecat(ctx, id)
	if err != nil {
		s.logger.Error("error in service layer while getting categories:DELETE ", logger.Error(err))
		return usersid, err
	}
	return usersid, nil
}

func (s categoriesser) Checkename(ctx context.Context, req string) (string, error) {
	msg, err := s.storage.Categories().Checkname(ctx, req)
	if err != nil {
		s.logger.Error("error in service layer while getting contact:DELETE ", logger.Error(err))
		return msg, err
	}
	return msg, nil
}
