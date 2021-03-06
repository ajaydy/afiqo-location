package models

import (
	"afiqo-location/helpers"
	"afiqo-location/util"
	"context"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	"strings"
	"time"
)

type (
	ShipmentModel struct {
		ID        uuid.UUID
		CourierID uuid.UUID
		OrderID   uuid.UUID
		Status    int
		IsDelete  bool
		CreatedBy uuid.UUID
		CreatedAt time.Time
		UpdatedBy uuid.NullUUID
		UpdatedAt pq.NullTime
	}

	ShipmentResponse struct {
		ID        uuid.UUID       `json:"id"`
		Courier   CourierResponse `json:"courier"`
		Order     OrderResponse   `json:"order"`
		Status    string          `json:"status"`
		IsDelete  bool            `json:"is_delete"`
		CreatedBy uuid.UUID       `json:"created_by"`
		CreatedAt time.Time       `json:"created_at"`
		UpdatedBy uuid.UUID       `json:"updated_by"`
		UpdatedAt time.Time       `json:"updated_at"`
	}
)

func (s ShipmentModel) Response(ctx context.Context, db *sql.DB, logger *helpers.Logger) (ShipmentResponse, error) {

	courier, err := GetOneCourier(ctx, db, s.CourierID)
	if err != nil {
		logger.Err.Printf(`model.shipment.go/GetOneCourier/%v`, err)
		return ShipmentResponse{}, nil
	}

	order, err := GetOneOrder(ctx, db, s.OrderID)
	if err != nil {
		logger.Err.Printf(`model.shipment.go/GetOneOrder/%v`, err)
		return ShipmentResponse{}, nil
	}

	orderResponse, err := order.Response(ctx, db, logger)
	if err != nil {
		logger.Err.Printf(`model.shipment.go/orderResponse/%v`, err)
		return ShipmentResponse{}, nil
	}

	status := util.GetShipmentStatus(s.Status)

	return ShipmentResponse{
		ID:        s.ID,
		Courier:   courier.Response(),
		Order:     orderResponse,
		Status:    status,
		IsDelete:  s.IsDelete,
		CreatedBy: s.CreatedBy,
		CreatedAt: s.CreatedAt,
		UpdatedBy: s.UpdatedBy.UUID,
		UpdatedAt: s.UpdatedAt.Time,
	}, nil
}

func GetOneShipment(ctx context.Context, db *sql.DB, shipmentID uuid.UUID) (ShipmentModel, error) {

	query := fmt.Sprintf(`
		SELECT
			id,
			courier_id,
			order_id,
			status,
			is_delete,
			created_by,
			created_at,
			updated_by,
			updated_at
		FROM 
			shipment
		WHERE 
			id = $1
	`)

	var shipment ShipmentModel
	err := db.QueryRowContext(ctx, query, shipmentID).Scan(
		&shipment.ID,
		&shipment.CourierID,
		&shipment.OrderID,
		&shipment.Status,
		&shipment.IsDelete,
		&shipment.CreatedBy,
		&shipment.CreatedAt,
		&shipment.UpdatedBy,
		&shipment.UpdatedAt,
	)

	if err != nil {
		return ShipmentModel{}, err
	}

	return shipment, nil

}

func GetAllShipment(ctx context.Context, db *sql.DB, filter helpers.Filter) ([]ShipmentModel, error) {

	var filters []string

	if filter.CourierID != uuid.Nil {
		filters = append(filters, fmt.Sprintf(`
			courier_id = '%s'`,
			filter.CourierID))
	}

	if filter.OrderID != uuid.Nil {
		filters = append(filters, fmt.Sprintf(`
			order_id = '%s'`,
			filter.OrderID))
	}

	filterJoin := strings.Join(filters, " AND ")
	if filterJoin != "" {
		filterJoin = fmt.Sprintf("AND %s", filterJoin)
	}

	query := fmt.Sprintf(`
		SELECT
			id,
			courier_id,
			order_id,
			status,
			is_delete,
			created_by,
			created_at,
			updated_by,
			updated_at
		FROM 
			shipment
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

	var shipments []ShipmentModel
	for rows.Next() {
		var shipment ShipmentModel

		rows.Scan(
			&shipment.ID,
			&shipment.CourierID,
			&shipment.OrderID,
			&shipment.Status,
			&shipment.IsDelete,
			&shipment.CreatedBy,
			&shipment.CreatedAt,
			&shipment.UpdatedBy,
			&shipment.UpdatedAt,
		)

		shipments = append(shipments, shipment)
	}

	return shipments, nil

}

func GetAllShipmentByCustomerID(ctx context.Context, db *sql.DB, filter helpers.Filter, customerID uuid.UUID) (
	[]ShipmentModel, error) {

	query := fmt.Sprintf(`
		SELECT
			s.id,
			s.courier_id,
			s.order_id,
			s.status,
			s.is_delete,
			s.created_by,
			s.created_at,
			s.updated_by,
			s.updated_at
		FROM 
			shipment s
		INNER JOIN
			"order" o
		ON 
			o.id = s.order_id
		WHERE 
			o.customer_id = $1
		LIMIT $2 OFFSET $3
	`)

	rows, err := db.QueryContext(ctx, query, customerID, filter.Limit, filter.Offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var shipments []ShipmentModel
	for rows.Next() {
		var shipment ShipmentModel

		rows.Scan(
			&shipment.ID,
			&shipment.CourierID,
			&shipment.OrderID,
			&shipment.Status,
			&shipment.IsDelete,
			&shipment.CreatedBy,
			&shipment.CreatedAt,
			&shipment.UpdatedBy,
			&shipment.UpdatedAt,
		)

		shipments = append(shipments, shipment)
	}

	return shipments, nil

}

func GetAllShipmentByCourierID(ctx context.Context, db *sql.DB, filter helpers.Filter, courierID uuid.UUID) (
	[]ShipmentModel, error) {

	query := fmt.Sprintf(`
		SELECT
			id,
			courier_id,
			order_id,
			status,
			is_delete,
			created_by,
			created_at,
			updated_by,
			updated_at
		FROM 
			shipment
		WHERE 
			is_delete = false
		AND 
			courier_id = $1
		ORDER BY
			updated_at %s ,created_at %s
		LIMIT $2 OFFSET $3`,
		filter.Dir, filter.Dir)

	rows, err := db.QueryContext(ctx, query, courierID, filter.Limit, filter.Offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var shipments []ShipmentModel
	for rows.Next() {
		var shipment ShipmentModel

		rows.Scan(
			&shipment.ID,
			&shipment.CourierID,
			&shipment.OrderID,
			&shipment.Status,
			&shipment.IsDelete,
			&shipment.CreatedBy,
			&shipment.CreatedAt,
			&shipment.UpdatedBy,
			&shipment.UpdatedAt,
		)

		shipments = append(shipments, shipment)
	}

	return shipments, nil

}

func (s *ShipmentModel) Insert(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		INSERT INTO shipment(
			courier_id,
			order_id,
			status,
			created_by,
			created_at)
		VALUES(
			$1,$2,$3,$4,now())
		RETURNING 
			id, created_at,is_delete
	`)

	err := db.QueryRowContext(ctx, query,
		s.CourierID, s.OrderID, s.Status, s.CreatedBy).Scan(
		&s.ID, &s.CreatedAt, &s.IsDelete,
	)

	if err != nil {
		return err
	}

	return nil

}

func (s *ShipmentModel) UpdateStatus(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		UPDATE shipment
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
func (s *ShipmentModel) Delete(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		UPDATE shipment
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
