package models

import (
	"afiqo-location/helpers"
	"afiqo-location/util"
	"context"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type (
	CustomerModel struct {
		ID          uuid.UUID
		Name        string
		Gender      int
		DateOfBirth time.Time
		Address     string
		PhoneNo     string
		Email       string
		Password    string
		IsActive    bool
		CreatedBy   uuid.UUID
		CreatedAt   time.Time
		UpdatedBy   uuid.NullUUID
		UpdatedAt   pq.NullTime
	}

	CustomerResponse struct {
		ID          uuid.UUID `json:"id"`
		Name        string    `json:"name"`
		Gender      string    `json:"gender"`
		DateOfBirth time.Time `json:"date_of_birth"`
		Address     string    `json:"address"`
		PhoneNo     string    `json:"phone_no"`
		Email       string    `json:"email"`
		IsActive    bool      `json:"is_active"`
		CreatedBy   uuid.UUID `json:"created_by"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedBy   uuid.UUID `json:"updated_by"`
		UpdatedAt   time.Time `json:"updated_at"`
	}
)

func (s CustomerModel) Response() CustomerResponse {

	gender := util.GetGender(s.Gender)

	return CustomerResponse{
		ID:          s.ID,
		Name:        s.Name,
		Gender:      gender,
		DateOfBirth: s.DateOfBirth,
		Address:     s.Address,
		PhoneNo:     s.PhoneNo,
		Email:       s.Email,
		IsActive:    s.IsActive,
		CreatedBy:   s.CreatedBy,
		CreatedAt:   s.CreatedAt,
		UpdatedBy:   s.UpdatedBy.UUID,
		UpdatedAt:   s.UpdatedAt.Time,
	}

}

func GetOneCustomer(ctx context.Context, db *sql.DB, customerID uuid.UUID) (CustomerModel, error) {

	query := fmt.Sprintf(`
		SELECT
			id,
			name,
			address,
			date_of_birth,	
			gender,
			email,
			password,
			phone_no,
			is_active,
			created_by,
			created_at,
			updated_by,
			updated_at
		FROM customer
		WHERE 
			id = $1
	`)

	var customer CustomerModel
	err := db.QueryRowContext(ctx, query, customerID).Scan(
		&customer.ID,
		&customer.Name,
		&customer.Address,
		&customer.DateOfBirth,
		&customer.Gender,
		&customer.Email,
		&customer.Password,
		&customer.PhoneNo,
		&customer.IsActive,
		&customer.CreatedBy,
		&customer.CreatedAt,
		&customer.UpdatedBy,
		&customer.UpdatedAt,
	)

	if err != nil {
		return CustomerModel{}, err
	}

	return customer, nil

}

func GetAllCustomer(ctx context.Context, db *sql.DB, filter helpers.Filter) ([]CustomerModel, error) {

	var searchQuery string

	if filter.Search != "" {
		searchQuery = fmt.Sprintf(`WHERE LOWER(name) LIKE LOWER('%%%s%%')`, filter.Search)
	}

	query := fmt.Sprintf(`
		SELECT
			id,
			name,
			address,
			date_of_birth,	
			gender,
			email,
			password,
			phone_no,
			is_active,
			created_by,
			created_at,
			updated_by,
			updated_at
		FROM customer
		%s
		ORDER BY name  %s
		LIMIT $1 OFFSET $2`, searchQuery, filter.Dir)

	rows, err := db.QueryContext(ctx, query, filter.Limit, filter.Offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var customers []CustomerModel
	for rows.Next() {
		var customer CustomerModel

		rows.Scan(
			&customer.ID,
			&customer.Name,
			&customer.Address,
			&customer.DateOfBirth,
			&customer.Gender,
			&customer.Email,
			&customer.Password,
			&customer.PhoneNo,
			&customer.IsActive,
			&customer.CreatedBy,
			&customer.CreatedAt,
			&customer.UpdatedBy,
			&customer.UpdatedAt,
		)

		customers = append(customers, customer)
	}

	return customers, nil

}

func GetOneCustomerByEmail(ctx context.Context, db *sql.DB, email string) (CustomerModel, error) {

	query := fmt.Sprintf(`
		SELECT
			id,
			name,
			address,
			date_of_birth,	
			gender,
			email,
			password,
			phone_no,
			is_active,
			created_by,
			created_at,
			updated_by,
			updated_at
		FROM customer
		WHERE 
			is_active = true AND email = $1 
	`)

	var customer CustomerModel
	err := db.QueryRowContext(ctx, query, email).Scan(
		&customer.ID,
		&customer.Name,
		&customer.Address,
		&customer.DateOfBirth,
		&customer.Gender,
		&customer.Email,
		&customer.Password,
		&customer.PhoneNo,
		&customer.IsActive,
		&customer.CreatedBy,
		&customer.CreatedAt,
		&customer.UpdatedBy,
		&customer.UpdatedAt,
	)

	if err != nil {
		return CustomerModel{}, err
	}

	return customer, nil

}

func (s *CustomerModel) Insert(ctx context.Context, db *sql.DB) error {

	password, err := bcrypt.GenerateFromPassword([]byte(s.Password), 12)

	if err != nil {
		return err
	}

	query := fmt.Sprintf(`
		INSERT INTO customer(
			name,
			address,
			date_of_birth,	
			gender,
			email,
			password,
			phone_no,
			created_by,
			created_at)
		VALUES(
		$1,$2,$3,$4,$5,$6,$7,$8,now())
		RETURNING id, created_at,is_active`)

	err = db.QueryRowContext(ctx, query,
		s.Name, s.Address, s.DateOfBirth, s.Gender, s.Email, password, s.PhoneNo, s.CreatedBy).Scan(
		&s.ID, &s.CreatedAt, &s.IsActive,
	)

	if err != nil {
		return err
	}

	return nil

}

func (s *CustomerModel) Update(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		UPDATE customer
		SET
			name=$1,
			address=$2,
			date_of_birth=$3,	
			gender=$4,
			phone_no=$5,
			updated_at=NOW(),
			updated_by=$6
		WHERE id=$7
		RETURNING id,created_at,updated_at,created_by,is_active,email`)

	err := db.QueryRowContext(ctx, query,
		s.Name, s.Address, s.DateOfBirth, s.Gender, s.PhoneNo, s.UpdatedBy, s.ID).Scan(
		&s.ID, &s.CreatedAt, &s.UpdatedAt, &s.CreatedBy, &s.IsActive, &s.Email,
	)

	if err != nil {
		return err
	}

	return nil

}

func (s *CustomerModel) PasswordUpdate(ctx context.Context, db *sql.DB) error {

	password, err := bcrypt.GenerateFromPassword([]byte(s.Password), 12)

	if err != nil {
		return err
	}

	query := fmt.Sprintf(`
		UPDATE customer
		SET
			password = $1,
			updated_at=NOW(),
			updated_by=$2
		WHERE id=$3
		RETURNING id,created_at,updated_at,created_by,is_active`)

	err = db.QueryRowContext(ctx, query,
		password, s.UpdatedBy, s.ID).Scan(
		&s.ID, &s.CreatedAt, &s.UpdatedAt, &s.CreatedBy, &s.IsActive,
	)

	if err != nil {
		return err
	}

	return nil

}

func (s *CustomerModel) Delete(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		UPDATE customer
		SET
			is_active=false,
			updated_by=$1,
			updated_at=NOW()
		WHERE id=$2`)

	_, err := db.ExecContext(ctx, query,
		s.UpdatedBy, s.ID)

	if err != nil {
		return err
	}

	return nil
}
