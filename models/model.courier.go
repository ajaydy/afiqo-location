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
	CourierModel struct {
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

	CourierResponse struct {
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

func (s CourierModel) Response() CourierResponse {
	return CourierResponse{
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

func GetOneCourier(ctx context.Context, db *sql.DB, courierID uuid.UUID) (CourierModel, error) {

	query := fmt.Sprintf(`
		SELECT
			id,
			name,
			email,
			password,
			phone_no,
			is_active,
			created_by,
			created_at,
			updated_by,
			updated_at
		FROM courier
		WHERE 
			id = $1
	`)

	var courier CourierModel
	err := db.QueryRowContext(ctx, query, courierID).Scan(
		&courier.ID,
		&courier.Name,
		&courier.Email,
		&courier.Password,
		&courier.PhoneNo,
		&courier.IsActive,
		&courier.CreatedBy,
		&courier.CreatedAt,
		&courier.UpdatedBy,
		&courier.UpdatedAt,
	)

	if err != nil {
		return CourierModel{}, err
	}

	return courier, nil

}

func GetAllCourier(ctx context.Context, db *sql.DB, filter helpers.Filter) ([]CourierModel, error) {

	var searchQuery string

	if filter.Search != "" {
		searchQuery = fmt.Sprintf(`WHERE LOWER(name) LIKE LOWER('%%%s%%')`, filter.Search)
	}

	query := fmt.Sprintf(`
		SELECT
			id,
			name,
			email,
			password,
			phone_no,
			is_active,
			created_by,
			created_at,
			updated_by,
			updated_at
		FROM courier
		%s
		ORDER BY name  %s
		LIMIT $1 OFFSET $2`, searchQuery, filter.Dir)

	rows, err := db.QueryContext(ctx, query, filter.Limit, filter.Offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var couriers []CourierModel
	for rows.Next() {
		var courier CourierModel

		rows.Scan(
			&courier.ID,
			&courier.Name,
			&courier.Email,
			&courier.Password,
			&courier.PhoneNo,
			&courier.IsActive,
			&courier.CreatedBy,
			&courier.CreatedAt,
			&courier.UpdatedBy,
			&courier.UpdatedAt,
		)

		couriers = append(couriers, courier)
	}

	return couriers, nil

}

func GetOneCourierByEmail(ctx context.Context, db *sql.DB, email string) (CourierModel, error) {

	query := fmt.Sprintf(`
		SELECT
			id,
			name,
			email,
			password,
			phone_no,
			is_active,
			created_by,
			created_at,
			updated_by,
			updated_at
		FROM courier
		WHERE 
			is_active = true AND email = $1 
	`)

	var courier CourierModel
	err := db.QueryRowContext(ctx, query, email).Scan(
		&courier.ID,
		&courier.Name,
		&courier.Email,
		&courier.Password,
		&courier.PhoneNo,
		&courier.IsActive,
		&courier.CreatedBy,
		&courier.CreatedAt,
		&courier.UpdatedBy,
		&courier.UpdatedAt,
	)

	if err != nil {
		return CourierModel{}, err
	}

	return courier, nil

}

func (s *CourierModel) Insert(ctx context.Context, db *sql.DB) error {

	password, err := bcrypt.GenerateFromPassword([]byte(s.Password), 12)

	if err != nil {
		return err
	}

	query := fmt.Sprintf(`
		INSERT INTO courier(
			name,
			email,
			password,
			phone_no,
			created_by,
			created_at)
		VALUES(
		$1,$2,$3,$4,$5,now())
		RETURNING id, created_at,is_active`)

	err = db.QueryRowContext(ctx, query,
		s.Name, s.Email, password, s.PhoneNo, s.CreatedBy).Scan(
		&s.ID, &s.CreatedAt, &s.IsActive,
	)

	if err != nil {
		return err
	}

	return nil

}

func (s *CourierModel) Update(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		UPDATE courier
		SET
			name=$1,
			phone_no=$2,
			updated_at=NOW(),
			updated_by=$3
		WHERE id=$4
		RETURNING id,created_at,updated_at,created_by,is_active,email`)

	err := db.QueryRowContext(ctx, query,
		s.Name, s.PhoneNo, s.UpdatedBy, s.ID).Scan(
		&s.ID, &s.CreatedAt, &s.UpdatedAt, &s.CreatedBy, &s.IsActive, &s.Email,
	)

	if err != nil {
		return err
	}

	return nil

}

func (s *CourierModel) PasswordUpdate(ctx context.Context, db *sql.DB) error {

	password, err := bcrypt.GenerateFromPassword([]byte(s.Password), 12)

	if err != nil {
		return err
	}

	query := fmt.Sprintf(`
		UPDATE courier
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

func (s *CourierModel) Delete(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		UPDATE courier
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
