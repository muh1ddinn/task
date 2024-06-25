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

type categoriesRepro struct {
	db     *pgxpool.Pool
	logger logger.ILogger
}

func NewCategories(db *pgxpool.Pool, log logger.ILogger) categoriesRepro {
	return categoriesRepro{
		db:     db,
		logger: log,
	}
}

func (c *categoriesRepro) Createcat(ctx context.Context, categories model.Categories) (model.Getcategoriest, error) {

	id := uuid.New()

	query := `INSERT INTO categories(
		id,
		name,
		created_at
	) VALUES($1, $2, $3)`

	_, err := c.db.Exec(ctx, query,
		id.String(),
		categories.Name,
		time.Now(),
	)
	fmt.Println(err)
	if err != nil {
		c.logger.Error("error while creating ategories:", logger.Error(err))
		return model.Getcategoriest{}, err
	}

	a, err := c.GetByIDcat(ctx, id.String())
	if err != nil {
		log.Println("error while getting ategories by id")
		return a, err
	}

	return a, nil
}

func (c *categoriesRepro) GetByIDcat(ctx context.Context, id string) (model.Getcategoriest, error) {

	categories := model.Getcategoriest{}

	query := `SELECT 
	id,
	name
	FROM categories WHERE id=$1 AND deleted_at =0`

	row := c.db.QueryRow(ctx, query, id)

	err := row.Scan(
		&categories.Id,
		&categories.Name,
	)
	if err != nil {
		fmt.Println(err)
		c.logger.Error("error while getting categories by ID: ", logger.Error(err))
		return categories, err
	}

	return categories, nil
}

func (c *categoriesRepro) SoftDeletecat(ctx context.Context, id string) (string, error) {
	query := `UPDATE categories SET 
        deleted_at = 1
        WHERE id = $1 AND deleted_at=0`

	_, err := c.db.Exec(ctx, query, id)

	fmt.Println(err)
	if err != nil {
		c.logger.Error("failed to delete user from database", logger.Error(err))
		return id, err
	}

	return id, nil
}

func (c *categoriesRepro) Deletecat(ctx context.Context, id string) (string, error) {
	query := `DELETE FROM categories
        WHERE id = $1`

	_, err := c.db.Exec(ctx, query, id)
	fmt.Println(err)
	if err != nil {
		c.logger.Error("failed to delete user from database", logger.Error(err))
		return id, err
	}

	return id, nil
}

func (c *categoriesRepro) Patchcat(ctx context.Context, categories model.Patchcategories) (model.Getcategoriest, error) {
	fmt.Println(categories.Name, "name")

	query := `UPDATE categories SET 
	id=$1,
	name= $2,
	updated_at = CURRENT_TIMESTAMP
	WHERE id = $1 AND deleted_at =0`
	_, err := c.db.Exec(ctx, query, categories.Id,
		categories.Name)

	fmt.Println(err, "error")
	if err != nil {
		c.logger.Error("failed to pacth user from database", logger.Error(err))
		return model.Getcategoriest{}, err
	}
	a, err := c.GetByIDcat(ctx, categories.Id)
	if err != nil {
		log.Println("error while getting category by id")
		return a, err
	}

	return a, nil
}

func (c *categoriesRepro) GetAllcat(ctx context.Context, req model.GetAllCategoriestRequest) (model.GetAllcategoriesResponse, error) {
	resp := model.GetAllcategoriesResponse{}

	offset := (req.Page - 1) * req.Limit
	baseFilter := " WHERE deleted_at = 0"
	searchFilter := ""
	args := []interface{}{}
	argIndex := 1

	if req.Search != "" {
		searchTerm := fmt.Sprintf("%%%s%%", req.Search)
		searchFilter = ` AND (name ILIKE $` + fmt.Sprint(argIndex) + `)`
		args = append(args, searchTerm)
		argIndex++
	}

	args = append(args, offset, req.Limit)

	query := `
        SELECT 
            count(id) OVER() AS total_count,
            id,
            name,
            created_at
        FROM 
            categories` + baseFilter + searchFilter + `
        ORDER BY 
            created_at DESC
        OFFSET $` + fmt.Sprint(argIndex) + ` LIMIT $` + fmt.Sprint(argIndex+1)

	fmt.Println("query: ", query)
	fmt.Println("args: ", args)

	rows, err := c.db.Query(ctx, query, args...)
	if err != nil {
		return resp, err
	}
	defer rows.Close()

	for rows.Next() {
		var category model.Getcategoriest
		var createdAt time.Time

		errScan := rows.Scan(
			&resp.Count,
			&category.Id,
			&category.Name,
			&createdAt,
		)
		if errScan != nil {
			c.logger.Error("error while scanning category info: ", logger.Error(errScan))
			return resp, errScan
		}

		category.Created_at = createdAt.Format("2006-01-02 15:04:05")
		resp.Categories = append(resp.Categories, category)
	}

	if err := rows.Err(); err != nil {
		c.logger.Error("error while iterating over rows: ", logger.Error(err))
		return model.GetAllcategoriesResponse{}, err
	}

	fmt.Printf("Final response: %+v\n", resp)
	return resp, nil
}

func (c categoriesRepro) Checkname(ctx context.Context, name string) (string, error) {

	contact := model.GetAllContact{}
	fmt.Println(name)
	query := `SELECT 
		name
	FROM categories WHERE name=$1 AND deleted_at = 0`

	row := c.db.QueryRow(ctx, query, name)

	err := row.Scan(
		&contact.Email,
	)
	if err != nil {
		if err.Error() == "no rows in result set" {
			fmt.Println("No contact found with name:", name)
			return "", nil
		}
		fmt.Println(err)
		c.logger.Error("error while getting contact by name: ", logger.Error(err))
		return "", err
	}

	return contact.Email, nil
}
