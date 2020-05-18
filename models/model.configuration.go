package models

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"time"
)

type (
	ConfigurationModel struct {
		ID          uuid.UUID
		DeliveryFee decimal.Decimal
		IsDelete    bool
		CreatedBy   uuid.UUID
		CreatedAt   time.Time
		UpdatedBy   uuid.NullUUID
		UpdatedAt   pq.NullTime
	}

	ConfigurationResponse struct {
		ID          uuid.UUID       `json:"id"`
		DeliveryFee decimal.Decimal `json:"delivery_fee"`
		IsDelete    bool            `json:"is_delete"`
		CreatedBy   uuid.UUID       `json:"created_by"`
		CreatedAt   time.Time       `json:"created_at"`
		UpdatedBy   uuid.UUID       `json:"updated_by"`
		UpdatedAt   time.Time       `json:"updated_at"`
	}
)

func (s ConfigurationModel) Response() ConfigurationResponse {
	return ConfigurationResponse{
		ID:          s.ID,
		DeliveryFee: s.DeliveryFee,
		IsDelete:    s.IsDelete,
		CreatedBy:   s.CreatedBy,
		CreatedAt:   s.CreatedAt,
		UpdatedBy:   s.UpdatedBy.UUID,
		UpdatedAt:   s.UpdatedAt.Time,
	}
}

func GetConfiguration(ctx context.Context, db *sql.DB) (ConfigurationModel, error) {

	query := fmt.Sprintf(`
		SELECT
			id,
			delivery_fee,
			is_delete,
			created_by,
			created_at,
			updated_by,
			updated_at
		FROM 
			configuration
	`)

	var configuration ConfigurationModel
	err := db.QueryRowContext(ctx, query).Scan(
		&configuration.ID,
		&configuration.DeliveryFee,
		&configuration.IsDelete,
		&configuration.CreatedBy,
		&configuration.CreatedAt,
		&configuration.UpdatedBy,
		&configuration.UpdatedAt,
	)

	if err != nil {
		return ConfigurationModel{}, err
	}

	return configuration, nil

}

func (s *ConfigurationModel) Update(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		UPDATE configuration
		SET
			delivery_fee=$1
			updated_at=NOW(),
			updated_by=$2
		`)

	db.QueryRowContext(ctx, query,
		s.DeliveryFee, s.UpdatedBy)

	return nil

}
