package models

import (
	"afiqo-location/helpers"
	"context"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	"time"
)

type (
	CategoryModel struct {
		ID          uuid.UUID
		Name        string
		Description string
		IsDelete    bool
		CreatedBy   uuid.UUID
		CreatedAt   time.Time
		UpdatedBy   uuid.NullUUID
		UpdatedAt   pq.NullTime
	}

	CategoryResponse struct {
		ID          uuid.UUID `json:"id"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		IsDelete    bool      `json:"is_delete"`
		CreatedBy   uuid.UUID `json:"created_by"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedBy   uuid.UUID `json:"updated_by"`
		UpdatedAt   time.Time `json:"updated_at"`
	}
)

func (s CategoryModel) Response() CategoryResponse {
	return CategoryResponse{
		ID:          s.ID,
		Name:        s.Name,
		Description: s.Description,
		IsDelete:    s.IsDelete,
		CreatedBy:   s.CreatedBy,
		CreatedAt:   s.CreatedAt,
		UpdatedBy:   s.UpdatedBy.UUID,
		UpdatedAt:   s.UpdatedAt.Time,
	}
}

func GetOneCategory(ctx context.Context, db *sql.DB, categoryID uuid.UUID) (CategoryModel, error) {

	query := fmt.Sprintf(`
		SELECT
			id,
			name,
			description,
			is_delete,
			created_by,
			created_at,
			updated_by,
			updated_at
		FROM category
		WHERE 
			id = $1
	`)

	var category CategoryModel
	err := db.QueryRowContext(ctx, query, categoryID).Scan(
		&category.ID,
		&category.Name,
		&category.Description,
		&category.IsDelete,
		&category.CreatedBy,
		&category.CreatedAt,
		&category.UpdatedBy,
		&category.UpdatedAt,
	)

	if err != nil {
		return CategoryModel{}, err
	}

	return category, nil

}

func GetAllCategory(ctx context.Context, db *sql.DB, filter helpers.Filter) ([]CategoryModel, error) {

	var searchQuery string

	if filter.Search != "" {
		searchQuery = fmt.Sprintf(`AND LOWER(name) LIKE LOWER('%%%s%%')`, filter.Search)
	}

	query := fmt.Sprintf(`
		SELECT
			id,
			name,
			description,
			is_delete,
			created_by,
			created_at,
			updated_by,
			updated_at
		FROM category
		WHERE 
			is_delete = false
		%s 
		ORDER BY
			name  %s
		LIMIT $1 OFFSET $2`,
		searchQuery, filter.Dir)

	rows, err := db.QueryContext(ctx, query, filter.Limit, filter.Offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var categories []CategoryModel
	for rows.Next() {
		var category CategoryModel

		rows.Scan(
			&category.ID,
			&category.Name,
			&category.Description,
			&category.IsDelete,
			&category.CreatedBy,
			&category.CreatedAt,
			&category.UpdatedBy,
			&category.UpdatedAt,
		)

		categories = append(categories, category)
	}

	return categories, nil

}

func (s *CategoryModel) Insert(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		INSERT INTO category(
			name,
			description,
			created_by,
			created_at
		)VALUES(
			$1,$2,$3,now())
		RETURNING
			id, created_at
	`)

	err := db.QueryRowContext(ctx, query,
		s.Name, s.Description, s.CreatedBy).Scan(
		&s.ID, &s.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil

}

func (s *CategoryModel) Update(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		UPDATE category
		SET
			name=$1,
			description=$2,
			updated_at=NOW(),
			updated_by=$3
		WHERE 
			id=$4
		RETURNING 
			id,created_at,updated_at,created_by,is_delete
	`)

	err := db.QueryRowContext(ctx, query,
		s.Name, s.Description, s.UpdatedBy, s.ID).Scan(
		&s.ID, &s.CreatedAt, &s.UpdatedAt, &s.CreatedBy, &s.IsDelete,
	)

	if err != nil {
		return err
	}

	return nil

}

func (s *CategoryModel) Delete(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		UPDATE category
		SET
			is_delete=true,
			updated_by=$1,
			updated_at=NOW()
		WHERE 
			id=$2
	`)

	_, err := db.ExecContext(ctx, query,
		s.UpdatedBy, s.ID)

	if err != nil {
		return err
	}

	return nil
}
