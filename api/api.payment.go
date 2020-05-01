package api

import (
	"afiqo-location/helpers"
	"afiqo-location/models"
	"context"
	"database/sql"
	"github.com/gomodule/redigo/redis"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

type (
	PaymentModule struct {
		db     *sql.DB
		cache  *redis.Pool
		logger *helpers.Logger
		name   string
	}

	PaymentUpdateParam struct {
		ID uuid.UUID `json:"id"`
	}
)

func NewPaymentModule(db *sql.DB, cache *redis.Pool, logger *helpers.Logger) *PaymentModule {
	return &PaymentModule{
		db:     db,
		cache:  cache,
		logger: logger,
		name:   "module/payment",
	}
}

func (s PaymentModule) Update(ctx context.Context, param PaymentUpdateParam) (interface{}, *helpers.Error) {

	payment, err := models.GetOnePayment(ctx, s.db, param.ID)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Update/GetOnePayment", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	order, err := models.GetOneOrder(ctx, s.db, payment.OrderID)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Update/GetOneOrder", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	if order.CustomerID != uuid.FromStringOrNil(ctx.Value("user_id").(string)) {
		return nil, helpers.ErrorWrap(err, s.name, "Update/GetOneOrder", helpers.OrderErrorMessage,
			http.StatusInternalServerError)
	}

	payment = models.PaymentModel{
		ID:     param.ID,
		Status: 1,
		UpdatedBy: uuid.NullUUID{
			UUID:  uuid.FromStringOrNil(ctx.Value("user_id").(string)),
			Valid: true,
		},
	}

	err = payment.Update(ctx, s.db)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Update/Update", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	order = models.OrderModel{
		ID:     payment.OrderID,
		Status: 1,
		UpdatedBy: uuid.NullUUID{
			UUID:  uuid.FromStringOrNil(ctx.Value("user_id").(string)),
			Valid: true,
		},
	}

	err = order.UpdateStatus(ctx, s.db)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Update/UpdateStatus", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	response, err := payment.Response(ctx, s.db, s.logger)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Update/Response", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return response, nil

}
