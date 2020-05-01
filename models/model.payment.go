package models

import (
	"afiqo-location/helpers"
	"context"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	"time"
)

type (
	PaymentModel struct {
		ID        uuid.UUID
		OrderID   uuid.UUID
		Status    int
		IsDelete  bool
		CreatedBy uuid.UUID
		CreatedAt time.Time
		UpdatedBy uuid.NullUUID
		UpdatedAt pq.NullTime
	}

	PaymentResponse struct {
		ID        uuid.UUID     `json:"id"`
		Order     OrderResponse `json:"order"`
		Status    int           `json:"status"`
		IsDelete  bool          `json:"is_delete"`
		CreatedBy uuid.UUID     `json:"created_by"`
		CreatedAt time.Time     `json:"created_at"`
		UpdatedBy uuid.UUID     `json:"updated_by"`
		UpdatedAt time.Time     `json:"updated_at"`
	}
)

func (s PaymentModel) Response(ctx context.Context, db *sql.DB, logger *helpers.Logger) (PaymentResponse, error) {

	order, err := GetOneOrder(ctx, db, s.OrderID)
	if err != nil {
		logger.Err.Printf(`model.payment.go/GetOneOrder/%v`, err)
		return PaymentResponse{}, err
	}

	orderResponse, err := order.Response(ctx, db, logger)
	if err != nil {
		logger.Err.Printf(`model.payment.go/OrderResponse/%v`, err)
		return PaymentResponse{}, err
	}

	return PaymentResponse{
		Order:     orderResponse,
		Status:    s.Status,
		IsDelete:  s.IsDelete,
		CreatedBy: s.CreatedBy,
		CreatedAt: s.CreatedAt,
		UpdatedBy: s.UpdatedBy.UUID,
		UpdatedAt: s.UpdatedAt.Time,
	}, nil

}

func GetOnePayment(ctx context.Context, db *sql.DB, paymentID uuid.UUID) (PaymentModel, error) {

	query := fmt.Sprintf(`
		SELECT
				id,
				order_id,
				status,
				is_delete,
				created_by,
				created_at,
				updated_by,
				updated_at
		FROM payment 
		WHERE
			id = $1
	`)

	var payment PaymentModel

	err := db.QueryRowContext(ctx, query, paymentID).Scan(
		&payment.ID,
		&payment.OrderID,
		&payment.Status,
		&payment.IsDelete,
		&payment.CreatedBy,
		&payment.CreatedAt,
		&payment.UpdatedBy,
		&payment.UpdatedAt,
	)

	if err != nil {
		return PaymentModel{}, err
	}

	return payment, nil

}

func GetAllPayment(ctx context.Context, db *sql.DB, filter helpers.Filter) ([]PaymentModel, error) {

	query := fmt.Sprintf(`
		SELECT
			id,
			order_id,
			status,
			is_delete,
			created_by,
			created_at,
			updated_by,
			updated_at
		FROM payment
		ORDER BY order_id  %s
		LIMIT $1 OFFSET $2`, filter.Dir)

	rows, err := db.QueryContext(ctx, query, filter.Limit, filter.Offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var payments []PaymentModel
	for rows.Next() {
		var payment PaymentModel

		rows.Scan(
			&payment.ID,
			&payment.OrderID,
			&payment.Status,
			&payment.IsDelete,
			&payment.CreatedBy,
			&payment.CreatedAt,
			&payment.UpdatedBy,
			&payment.UpdatedAt,
		)

		payments = append(payments, payment)
	}

	return payments, nil

}

func (s *PaymentModel) Insert(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		INSERT INTO payment(
			order_id,
			status,
			created_by,
			created_at)
		VALUES(
		$1,$2,$3,now())
		RETURNING id, created_at,is_delete`)

	err := db.QueryRowContext(ctx, query,
		s.OrderID, s.Status, s.CreatedBy).Scan(
		&s.ID, &s.CreatedAt, &s.IsDelete,
	)

	if err != nil {
		return err
	}

	return nil

}

func (s *PaymentModel) Update(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		UPDATE payment
		SET
			status=$1
			updated_at=NOW(),
			updated_by=$2
		WHERE id=$3
		RETURNING id,created_at,updated_at,created_by`)

	err := db.QueryRowContext(ctx, query,
		s.Status, s.UpdatedBy, s.ID).Scan(
		&s.ID, &s.CreatedAt, &s.UpdatedAt, &s.CreatedBy,
	)

	if err != nil {
		return err
	}

	return nil

}

func (s *PaymentModel) Delete(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		UPDATE payment
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
