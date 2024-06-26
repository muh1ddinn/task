package service

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"task/api/model"
	"task/pkg/logger"
	"task/pkg/validate"
	"task/storage"
)

type contactser struct {
	storage storage.IStorage
	logger  logger.ILogger
}

func NewContact(storage storage.IStorage, logger logger.ILogger) contactser {

	return contactser{
		storage: storage,
		logger:  logger,
	}
}

func (u contactser) Create(ctx context.Context, contact model.Contact) (model.GetAllContact, error) {

	if err := validate.ValidatePhone(contact.Phone); err != nil {
		fmt.Println(contact.Phone)
		return model.GetAllContact{}, err

	}

	msg, err := u.storage.Contacts().Create(ctx, contact)
	if err != nil {
		u.logger.Error("ERROR in service layer while create :Create", logger.Error(err))
		return model.GetAllContact{}, err
	}

	return msg, nil
}

func (s contactser) Patch(ctx context.Context, contact model.PatchContact) (model.GetAllContact, error) {

	if err := validate.ValidatePhone(contact.Phone); err != nil {
		fmt.Println(contact.Phone)
		return model.GetAllContact{}, err

	}

	msg, err := s.storage.Contacts().Patch(ctx, contact)
	if err != nil {
		fmt.Println(err)
		s.logger.Error("error in service layer while getting contact:Patch", logger.Error(err))
		return model.GetAllContact{}, err
	}
	return msg, nil
}

func (s contactser) GetAll(ctx context.Context, contact model.GetAllContactRequest) ([]model.GetAllContact, error) {
	users, err := s.storage.Contacts().GetAll(ctx, contact)
	if err != nil {
		s.logger.Error("error in service layer while getting contact:GetAll ", logger.Error(err))
		return nil, err
	}
	return users.Contact, nil
}

func (s contactser) GetByID(ctx context.Context, id string) (model.GetAllContact, error) {
	users, err := s.storage.Contacts().GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get contact by ID:GetByID", logger.Error(err))
		return users, err
	}

	return users, nil
}

func (s contactser) Deletesoft(ctx context.Context, id string) (string, error) {
	usersid, err := s.storage.Contacts().SoftDelete(ctx, id)
	if err != nil {
		s.logger.Error("error in service layer while getting contact:DELETEsoft", logger.Error(err))
		return usersid, err
	}
	return usersid, nil
}

func (s contactser) Delet(ctx context.Context, id string) (string, error) {
	usersid, err := s.storage.Contacts().Delete(ctx, id)
	if err != nil {
		s.logger.Error("error in service layer while getting contact:DELETE ", logger.Error(err))
		return usersid, err
	}
	return usersid, nil
}

func (s contactser) Checkemail(ctx context.Context, req string) (string, error) {
	msg, err := s.storage.Contacts().CheckEmail(ctx, req)
	if err != nil {
		s.logger.Error("error in service layer while getting contact:DELETE ", logger.Error(err))
		return msg, err
	}
	return msg, nil
}

func (s contactser) History(ctx context.Context, id string) ([]model.ContactHistory, error) {
	users, err := s.storage.Contacts().History(ctx, id)
	fmt.Println(users)
	if err != nil {
		s.logger.Error("failed to get contact by history", logger.Error(err))
		return users, err
	}

	return users, nil
}

func (s contactser) ImportContacts(ctx context.Context, file io.Reader) error {
	contacts, err := s.parseCSV(file)
	if err != nil {
		return err
	}

	for _, contact := range contacts {
		if _, err := s.storage.Contacts().Createcsv(ctx, contact); err != nil {
			s.logger.Error("Error importing contact", logger.Error(err))
			return err
		}
	}

	return nil
}

func (s contactser) parseCSV(file io.Reader) ([]model.GetAllContact, error) {
	reader := csv.NewReader(file)
	var contacts []model.GetAllContact

	// Read header
	_, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV header: %v", err)
	}

	// Read rows
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading CSV row: %v", err)
		}

		contact := model.GetAllContact{
			Id:       record[0],
			Phone:    record[1],
			Name:     record[2],
			Email:    record[3],
			Address:  record[4],
			Category: record[5],
		}
		contacts = append(contacts, contact)
	}
	return contacts, nil
}
