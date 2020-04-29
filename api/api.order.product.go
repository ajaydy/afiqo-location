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
	OrderProductModule struct {
		db     *sql.DB
		cache  *redis.Pool
		logger *helpers.Logger
		name   string
	}

	OrderProductDetailParam struct {
		ID uuid.UUID `json:"id"`
	}

	OrderProductAddParam struct {
		OrderID   uuid.UUID `json:"order_id"`
		ProductID uuid.UUID `json:"product_id"`
		Quantity  uint      `json:"quantity"`
	}

	OrderProductDeleteParam struct {
		ID uuid.UUID `json:"id"`
	}
)

func NewOrderProductModule(db *sql.DB, cache *redis.Pool, logger *helpers.Logger) *OrderProductModule {
	return &OrderProductModule{
		db:     db,
		cache:  cache,
		logger: logger,
		name:   "module/orderProduct",
	}
}

func (s OrderProductModule) Detail(ctx context.Context, param OrderProductDetailParam) (interface{}, *helpers.Error) {
	orderProduct, err := models.GetOneOrderProduct(ctx, s.db, param.ID)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Detail/GetOneOrderProduct", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	response, err := orderProduct.Response(ctx, s.db, s.logger)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Detail/Response", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return response, nil
}

func (s OrderProductModule) List(ctx context.Context, filter helpers.Filter) (interface{}, *helpers.Error) {

	orderProducts, err := models.GetAllOrderProduct(ctx, s.db, filter)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "List/GetAllOrderProduct", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	var orderProductResponse []models.OrderProductResponse
	for _, orderProduct := range orderProducts {
		response, err := orderProduct.Response(ctx, s.db, s.logger)
		if err != nil {
			return nil, helpers.ErrorWrap(err, s.name, "List/Response", helpers.InternalServerError,
				http.StatusInternalServerError)
		}

		orderProductResponse = append(orderProductResponse, response)
	}

	return orderProductResponse, nil
}

func (s OrderProductModule) Delete(ctx context.Context, param OrderProductDeleteParam) (interface{}, *helpers.Error) {

	orderProduct := models.OrderProductModel{
		ID: param.ID,
		UpdatedBy: uuid.NullUUID{
			UUID:  uuid.FromStringOrNil(ctx.Value("user_id").(string)),
			Valid: true,
		},
	}

	err := orderProduct.Delete(ctx, s.db)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Delete/Delete", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return nil, nil

}
