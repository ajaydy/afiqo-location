package models

import (
	"afiqo-location/helpers"
	"context"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	"strings"
	"time"
)

type (
	WarehouseProductModel struct {
		ID          uuid.UUID
		ProductID   uuid.UUID
		WarehouseID uuid.UUID
		Stock       uint
		IsDelete    bool
		CreatedBy   uuid.UUID
		CreatedAt   time.Time
		UpdatedBy   uuid.NullUUID
		UpdatedAt   pq.NullTime
	}

	WarehouseProductResponse struct {
		ID        uuid.UUID         `json:"id"`
		Product   ProductResponse   `json:"product"`
		Warehouse WarehouseResponse `json:"warehouse"`
		Stock     uint              `json:"stock"`
		IsDelete  bool              `json:"is_delete"`
		CreatedBy uuid.UUID         `json:"created_by"`
		CreatedAt time.Time         `json:"created_at"`
		UpdatedBy uuid.UUID         `json:"updated_by"`
		UpdatedAt time.Time         `json:"updated_at"`
	}
)

func (s WarehouseProductModel) Response(ctx context.Context, db *sql.DB, logger *helpers.Logger) (
	WarehouseProductResponse, error) {

	product, err := GetOneProduct(ctx, db, s.ProductID)
	if err != nil {
		logger.Err.Printf(`model.warehouse.product.go/GetOneProduct/%v`, err)
		return WarehouseProductResponse{}, nil
	}

	productResponse, err := product.Response(ctx, db, logger)

	if err != nil {
		logger.Err.Printf(`model.warehouse.product.go/product.Response/%v`, err)
		return WarehouseProductResponse{}, nil
	}

	warehouse, err := GetOneWarehouse(ctx, db, s.WarehouseID)
	if err != nil {
		logger.Err.Printf(`model.warehouse.product.go/GetOneWarehouse/%v`, err)
		return WarehouseProductResponse{}, nil
	}

	return WarehouseProductResponse{
		ID:        s.ID,
		Product:   productResponse,
		Warehouse: warehouse.Response(),
		Stock:     s.Stock,
		IsDelete:  s.IsDelete,
		CreatedBy: s.CreatedBy,
		CreatedAt: s.CreatedAt,
		UpdatedBy: s.UpdatedBy.UUID,
		UpdatedAt: s.UpdatedAt.Time,
	}, nil

}

func GetOneWarehouseProduct(ctx context.Context, db *sql.DB, warehouseProductID uuid.UUID) (
	WarehouseProductModel, error) {

	query := fmt.Sprintf(`
		SELECT
			id,
			warehouse_id,
			product_id,
			stock,
			is_delete,
			created_by,
			created_at,
			updated_by,
			updated_at
		FROM warehouse_product
		WHERE 
			id = $1
	`)

	var warehouse WarehouseProductModel
	err := db.QueryRowContext(ctx, query, warehouseProductID).Scan(
		&warehouse.ID,
		&warehouse.WarehouseID,
		&warehouse.ProductID,
		&warehouse.Stock,
		&warehouse.IsDelete,
		&warehouse.CreatedBy,
		&warehouse.CreatedAt,
		&warehouse.UpdatedBy,
		&warehouse.UpdatedAt,
	)

	if err != nil {
		return WarehouseProductModel{}, err
	}

	return warehouse, nil

}

func GetAllWarehouseProduct(ctx context.Context, db *sql.DB, filter helpers.Filter) (
	[]WarehouseProductModel, error) {

	var filters []string

	if filter.WarehouseID != uuid.Nil {
		filters = append(filters, fmt.Sprintf(`
			warehouse_id = '%s'`,
			filter.WarehouseID))
	}

	if filter.ProductID != uuid.Nil {
		filters = append(filters, fmt.Sprintf(`
			product_id = '%s'`,
			filter.WarehouseID))
	}

	filterJoin := strings.Join(filters, " AND ")
	if filterJoin != "" {
		filterJoin = fmt.Sprintf("WHERE %s", filterJoin)
	}

	query := fmt.Sprintf(`
		SELECT
			id,
			warehouse_id,
			product_id,
			stock,
			is_delete,
			created_by,
			created_at,
			updated_by,
			updated_at
		FROM warehouse_product
		LIMIT $1 OFFSET $2`)

	rows, err := db.QueryContext(ctx, query, filter.Limit, filter.Offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var warehouses []WarehouseProductModel
	for rows.Next() {
		var warehouse WarehouseProductModel

		rows.Scan(
			&warehouse.ID,
			&warehouse.WarehouseID,
			&warehouse.ProductID,
			&warehouse.Stock,
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

func (s *WarehouseProductModel) Insert(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		INSERT INTO warehouse_product(
			warehouse_id,
			product_id,
			stock,
			created_by,
			created_at)
		VALUES(
		$1,$2,$3,$4,now())
		RETURNING id, created_at`)

	err := db.QueryRowContext(ctx, query,
		s.WarehouseID, s.ProductID, s.Stock, s.CreatedBy).Scan(
		&s.ID, &s.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil

}
func (s *WarehouseProductModel) Update(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		UPDATE warehouse_product
		SET
			stock = $1
			updated_at=NOW(),
			updated_by=$2
		WHERE id=$3
		RETURNING id,created_at,updated_at,created_by`)

	err := db.QueryRowContext(ctx, query,
		s.Stock, s.UpdatedBy, s.ID).Scan(
		&s.ID, &s.CreatedAt, &s.UpdatedAt, &s.CreatedBy,
	)

	if err != nil {
		return err
	}

	return nil

}

func (s *WarehouseProductModel) Delete(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		UPDATE warehouse_product
		SET
			is_delete=true,
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
