package service

import (
	"task/pkg/logger"
	"task/storage"
)

type IServiceMangaer interface {
	Categoriess() categoriesser
	Contacts() contactser
}

type Service struct {
	categoriesser categoriesser
	contactser    contactser

	logger logger.ILogger
}

func New(storage storage.IStorage, log logger.ILogger) Service {
	return Service{

		categoriesser: NewCategories(storage, log),
		contactser:    NewContact(storage, log),

		logger: log,
	}
}

func (s Service) Categoriess() categoriesser {
	return s.categoriesser
}

func (s Service) Contacts() contactser {

	return s.contactser
}
