package storage

import (
	"context"
	"task/api/model"
)

type IStorage interface {
	CloseDB()
	Contacts() IContactStorage
	Categories() ICategoriestStorage
	Contactcsv() IContactcsvtStorage
}

type IContactStorage interface {
	Create(context.Context, model.Contact) (model.GetAllContact, error)
	Createcsv(context.Context, model.GetAllContact) (model.GetAllContact, error)
	GetAll(context.Context, model.GetAllContactRequest) (model.GetAllContactResponse, error)
	Delete(context.Context, string) (string, error)
	SoftDelete(context.Context, string) (string, error)
	GetByID(context.Context, string) (model.GetAllContact, error)
	Patch(context.Context, model.PatchContact) (model.GetAllContact, error)
	CheckEmail(context.Context, string) (string, error)
	History(context.Context, string) ([]model.ContactHistory, error)
}

type ICategoriestStorage interface {
	Createcat(context.Context, model.Categories) (model.Getcategoriest, error)
	GetAllcat(context.Context, model.GetAllCategoriestRequest) (model.GetAllcategoriesResponse, error)
	Deletecat(context.Context, string) (string, error)
	SoftDeletecat(context.Context, string) (string, error)
	GetByIDcat(context.Context, string) (model.Getcategoriest, error)
	Patchcat(context.Context, model.Patchcategories) (model.Getcategoriest, error)
	Checkname(context.Context, string) (string, error)
}

type IContactcsvtStorage interface {
	ExportToCSV(context.Context, string, string) error
}
