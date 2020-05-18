package models

import (
	"afiqo-location/helpers"
	"afiqo-location/util"
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
	OrderModel struct {
		ID               uuid.UUID
		CustomerID       uuid.UUID
		WarehouseID      uuid.UUID
		DeliveryDatetime time.Time
		DeliveryAddress  string
		Latitude         decimal.Decimal
		Longitude        decimal.Decimal
		Status           int
		TotalPrice       decimal.Decimal
		IsDelete         bool
		CreatedBy        uuid.UUID
		CreatedAt        time.Time
		UpdatedBy        uuid.NullUUID
		UpdatedAt        pq.NullTime
	}

	OrderResponse struct {
		ID               uuid.UUID         `json:"id"`
		Customer         CustomerResponse  `json:"customer"`
		Warehouse        WarehouseResponse `json:"warehouse"`
		DeliveryDatetime time.Time         `json:"delivery_datetime"`
		DeliveryAddress  string            `json:"delivery_address"`
		Latitude         decimal.Decimal   `json:"latitude"`
		Longitude        decimal.Decimal   `json:"longitude"`
		Status           string            `json:"status"`
		TotalPrice       decimal.Decimal   `json:"total_price"`
		IsDelete         bool              `json:"is_delete"`
		CreatedBy        uuid.UUID         `json:"created_by"`
		CreatedAt        time.Time         `json:"created_at"`
		UpdatedBy        uuid.UUID         `json:"updated_by"`
		UpdatedAt        time.Time         `json:"updated_at"`
	}
)

func (s OrderModel) Response(ctx context.Context, db *sql.DB, logger *helpers.Logger) (OrderResponse, error) {

	customer, err := GetOneCustomer(ctx, db, s.CustomerID)
	if err != nil {
		logger.Err.Printf(`model.order.go/GetOneCustomer/%v`, err)
		return OrderResponse{}, nil
	}

	warehouse, err := GetOneWarehouse(ctx, db, s.WarehouseID)
	if err != nil {
		logger.Err.Printf(`model.order.go/GetOneWarehouse/%v`, err)
		return OrderResponse{}, nil
	}

	status := util.GetOrderStatus(s.Status)

	return OrderResponse{
		ID:               s.ID,
		Customer:         customer.Response(),
		Warehouse:        warehouse.Response(),
		DeliveryDatetime: s.DeliveryDatetime,
		DeliveryAddress:  s.DeliveryAddress,
		Longitude:        s.Longitude,
		Latitude:         s.Latitude,
		Status:           status,
		TotalPrice:       s.TotalPrice,
		IsDelete:         s.IsDelete,
		CreatedBy:        s.CreatedBy,
		CreatedAt:        s.CreatedAt,
		UpdatedBy:        s.UpdatedBy.UUID,
		UpdatedAt:        s.UpdatedAt.Time,
	}, nil

}

func GetOneOrder(ctx context.Context, db *sql.DB, orderID uuid.UUID) (OrderModel, error) {

	query := fmt.Sprintf(`
		SELECT
			id,
			warehouse_id,
			customer_id,
			delivery_datetime,
			delivery_address,
			latitude,
			longitude,
			status,
			total_price,
			is_delete,
			created_by,
			created_at,
			updated_by,
			updated_at
		FROM 
			"order"
		WHERE 
			id = $1
	`)

	var order OrderModel
	err := db.QueryRowContext(ctx, query, orderID).Scan(
		&order.ID,
		&order.WarehouseID,
		&order.CustomerID,
		&order.DeliveryDatetime,
		&order.DeliveryAddress,
		&order.Latitude,
		&order.Longitude,
		&order.Status,
		&order.TotalPrice,
		&order.IsDelete,
		&order.CreatedBy,
		&order.CreatedAt,
		&order.UpdatedBy,
		&order.UpdatedAt,
	)

	if err != nil {
		return OrderModel{}, err
	}

	return order, nil

}

