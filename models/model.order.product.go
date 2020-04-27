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
	OrderProductModel struct {
		ID        uuid.UUID
		OrderID   uuid.UUID
		ProductID uuid.UUID
		Quantity  uint
		SubTotal  decimal.Decimal
		IsDelete  bool
		CreatedBy uuid.UUID
		CreatedAt time.Time
		UpdatedBy uuid.NullUUID
		UpdatedAt pq.NullTime
	}

	OrderProductResponse struct {
		ID        uuid.UUID       `json:"id"`
		Order     OrderResponse   `json:"order"`
		Product   ProductResponse `json:"product"`
		Quantity  uint            `json:"quantity"`
		SubTotal  decimal.Decimal `json:"sub_total"`
		IsDelete  bool            `json:"is_delete"`
		CreatedBy uuid.UUID       `json:"created_by"`
		CreatedAt time.Time       `json:"created_at"`
		UpdatedBy uuid.UUID       `json:"updated_by"`
		UpdatedAt time.Time       `json:"updated_at"`
	}
)

func (s OrderProductModel) Response(ctx context.Context, db *sql.DB, logger *helpers.Logger) (
	OrderProductResponse, error) {

	order, err := GetOneOrder(ctx, db, s.OrderID)
	if err != nil {
		logger.Err.Printf(`model.order.product.go/GetOneOrder/%v`, err)
		return OrderProductResponse{}, err
	}

	orderResponse, err := order.Response(ctx, db, logger)
	if err != nil {
		logger.Err.Printf(`model.order.product.go/OrderResponse/%v`, err)
		return OrderProductResponse{}, err
	}

	product, err := GetOneProduct(ctx, db, s.ProductID)
	if err != nil {
		logger.Err.Printf(`model.order.product.go/GetOneProduct/%v`, err)
		return OrderProductResponse{}, err
	}

	productResponse, err := product.Response(ctx, db, logger)
	if err != nil {
		logger.Err.Printf(`model.order.product.go/ProductResponse/%v`, err)
		return OrderProductResponse{}, err
	}

	return OrderProductResponse{
		Order:     orderResponse,
		Product:   productResponse,
		Quantity:  s.Quantity,
		SubTotal:  s.SubTotal,
		IsDelete:  s.IsDelete,
		CreatedBy: s.CreatedBy,
		CreatedAt: s.CreatedAt,
		UpdatedBy: s.UpdatedBy.UUID,
		UpdatedAt: s.UpdatedAt.Time,
	}, nil

}

func GetOneOrderProduct(ctx context.Context, db *sql.DB, orderProductID uuid.UUID) (OrderProductModel, error) {

	query := fmt.Sprintf(`
		SELECT
				id,
				order_id,
				product_id,
				quantity,
				subtotal,
				is_delete,
				created_by,
				created_at,
				updated_by,
				updated_at
		FROM order_product 
		WHERE
			id = $1
	`)

	var orderProduct OrderProductModel

	err := db.QueryRowContext(ctx, query, orderProductID).Scan(
		&orderProduct.ID,
		&orderProduct.OrderID,
		&orderProduct.ProductID,
		&orderProduct.Quantity,
		&orderProduct.SubTotal,
		&orderProduct.IsDelete,
		&orderProduct.CreatedBy,
		&orderProduct.CreatedAt,
		&orderProduct.UpdatedBy,
		&orderProduct.UpdatedAt,
	)

	if err != nil {
		return OrderProductModel{}, err
	}

	return orderProduct, nil

}

func GetAllOrderProduct(ctx context.Context, db *sql.DB, filter helpers.Filter) ([]OrderProductModel, error) {

	query := fmt.Sprintf(`
		SELECT
				id,
				order_id,
				product_id,
				quantity,
				subtotal,
				is_delete,
				created_by,
				created_at,
				updated_by,
				updated_at
		FROM order_product 
		ORDER BY order_id  %s
		LIMIT $1 OFFSET $2`, filter.Dir)

	rows, err := db.QueryContext(ctx, query, filter.Limit, filter.Offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var orderProducts []OrderProductModel
	for rows.Next() {
		var orderProduct OrderProductModel

		rows.Scan(
			&orderProduct.ID,
			&orderProduct.OrderID,
			&orderProduct.ProductID,
			&orderProduct.Quantity,
			&orderProduct.SubTotal,
			&orderProduct.IsDelete,
			&orderProduct.CreatedBy,
			&orderProduct.CreatedAt,
			&orderProduct.UpdatedBy,
			&orderProduct.UpdatedAt,
		)

		orderProducts = append(orderProducts, orderProduct)
	}

	return orderProducts, nil

}

func (s *OrderProductModel) Insert(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		INSERT INTO order_product(
			order_id,
			product_id,
			quantity,
			subtotal,
			created_by,
			created_at)
		VALUES(
		$1,$2,$3,$4,$5,now())
		RETURNING id, created_at,is_delete`)

	err := db.QueryRowContext(ctx, query,
		s.OrderID, s.ProductID, s.Quantity, s.SubTotal, s.CreatedBy).Scan(
		&s.ID, &s.CreatedAt, &s.IsDelete,
	)

	if err != nil {
		return err
	}

	return nil

}

func (s *OrderProductModel) Delete(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		UPDATE order_product
		SET
			is_delete = true,
			updated_by = $1,
			updated_at = NOW()
		WHERE id=$2`)

	_, err := db.ExecContext(ctx, query,
		s.UpdatedBy, s.ID)

	if err != nil {
		return err
	}

	return nil
}
