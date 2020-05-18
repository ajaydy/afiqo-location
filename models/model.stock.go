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
	StockModel struct {
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

	StockResponse struct {
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

func (s StockModel) Response(ctx context.Context, db *sql.DB, logger *helpers.Logger) (
	StockResponse, error) {

	product, err := GetOneProduct(ctx, db, s.ProductID)
	if err != nil {
		logger.Err.Printf(`model.stock.go/GetOneProduct/%v`, err)
		return StockResponse{}, nil
	}

	productResponse, err := product.Response(ctx, db, logger)

	if err != nil {
		logger.Err.Printf(`model.stock.go/product.Response/%v`, err)
		return StockResponse{}, nil
	}

	warehouse, err := GetOneWarehouse(ctx, db, s.WarehouseID)
	if err != nil {
		logger.Err.Printf(`model.stock.go/GetOneWarehouse/%v`, err)
		return StockResponse{}, nil
	}

	return StockResponse{
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

func GetOneStock(ctx context.Context, db *sql.DB, stockID uuid.UUID) (
	StockModel, error) {

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
		FROM 
			stock
		WHERE 
			id = $1
	`)

	var stock StockModel
	err := db.QueryRowContext(ctx, query, stockID).Scan(
		&stock.ID,
		&stock.WarehouseID,
		&stock.ProductID,
		&stock.Stock,
		&stock.IsDelete,
		&stock.CreatedBy,
		&stock.CreatedAt,
		&stock.UpdatedBy,
		&stock.UpdatedAt,
	)

	if err != nil {
		return StockModel{}, err
	}

	return stock, nil

}

func GetOneStockByProductAndWarehouse(ctx context.Context, db *sql.DB, warehouseID, productID uuid.UUID) (
	StockModel, error) {

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
		FROM 
			stock
		WHERE 
			warehouse_id=$1
		AND 
			product_id=$2
	`)

	var stock StockModel
	err := db.QueryRowContext(ctx, query, warehouseID, productID).Scan(
		&stock.ID,
		&stock.WarehouseID,
		&stock.ProductID,
		&stock.Stock,
		&stock.IsDelete,
		&stock.CreatedBy,
		&stock.CreatedAt,
		&stock.UpdatedBy,
		&stock.UpdatedAt,
	)

	if err != nil {
		return StockModel{}, err
	}

	return stock, nil

}

func GetAllStock(ctx context.Context, db *sql.DB, filter helpers.Filter) (
	[]StockModel, error) {

	var filters []string

	if filter.WarehouseID != uuid.Nil {
		filters = append(filters, fmt.Sprintf(`
			warehouse_id = '%s'`,
			filter.WarehouseID))
	}

	if filter.ProductID != uuid.Nil {
		filters = append(filters, fmt.Sprintf(`
			product_id = '%s'`,
			filter.ProductID))
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
		FROM 
			stock
		%s
		LIMIT $1 OFFSET $2`,
		filterJoin)

	rows, err := db.QueryContext(ctx, query, filter.Limit, filter.Offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var stocks []StockModel
	for rows.Next() {
		var stock StockModel

		rows.Scan(
			&stock.ID,
			&stock.WarehouseID,
			&stock.ProductID,
			&stock.Stock,
			&stock.IsDelete,
			&stock.CreatedBy,
			&stock.CreatedAt,
			&stock.UpdatedBy,
			&stock.UpdatedAt,
		)

		stocks = append(stocks, stock)
	}

	return stocks, nil

}

func GetAllStockBySupplierID(ctx context.Context, db *sql.DB, filter helpers.Filter, supplierID uuid.UUID) (
	[]StockModel, error) {

	var filters []string

	if filter.WarehouseID != uuid.Nil {
		filters = append(filters, fmt.Sprintf(`
			warehouse_id = '%s'`,
			filter.WarehouseID))
	}

	if filter.ProductID != uuid.Nil {
		filters = append(filters, fmt.Sprintf(`
			product_id = '%s'`,
			filter.ProductID))
	}

	filterJoin := strings.Join(filters, " AND ")
	if filterJoin != "" {
		filterJoin = fmt.Sprintf("AND %s", filterJoin)
	}

	query := fmt.Sprintf(`
		SELECT
			s.id,
			s.warehouse_id,
			s.product_id,
			s.stock,
			s.is_delete,
			s.created_by,
			s.created_at,
			s.updated_by,
			s.updated_at
		FROM 
			stock s
		INNER JOIN
			product p
		ON
			p.id = s.product_id
		WHERE 
			p.supplier_id = $1
		%s
		ORDER BY
			s.updated_at %s ,s.created_at %s
		LIMIT $2 OFFSET $3`,
		filterJoin, filter.Dir, filter.Dir)

	rows, err := db.QueryContext(ctx, query, supplierID, filter.Limit, filter.Offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var stocks []StockModel
	for rows.Next() {
		var stock StockModel

		rows.Scan(
			&stock.ID,
			&stock.WarehouseID,
			&stock.ProductID,
			&stock.Stock,
			&stock.IsDelete,
			&stock.CreatedBy,
			&stock.CreatedAt,
			&stock.UpdatedBy,
			&stock.UpdatedAt,
		)

		stocks = append(stocks, stock)
	}

	return stocks, nil

}

func GetAllStockByProductID(ctx context.Context, db *sql.DB, productID uuid.UUID) (
	[]StockModel, error) {

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
		FROM 
			stock
		WHERE 
			product_id=$1
	`)

	rows, err := db.QueryContext(ctx, query, productID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var stocks []StockModel
	for rows.Next() {
		var stock StockModel

		rows.Scan(
			&stock.ID,
			&stock.WarehouseID,
			&stock.ProductID,
			&stock.Stock,
			&stock.IsDelete,
			&stock.CreatedBy,
			&stock.CreatedAt,
			&stock.UpdatedBy,
			&stock.UpdatedAt,
		)

		stocks = append(stocks, stock)
	}

	return stocks, nil

}

func (s *StockModel) Insert(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		INSERT INTO stock(
			warehouse_id,
			product_id,
			stock,
			created_by,
			created_at)
		VALUES(
			$1,$2,$3,$4,now())
		RETURNING 
			id, created_at
	`)

	err := db.QueryRowContext(ctx, query,
		s.WarehouseID, s.ProductID, s.Stock, s.CreatedBy).Scan(
		&s.ID, &s.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil

}
func (s *StockModel) Update(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		UPDATE stock
		SET
			stock = $1,
			updated_at=NOW(),
			updated_by=$2
		WHERE 
			id=$3
		RETURNING 
			id,created_at,updated_at,created_by
	`)

	err := db.QueryRowContext(ctx, query,
		s.Stock, s.UpdatedBy, s.ID).Scan(
		&s.ID, &s.CreatedAt, &s.UpdatedAt, &s.CreatedBy,
	)

	if err != nil {
		return err
	}

	return nil

}

func (s *StockModel) Delete(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		UPDATE stock
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