func GetAllOrder(ctx context.Context, db *sql.DB, filter helpers.Filter) ([]OrderModel, error) {

	var filters []string

	if filter.CustomerID != uuid.Nil {
		filters = append(filters, fmt.Sprintf(`
			customer_id = '%s'`,
			filter.CustomerID))
	}

	if filter.WarehouseID != uuid.Nil {
		filters = append(filters, fmt.Sprintf(`
			warehouse_id = '%s'`,
			filter.WarehouseID))
	}

	filterJoin := strings.Join(filters, " AND ")
	if filterJoin != "" {
		filterJoin = fmt.Sprintf("AND %s", filterJoin)
	}

	query := fmt.Sprintf(`
		SELECT
			id,
			warehouse_id,
			customer_id,
			delivery_datetime,
			delivery_address,
			latitude,
			longitude,
			status,
			total_price,
			is_delete,
			created_by,
			created_at,
			updated_by,
			updated_at
		FROM 
			"order"
		WHERE 
			is_delete = false
		%s
		LIMIT $1 OFFSET $2`,
		filterJoin)

	rows, err := db.QueryContext(ctx, query, filter.Limit, filter.Offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var orders []OrderModel
	for rows.Next() {
		var order OrderModel

		rows.Scan(
			&order.ID,
			&order.WarehouseID,
			&order.CustomerID,
			&order.DeliveryDatetime,
			&order.DeliveryAddress,
			&order.Latitude,
			&order.Longitude,
			&order.Status,
			&order.TotalPrice,
			&order.IsDelete,
			&order.CreatedBy,
			&order.CreatedAt,
			&order.UpdatedBy,
			&order.UpdatedAt,
		)

		orders = append(orders, order)
	}

	return orders, nil

}

func GetAllOrderByCustomerID(ctx context.Context, db *sql.DB, filter helpers.Filter, customerID uuid.UUID) (
	[]OrderModel, error) {

	query := fmt.Sprintf(`
		SELECT
			id,
			warehouse_id,
			customer_id,
			delivery_datetime,
			delivery_address,
			latitude,
			longitude,
			status,
			total_price,
			is_delete,
			created_by,
			created_at,
			updated_by,
			updated_at
		FROM 
			"order"
		WHERE 
			is_delete = false
		AND 
			customer_id = $1
		ORDER BY
			updated_at %s ,created_at %s 
		LIMIT $2 OFFSET $3`, filter.Dir, filter.Dir)

	rows, err := db.QueryContext(ctx, query, customerID, filter.Limit, filter.Offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var orders []OrderModel
	for rows.Next() {
		var order OrderModel

		rows.Scan(
			&order.ID,
			&order.WarehouseID,
			&order.CustomerID,
			&order.DeliveryDatetime,
			&order.DeliveryAddress,
			&order.Latitude,
			&order.Longitude,
			&order.Status,
			&order.TotalPrice,
			&order.IsDelete,
			&order.CreatedBy,
			&order.CreatedAt,
			&order.UpdatedBy,
			&order.UpdatedAt,
		)

		orders = append(orders, order)
	}

	return orders, nil

}

func (s *OrderModel) Insert(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		INSERT INTO "order"(
			warehouse_id,
			customer_id,
			delivery_datetime,
			delivery_address,
			latitude,
			longitude,
			status,
			total_price,
			created_by,
			created_at
		)VALUES(
			$1,$2,$3,$4,$5,$6,$7,$8,$9,now())
		RETURNING 
			id, created_at,is_delete
	`)

	err := db.QueryRowContext(ctx, query,
		s.WarehouseID, s.CustomerID, s.DeliveryDatetime, s.DeliveryAddress, s.Latitude, s.Longitude, s.Status, s.TotalPrice,
		s.CreatedBy).Scan(
		&s.ID, &s.CreatedAt, &s.IsDelete,
	)

	if err != nil {
		return err
	}

	return nil

}

func (s *OrderModel) UpdateStatus(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		UPDATE "order"
		SET
			status=$1,
			updated_at=NOW(),
			updated_by=$2
		WHERE 
			id=$3
		RETURNING 
			id,created_at,updated_at,created_by
	`)

	err := db.QueryRowContext(ctx, query,
		s.Status, s.UpdatedBy, s.ID).Scan(
		&s.ID, &s.CreatedAt, &s.UpdatedAt, &s.CreatedBy,
	)

	if err != nil {
		return err
	}

	return nil

}

func (s *OrderModel) UpdatePrice(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		UPDATE "order"
		SET
			total_price=$1,
			updated_at=NOW(),
			updated_by=$2
		WHERE 
			id=$3
		RETURNING 
			id,created_at,updated_at,created_by
	`)

	err := db.QueryRowContext(ctx, query,
		s.TotalPrice, s.UpdatedBy, s.ID).Scan(
		&s.ID, &s.CreatedAt, &s.UpdatedAt, &s.CreatedBy,
	)

	if err != nil {
		return err
	}

	return nil

}

func (s *OrderModel) Delete(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		UPDATE order
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
