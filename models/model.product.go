package models

import (
	"afiqo-location/helpers"
	"context"
	"database/sql"
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"time"
)

type (
	ProductModel struct {
		ID          uuid.UUID
		SupplierID  uuid.UUID
		CategoryID  uuid.UUID
		Name        string
		Stock       uint
		Price       decimal.Decimal
		Description string
		IsDelete    bool
		CreatedBy   uuid.UUID
		CreatedAt   time.Time
		UpdatedBy   uuid.NullUUID
		UpdatedAt   pq.NullTime
	}

	ProductResponse struct {
		ID          uuid.UUID        `json:"id"`
		Supplier    SupplierResponse `json:"supplier"`
		Category    CategoryResponse `json:"category"`
		Name        string           `json:"name"`
		Stock       uint             `json:"stock"`
		Price       decimal.Decimal  `json:"price"`
		Description string           `json:"description"`
		CreatedBy   uuid.UUID        `json:"created_by"`
		CreatedAt   time.Time        `json:"created_at"`
		UpdatedBy   uuid.UUID        `json:"updated_by"`
		UpdatedAt   time.Time        `json:"updated_at"`
	}
)

func (s ProductModel) Response(ctx context.Context, db *sql.DB, logger *helpers.Logger) (ProductResponse, error) {

	supplier, err := GetOneSupplier(ctx, db, s.SupplierID)
	if err != nil {
		logger.Err.Printf(`model.product.go/GetOneProduct/%v`, err)
		return ProductResponse{}, nil
	}

	category, err := GetOneCategory(ctx, db, s.CategoryID)
	if err != nil {
		logger.Err.Printf(`model.product.go/GetOneCategory/%v`, err)
		return ProductResponse{}, nil
	}

	return ProductResponse{
		ID:          s.ID,
		Supplier:    supplier.Response(),
		Category:    category.Response(),
		Name:        s.Name,
		Stock:       s.Stock,
		Price:       s.Price,
		Description: s.Description,
		CreatedBy:   s.CreatedBy,
		CreatedAt:   s.CreatedAt,
		UpdatedBy:   s.UpdatedBy.UUID,
		UpdatedAt:   s.UpdatedAt.Time,
	}, nil
}
