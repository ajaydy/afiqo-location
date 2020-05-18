package models

import (
	"afiqo-location/helpers"
	"context"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"strings"
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

func GetOneProduct(ctx context.Context, db *sql.DB, productID uuid.UUID) (ProductModel, error) {

	query := fmt.Sprintf(`
		SELECT
			id,
			supplier_id,
			category_id,
			name,
			stock,
			price,
			description,
			is_delete,
			created_by,
			created_at,
			updated_by,
			updated_at
		FROM 
			product
		WHERE 
			id = $1
	`)

	var product ProductModel
	err := db.QueryRowContext(ctx, query, productID).Scan(
		&product.ID,
		&product.SupplierID,
		&product.CategoryID,
		&product.Name,
		&product.Stock,
		&product.Price,
		&product.Description,
		&product.IsDelete,
		&product.CreatedBy,
		&product.CreatedAt,
		&product.UpdatedBy,
		&product.UpdatedAt,
	)

	if err != nil {
		return ProductModel{}, err
	}

	return product, nil

}

func GetAllProduct(ctx context.Context, db *sql.DB, filter helpers.Filter) ([]ProductModel, error) {

	var filters []string

	if filter.CategoryID != uuid.Nil {
		filters = append(filters, fmt.Sprintf(`
			category_id = '%s'`,
			filter.CategoryID))
	}

	if filter.SupplierID != uuid.Nil {
		filters = append(filters, fmt.Sprintf(`
			supplier_id = '%s'`,
			filter.SupplierID))
	}

	filterJoin := strings.Join(filters, " AND ")
	if filterJoin != "" {
		filterJoin = fmt.Sprintf("AND %s", filterJoin)
	}

	query := fmt.Sprintf(`
		SELECT
			id,
			supplier_id,
			category_id,
			name,
			stock,
			price,
			description,
			is_delete,
			created_by,
			created_at,
			updated_by,
			updated_at
		FROM 
			product
		WHERE 
			is_delete = false
		%s
		ORDER BY 
			name  %s   
		LIMIT $1 OFFSET $2`,
		filterJoin, filter.Dir)

	rows, err := db.QueryContext(ctx, query, filter.Limit, filter.Offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var products []ProductModel
	for rows.Next() {
		var product ProductModel

		rows.Scan(
			&product.ID,
			&product.SupplierID,
			&product.CategoryID,
			&product.Name,
			&product.Stock,
			&product.Price,
			&product.Description,
			&product.IsDelete,
			&product.CreatedBy,
			&product.CreatedAt,
			&product.UpdatedBy,
			&product.UpdatedAt,
		)

		products = append(products, product)
	}

	return products, nil

}

func GetAllProductForCustomer(ctx context.Context, db *sql.DB, filter helpers.Filter, warehouseID uuid.UUID) (
	[]ProductModel, error) {

	var filters []string

	if filter.CategoryID != uuid.Nil {
		filters = append(filters, fmt.Sprintf(`
			category_id = '%s'`,
			filter.CategoryID))
	}

	if filter.SupplierID != uuid.Nil {
		filters = append(filters, fmt.Sprintf(`
			supplier_id = '%s'`,
			filter.SupplierID))
	}

	filterJoin := strings.Join(filters, " AND ")
	if filterJoin != "" {
		filterJoin = fmt.Sprintf("AND %s", filterJoin)
	}

	query := fmt.Sprintf(`
		SELECT
			p.id,
			p.supplier_id,
			p.category_id,
			p.name,
			s.stock,
			p.price,
			p.description,
			p.is_delete,
			p.created_by,
			p.created_at,
			p.updated_by,
			p.updated_at
		FROM 
			product p
		INNER JOIN
			stock s
		ON
			p.id = s.product_id
		WHERE 
			p.is_delete = false
		AND 
			s.warehouse_id = $1
		%s
		ORDER BY 
			name  %s   
		LIMIT $2 OFFSET $3`,
		filterJoin, filter.Dir)

	rows, err := db.QueryContext(ctx, query, warehouseID, filter.Limit, filter.Offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var products []ProductModel
	for rows.Next() {
		var product ProductModel

		rows.Scan(
			&product.ID,
			&product.SupplierID,
			&product.CategoryID,
			&product.Name,
			&product.Stock,
			&product.Price,
			&product.Description,
			&product.IsDelete,
			&product.CreatedBy,
			&product.CreatedAt,
			&product.UpdatedBy,
			&product.UpdatedAt,
		)

		products = append(products, product)
	}

	return products, nil

}

func GetAllProductBySupplierID(ctx context.Context, db *sql.DB, filter helpers.Filter, supplierID uuid.UUID) (
	[]ProductModel, error) {

	query := fmt.Sprintf(`
		SELECT
			id,
			supplier_id,
			category_id,
			name,
			stock,
			price,
			description,
			is_delete,
			created_by,
			created_at,
			updated_by,
			updated_at
		FROM
			product
		WHERE 
			is_delete = false
		AND 
			supplier_id = $1
		ORDER BY
			updated_at %s ,created_at %s 
		LIMIT $2 OFFSET $3`,
		filter.Dir, filter.Dir)

	rows, err := db.QueryContext(ctx, query, supplierID, filter.Limit, filter.Offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var products []ProductModel
	for rows.Next() {
		var product ProductModel

		rows.Scan(
			&product.ID,
			&product.SupplierID,
			&product.CategoryID,
			&product.Name,
			&product.Stock,
			&product.Price,
			&product.Description,
			&product.IsDelete,
			&product.CreatedBy,
			&product.CreatedAt,
			&product.UpdatedBy,
			&product.UpdatedAt,
		)

		products = append(products, product)
	}

	return products, nil

}

func (s *ProductModel) Insert(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		INSERT INTO product(
			supplier_id,
			category_id,
			name,
			stock,
			price,
			description,
			created_by,
			created_at)
		VALUES(
			$1,$2,$3,$4,$5,$6,$7,now())
		RETURNING 
			id, created_at,is_delete
	`)

	err := db.QueryRowContext(ctx, query,
		s.SupplierID, s.CategoryID, s.Name, s.Stock, s.Price, s.Description, s.CreatedBy).Scan(
		&s.ID, &s.CreatedAt, &s.IsDelete,
	)

	if err != nil {
		return err
	}

	return nil

}

func (s *ProductModel) Update(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		UPDATE product
		SET
			name=$1,
			stock=$2,
			price=$3,
			description=$4,
			updated_at=NOW(),
			updated_by=$5
		WHERE 
			id=$6
		RETURNING 
			id,created_at,updated_at,created_by,is_delete
	`)

	err := db.QueryRowContext(ctx, query,
		s.Name, s.Stock, s.Price, s.Description, s.UpdatedBy, s.ID).Scan(
		&s.ID, &s.CreatedAt, &s.UpdatedAt, &s.CreatedBy, &s.IsDelete,
	)

	if err != nil {
		return err
	}

	return nil

}

func (s *ProductModel) StockUpdate(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		UPDATE product
		SET
			stock=$1,
			updated_at=NOW(),
			updated_by=$2
		WHERE 
			id=$3
		RETURNING 
			id,created_at,updated_at,created_by,is_delete
	`)

	err := db.QueryRowContext(ctx, query,
		s.Stock, s.UpdatedBy, s.ID).Scan(
		&s.ID, &s.CreatedAt, &s.UpdatedAt, &s.CreatedBy, &s.IsDelete,
	)

	if err != nil {
		return err
	}

	return nil

}

func (s *ProductModel) Delete(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		UPDATE product
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
