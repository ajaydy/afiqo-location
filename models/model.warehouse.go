package models

import (
	"afiqo-location/helpers"
	"context"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"time"
)

type (
	WarehouseModel struct {
		ID        uuid.UUID
		Name      string
		Address   string
		Latitude  decimal.Decimal
		Longitude decimal.Decimal
		PhoneNo   string
		IsDelete  bool
		CreatedBy uuid.UUID
		CreatedAt time.Time
		UpdatedBy uuid.NullUUID
		UpdatedAt pq.NullTime
	}

	WarehouseResponse struct {
		ID        uuid.UUID       `json:"id"`
		Name      string          `json:"name"`
		Address   string          `json:"address"`
		Latitude  decimal.Decimal `json:"latitude"`
		Longitude decimal.Decimal `json:"longitude"`
		PhoneNo   string          `json:"phone_no"`
		IsDelete  bool            `json:"is_delete"`
		CreatedBy uuid.UUID       `json:"created_by"`
		CreatedAt time.Time       `json:"created_at"`
		UpdatedBy uuid.UUID       `json:"updated_by"`
		UpdatedAt time.Time       `json:"updated_at"`
	}
)

func (s WarehouseModel) Response() WarehouseResponse {
	return WarehouseResponse{
		ID:        s.ID,
		Name:      s.Name,
		Address:   s.Address,
		Latitude:  s.Latitude,
		Longitude: s.Longitude,
		PhoneNo:   s.PhoneNo,
		IsDelete:  s.IsDelete,
		CreatedBy: s.CreatedBy,
		CreatedAt: s.CreatedAt,
		UpdatedBy: s.UpdatedBy.UUID,
		UpdatedAt: s.UpdatedAt.Time,
	}
}

func GetOneWarehouse(ctx context.Context, db *sql.DB, warehouseID uuid.UUID) (WarehouseModel, error) {

	query := fmt.Sprintf(`
		SELECT
			id,
			name,
			address,
			latitude,
			longitude,
			phone_no,
			is_delete,
			created_by,
			created_at,
			updated_by,
			updated_at
		FROM 
			warehouse
		WHERE 
			id = $1
	`)

	var warehouse WarehouseModel
	err := db.QueryRowContext(ctx, query, warehouseID).Scan(
		&warehouse.ID,
		&warehouse.Name,
		&warehouse.Address,
		&warehouse.Latitude,
		&warehouse.Longitude,
		&warehouse.PhoneNo,
		&warehouse.IsDelete,
		&warehouse.CreatedBy,
		&warehouse.CreatedAt,
		&warehouse.UpdatedBy,
		&warehouse.UpdatedAt,
	)

	if err != nil {
		return WarehouseModel{}, err
	}

	return warehouse, nil

}

func GetAllWarehouseWithDistance(ctx context.Context, db *sql.DB, filter helpers.Filter) ([]WarehouseModel, error) {

	query := fmt.Sprintf(`
		SELECT
			id,
			name,
			address,
			latitude,
			longitude,
			phone_no,
			is_delete,
			created_by,
			created_at,
			updated_by,
			updated_at,
			SQRT(
				POW(69.1 * (latitude::FLOAT8 - $1), 2) +
				POW(69.1 * ($2 - longitude::FLOAT8) * COS(latitude::FLOAT8 / 57.3), 2)) AS distance  
		FROM 
			warehouse
		ORDER BY 
			distance
		LIMIT $3 OFFSET $4`)

	rows, err := db.QueryContext(ctx, query, filter.Latitude, filter.Longitude, filter.Limit, filter.Offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var warehouses []WarehouseModel
	for rows.Next() {
		var warehouse WarehouseModel
		var distance decimal.Decimal

		rows.Scan(
			&warehouse.ID,
			&warehouse.Name,
			&warehouse.Address,
			&warehouse.Latitude,
			&warehouse.Longitude,
			&warehouse.PhoneNo,
			&warehouse.IsDelete,
			&warehouse.CreatedBy,
			&warehouse.CreatedAt,
			&warehouse.UpdatedBy,
			&warehouse.UpdatedAt,
			&distance,
		)

		warehouses = append(warehouses, warehouse)

	}

	return warehouses, nil

}

func GetAllWarehouse(ctx context.Context, db *sql.DB, filter helpers.Filter) ([]WarehouseModel, error) {

	var searchQuery string

	if filter.Search != "" {
		searchQuery = fmt.Sprintf(`WHERE LOWER(name) LIKE LOWER('%%%s%%')`, filter.Search)
	}

	query := fmt.Sprintf(`
		SELECT
			id,
			name,
			address,
			latitude,
			longitude,
			phone_no,
			is_delete,
			created_by,
			created_at,
			updated_by,
			updated_at
		FROM 
			warehouse
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

	var warehouses []WarehouseModel
	for rows.Next() {
		var warehouse WarehouseModel

		rows.Scan(
			&warehouse.ID,
			&warehouse.Name,
			&warehouse.Address,
			&warehouse.Latitude,
			&warehouse.Longitude,
			&warehouse.PhoneNo,
			&warehouse.IsDelete,
			&warehouse.CreatedBy,
			&warehouse.CreatedAt,
			&warehouse.UpdatedBy,
			&warehouse.UpdatedAt,
		)

		warehouses = append(warehouses, warehouse)
	}

	return warehouses, nil

}

func (s *WarehouseModel) Insert(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		INSERT INTO warehouse(
			name,
			address,
			latitude,
			longitude,
			phone_no,
			created_by,
			created_at)
		VALUES(
			$1,$2,$3,$4,$5,$6,now())
		RETURNING 
			id, created_at
	`)

	err := db.QueryRowContext(ctx, query,
		s.Name, s.Address, s.Latitude, s.Longitude, s.PhoneNo, s.CreatedBy).Scan(
		&s.ID, &s.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil

}

func (s *WarehouseModel) Update(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		UPDATE warehouse
		SET
			name=$1,
			address=$2,
			latitude=$3,
			longitude=$4,
			phone_no=$5,
			updated_at=NOW(),
			updated_by=$6
		WHERE 
			id=$7
		RETURNING 
			id,created_at,updated_at,created_by
	`)

	err := db.QueryRowContext(ctx, query,
		s.Name, s.Address, s.Latitude, s.Longitude, s.PhoneNo, s.UpdatedBy, s.ID).Scan(
		&s.ID, &s.CreatedAt, &s.UpdatedAt, &s.CreatedBy,
	)

	if err != nil {
		return err
	}

	return nil

}

func (s *WarehouseModel) Delete(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		UPDATE warehouse
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
