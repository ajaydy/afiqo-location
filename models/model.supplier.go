package models

import (
	"afiqo-location/helpers"
	"context"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type (
	SupplierModel struct {
		ID        uuid.UUID
		Name      string
		PhoneNo   string
		Email     string
		Password  string
		IsActive  bool
		CreatedBy uuid.UUID
		CreatedAt time.Time
		UpdatedBy uuid.NullUUID
		UpdatedAt pq.NullTime
	}

	SupplierResponse struct {
		ID        uuid.UUID `json:"id"`
		Name      string    `json:"name"`
		PhoneNo   string    `json:"phone_no"`
		Email     string    `json:"email"`
		IsActive  bool      `json:"is_active"`
		CreatedBy uuid.UUID `json:"created_by"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedBy uuid.UUID `json:"updated_by"`
		UpdatedAt time.Time `json:"updated_at"`
	}
)

func (s SupplierModel) Response() SupplierResponse {
	return SupplierResponse{
		ID:        s.ID,
		Name:      s.Name,
		PhoneNo:   s.PhoneNo,
		Email:     s.Email,
		IsActive:  s.IsActive,
		CreatedBy: s.CreatedBy,
		CreatedAt: s.CreatedAt,
		UpdatedBy: s.UpdatedBy.UUID,
		UpdatedAt: s.UpdatedAt.Time,
	}

}

func GetOneSupplier(ctx context.Context, db *sql.DB, supplierID uuid.UUID) (SupplierModel, error) {

	query := fmt.Sprintf(`
		SELECT
			id,
			name,
			phone_no,
			email,
			password,
			is_active,
			created_by,
			created_at,
			updated_by,
			updated_at
		FROM 
			supplier
		WHERE 
			id = $1
	`)

	var supplier SupplierModel
	err := db.QueryRowContext(ctx, query, supplierID).Scan(
		&supplier.ID,
		&supplier.Name,
		&supplier.PhoneNo,
		&supplier.Email,
		&supplier.Password,
		&supplier.IsActive,
		&supplier.CreatedBy,
		&supplier.CreatedAt,
		&supplier.UpdatedBy,
		&supplier.UpdatedAt,
	)

	if err != nil {
		return SupplierModel{}, err
	}

	return supplier, nil

}

func GetAllSupplier(ctx context.Context, db *sql.DB, filter helpers.Filter) ([]SupplierModel, error) {

	var searchQuery string

	if filter.Search != "" {
		searchQuery = fmt.Sprintf(`WHERE LOWER(name) LIKE LOWER('%%%s%%')`, filter.Search)
	}

	query := fmt.Sprintf(`
		SELECT
			id,
			name,
			phone_no,
			email,
			password,
			is_active,
			created_by,
			created_at,
			updated_by,
			updated_at
		FROM 
			supplier
		%s
		ORDER BY
			name %s
		LIMIT $1 OFFSET $2`,
		searchQuery, filter.Dir)

	rows, err := db.QueryContext(ctx, query, filter.Limit, filter.Offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var suppliers []SupplierModel
	for rows.Next() {
		var supplier SupplierModel

		rows.Scan(
			&supplier.ID,
			&supplier.Name,
			&supplier.PhoneNo,
			&supplier.Email,
			&supplier.Password,
			&supplier.IsActive,
			&supplier.CreatedBy,
			&supplier.CreatedAt,
			&supplier.UpdatedBy,
			&supplier.UpdatedAt,
		)

		suppliers = append(suppliers, supplier)
	}

	return suppliers, nil

}

func GetOneSupplierByEmail(ctx context.Context, db *sql.DB, email string) (SupplierModel, error) {

	query := fmt.Sprintf(`
		SELECT
			id,
			name,
			phone_no,
			email,
			password,
			is_active,
			created_by,
			created_at,
			updated_by,
			updated_at
		FROM 
			supplier
		WHERE 
			is_active = true
		AND 
			email = $1 
	`)

	var supplier SupplierModel
	err := db.QueryRowContext(ctx, query, email).Scan(
		&supplier.ID,
		&supplier.Name,
		&supplier.PhoneNo,
		&supplier.Email,
		&supplier.Password,
		&supplier.IsActive,
		&supplier.CreatedBy,
		&supplier.CreatedAt,
		&supplier.UpdatedBy,
		&supplier.UpdatedAt,
	)

	if err != nil {
		return SupplierModel{}, err
	}

	return supplier, nil

}

func (s *SupplierModel) Insert(ctx context.Context, db *sql.DB) error {

	password, err := bcrypt.GenerateFromPassword([]byte(s.Password), 12)

	if err != nil {
		return err
	}

	query := fmt.Sprintf(`
		INSERT INTO supplier(
			name,
			phone_no,
			email,
			password,
			created_by,
			created_at)
		VALUES(
			$1,$2,$3,$4,$5,now())
		RETURNING 
			id, created_at,is_active
	`)

	err = db.QueryRowContext(ctx, query,
		s.Name, s.PhoneNo, s.Email, password, s.CreatedBy).Scan(
		&s.ID, &s.CreatedAt, &s.IsActive,
	)

	if err != nil {
		return err
	}

	return nil

}

func (s *SupplierModel) Update(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		UPDATE supplier
		SET
			name=$1,
			phone_no=$2,
			updated_at=NOW(),
			updated_by=$3
		WHERE 
			id=$4
		RETURNING 
			id,created_at,updated_at,created_by,is_active,email
	`)

	err := db.QueryRowContext(ctx, query,
		s.Name, s.PhoneNo, s.UpdatedBy, s.ID).Scan(
		&s.ID, &s.CreatedAt, &s.UpdatedAt, &s.CreatedBy, &s.IsActive, &s.Email,
	)

	if err != nil {
		return err
	}

	return nil

}

func (s *SupplierModel) PasswordUpdate(ctx context.Context, db *sql.DB) error {

	password, err := bcrypt.GenerateFromPassword([]byte(s.Password), 12)

	if err != nil {
		return err
	}

	query := fmt.Sprintf(`
		UPDATE supplier
		SET
			password = $1,
			updated_at=NOW(),
			updated_by=$2
		WHERE 
			id=$3
		RETURNING 
			id,created_at,updated_at,created_by,is_active
	`)

	err = db.QueryRowContext(ctx, query,
		password, s.UpdatedBy, s.ID).Scan(
		&s.ID, &s.CreatedAt, &s.UpdatedAt, &s.CreatedBy, &s.IsActive,
	)

	if err != nil {
		return err
	}

	return nil

}

func (s *SupplierModel) Delete(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		UPDATE supplier
		SET
			is_active=false,
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
