package postgres

import (
	"context"
	"fmt"
	"log"
	"task/api/model"
	"task/pkg/logger"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type contactRepro struct {
	db     *pgxpool.Pool
	logger logger.ILogger
}

func Newcontact(db *pgxpool.Pool, log logger.ILogger) contactRepro {
	return contactRepro{
		db:     db,
		logger: log,
	}
}

func (c *contactRepro) Create(ctx context.Context, contact model.Contact) (model.GetAllContact, error) {

	id := uuid.New()

	query := `INSERT INTO contact(
		id,
		phone,
		name,
		email,
		address,
		Category,
		created_at
	) VALUES($1, $2, $3, $4, $5, $6,$7)`

	_, err := c.db.Exec(ctx, query,
		id.String(),
		contact.Phone,
		contact.Name,
		contact.Email,
		contact.Address,
		contact.Category,
		time.Now(),
	)
	fmt.Println(err)
	if err != nil {
		c.logger.Error("error while creating user:", logger.Error(err))
		return model.GetAllContact{}, err
	}

	a, err := c.GetByID(ctx, id.String())
	if err != nil {
		log.Println("error while getting category by id")
		return a, err
	}

	return a, nil
}

func (c *contactRepro) GetByID(ctx context.Context, id string) (model.GetAllContact, error) {

	contact := model.GetAllContact{}

	query := `SELECT 
	id,
	phone,
	name,
	email,
	address,
	Category
	--created_at
	FROM contact WHERE id=$1 AND deleted_at =0`

	row := c.db.QueryRow(ctx, query, id)

	err := row.Scan(
		&contact.Id,
		&contact.Phone,
		&contact.Name,
		&contact.Email,
		&contact.Address,
		&contact.Category,
		//&contact.Created_at,
	)
	if err != nil {
		fmt.Println(err)
		c.logger.Error("error while getting user by ID: ", logger.Error(err))
		return contact, err
	}

	return contact, nil
}

func (c *contactRepro) SoftDelete(ctx context.Context, id string) (string, error) {

	currentContact, err := c.GetByID(ctx, id)
	if err != nil {
		return "", err
	}

	err = c.logHistory(ctx, currentContact, "DELETE")
	if err != nil {
		return "", err
	}

	query := `UPDATE contact SET deleted_at = 1 WHERE id = $1 AND deleted_at =0`

	_, err = c.db.Exec(ctx, query, id)
	if err != nil {
		c.logger.Error("failed to soft delete user from database", logger.Error(err))
		return id, err
	}

	return id, nil
}

func (c *contactRepro) Delete(ctx context.Context, id string) (string, error) {

	deleteHistoryQuery := `DELETE FROM contact_history WHERE contact_id = $1`
	_, err := c.db.Exec(ctx, deleteHistoryQuery, id)
	if err != nil {
		c.logger.Error("failed to delete user history from database", logger.Error(err))
		return id, err
	}

	deleteContactQuery := `DELETE FROM contact WHERE id = $1`
	_, err = c.db.Exec(ctx, deleteContactQuery, id)

	fmt.Println(err)
	if err != nil {
		c.logger.Error("failed to delete user from database", logger.Error(err))
		return id, err
	}

	return id, nil
}

func (c *contactRepro) Patch(ctx context.Context, contact model.PatchContact) (model.GetAllContact, error) {
	fmt.Println(contact.Email, "email")
	fmt.Println(contact.Id, "uuid")

	currentContact, err := c.GetByID(ctx, contact.Id)
	if err != nil {
		return model.GetAllContact{}, err
	}

	err = c.logHistory(ctx, currentContact, "UPDATE")
	fmt.Println(err)
	if err != nil {
		return model.GetAllContact{}, err
	}

	query := `UPDATE contact SET 
		phone = $2,
		email = $3,
		updated_at = CURRENT_TIMESTAMP
	WHERE id = $1 AND deleted_at =0`
	_, err = c.db.Exec(ctx, query, contact.Id, contact.Phone, contact.Email)
	if err != nil {
		c.logger.Error("failed to patch user from database", logger.Error(err))
		return model.GetAllContact{}, err
	}

	updatedContact, err := c.GetByID(ctx, contact.Id)
	if err != nil {
		log.Println("error while getting contact by id")
		return updatedContact, err
	}

	return updatedContact, nil
}

func (c *contactRepro) GetAll(ctx context.Context, req model.GetAllContactRequest) (model.GetAllContactResponse, error) {
	resp := model.GetAllContactResponse{}

	offset := (req.Page - 1) * req.Limit
	baseFilter := " WHERE deleted_at = 0"
	searchFilter := ""
	args := []interface{}{}

	// If a search term is provided, add it to the filter
	if req.Search != "" {
		searchTerm := fmt.Sprintf("%%%s%%", req.Search)
		searchFilter = ` AND (
            email ILIKE $1 OR 
            name ILIKE $1 OR 
            category ILIKE $1
        )`
		args = append(args, searchTerm)
	}

	// Append LIMIT and OFFSET parameters to args slice
	args = append(args, offset, req.Limit)

	// Construct the SQL query
	query := `
        SELECT 
            count(id) OVER() AS total_count,
            id,
            phone,
            name,
            email,
            address,
            category,
            created_at
        FROM 
            contact` + baseFilter + searchFilter + `
        ORDER BY 
            created_at DESC
        OFFSET $` + fmt.Sprint(len(args)-1) + ` LIMIT $` + fmt.Sprint(len(args))

	fmt.Println("query: ", query)
	fmt.Println("args: ", args)

	// Execute the query
	rows, err := c.db.Query(ctx, query, args...)
	if err != nil {
		return resp, err
	}
	defer rows.Close()

	// Iterate over the rows
	for rows.Next() {
		var contact model.GetAllContact
		var createdAt time.Time

		errScan := rows.Scan(
			&resp.Count,
			&contact.Id,
			&contact.Phone,
			&contact.Name,
			&contact.Email,
			&contact.Address,
			&contact.Category,
			&createdAt,
		)
		if errScan != nil {
			c.logger.Error("error while scanning contact info: ", logger.Error(errScan))
			return resp, errScan
		}

		// Format created_at field
		contact.Created_at = createdAt.Format("2006-01-02 15:04:05")

		// Append contact to response slice
		resp.Contact = append(resp.Contact, contact)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		c.logger.Error("error while iterating over rows: ", logger.Error(err))
		return model.GetAllContactResponse{}, err
	}

	fmt.Printf("Final response: %+v\n", resp)
	return resp, nil
}

func (c *contactRepro) CheckEmail(ctx context.Context, email string) (string, error) {
	contact := model.GetAllContact{}
	fmt.Println(email)
	query := `SELECT 
		phone,
		email
	FROM contact WHERE email=$1 AND deleted_at = 0`

	row := c.db.QueryRow(ctx, query, email)

	err := row.Scan(
		&contact.Phone,
		&contact.Email,
	)
	if err != nil {
		if err.Error() == "no rows in result set" {
			fmt.Println("No contact found with email:", email)
			return "", nil
		}
		fmt.Println(err)
		c.logger.Error("error while getting contact by email: ", logger.Error(err))
		return "", err
	}

	return contact.Email, nil
}

func (c *contactRepro) logHistory(ctx context.Context, contact model.GetAllContact, changeType string) error {
	id := uuid.New()
	fmt.Println(contact)
	query := `INSERT INTO contact_history (
		id, contact_id, phone, name, email, address, category, changed_at, change_type
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := c.db.Exec(ctx, query, id, contact.Id, contact.Phone, contact.Name, contact.Email, contact.Address, contact.Category, time.Now(), changeType)
	if err != nil {
		fmt.Println(err)
		c.logger.Error("failed to log history: ", logger.Error(err))
	}
	return err
}

func (c *contactRepro) History(ctx context.Context, id string) ([]model.ContactHistory, error) {
	var resp []model.ContactHistory

	fmt.Println(id, "id in storage")
	query := `
        SELECT 
            id,
            contact_id,
            phone,
            name,
            email,
            address,
            category,
			changed_at,
            change_type  
        FROM 
            contact_history
        WHERE 
            contact_id = $1
        ORDER BY 
		    changed_at DESC
        `

	// Execute the query
	rows, err := c.db.Query(ctx, query, id)
	if err != nil {
		return resp, err
	}
	defer rows.Close()

	// Iterate over the rows
	for rows.Next() {
		var contact model.ContactHistory
		var changedAt time.Time

		errScan := rows.Scan(
			&contact.ID,
			&contact.ContactID,
			&contact.Phone,
			&contact.Name,
			&contact.Email,
			&contact.Address,
			&contact.Category,
			&changedAt,
			&contact.ChangeType,
		)
		if errScan != nil {
			c.logger.Error("error while scanning contact history info: ", logger.Error(errScan))
			return resp, errScan
		}

		// Format changed_at field
		contact.Changed_at = changedAt.Format("2006-01-02 15:04:05")

		// Append contact history to response slice
		resp = append(resp, contact)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		c.logger.Error("error while iterating over rows: ", logger.Error(err))
		return resp, err
	}

	fmt.Printf("Final response: %+v\n", resp)
	return resp, nil
}
