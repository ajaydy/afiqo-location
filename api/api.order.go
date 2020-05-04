package api

import (
	"afiqo-location/helpers"
	"afiqo-location/maps"
	"afiqo-location/models"
	"afiqo-location/util"
	"context"
	"database/sql"
	"github.com/gomodule/redigo/redis"
	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"net/http"
	"time"
)

type (
	Product struct {
		ID       uuid.UUID `json:"id"`
		Quantity uint      `json:"quantity"`
	}

	OrderModule struct {
		db     *sql.DB
		cache  *redis.Pool
		logger *helpers.Logger
		name   string
	}

	OrderDetailParam struct {
		ID uuid.UUID `json:"id"`
	}

	OrderParam struct {
		DeliveryAddress  string    `json:"delivery_address"`
		DeliveryDatetime time.Time `json:"delivery_datetime"`
		Product          []Product `json:"product"`
	}

	OrderDeleteParam struct {
		ID uuid.UUID `json:"id"`
	}
)

func NewOrderModule(db *sql.DB, cache *redis.Pool, logger *helpers.Logger) *OrderModule {
	return &OrderModule{
		db:     db,
		cache:  cache,
		logger: logger,
		name:   "module/order",
	}
}

func (s OrderModule) Detail(ctx context.Context, param OrderDetailParam) (interface{}, *helpers.Error) {
	order, err := models.GetOneOrder(ctx, s.db, param.ID)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Detail/GetOneOrder", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	response, err := order.Response(ctx, s.db, s.logger)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Detail/Response", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return response, nil
}

func (s OrderModule) List(ctx context.Context, filter helpers.Filter) (interface{}, *helpers.Error) {

	orders, err := models.GetAllOrder(ctx, s.db, filter)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "List/GetAllOrder", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	var orderResponse []models.OrderResponse
	for _, order := range orders {
		response, err := order.Response(ctx, s.db, s.logger)
		if err != nil {
			return nil, helpers.ErrorWrap(err, s.name, "List/Response", helpers.InternalServerError,
				http.StatusInternalServerError)
		}

		orderResponse = append(orderResponse, response)
	}

	return orderResponse, nil
}

func (s OrderModule) Order(ctx context.Context, param OrderParam) (interface{}, *helpers.Error) {

	order := models.OrderModel{
		CustomerID:       uuid.FromStringOrNil(ctx.Value("user_id").(string)),
		DeliveryDatetime: param.DeliveryDatetime,
		DeliveryAddress:  param.DeliveryAddress,
		Status:           0,
		TotalPrice:       decimal.NewFromFloat(0.00),
		CreatedBy:        uuid.FromStringOrNil(ctx.Value("user_id").(string)),
	}

	err := order.Insert(ctx, s.db)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Order/order.Insert", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	for _, orderProduct := range param.Product {

		product, err := models.GetOneProduct(ctx, s.db, orderProduct.ID)

		if err != nil {
			return nil, helpers.ErrorWrap(err, s.name, "Order/GetOneProduct", helpers.InternalServerError,
				http.StatusInternalServerError)
		}

		stocks, err := models.GetAllStock(ctx, s.db, helpers.Filter{
			FilterOption: helpers.FilterOption{
				Limit:  999,
				Offset: 0,
			},
			ProductID: product.ID,
		})

		if err != nil {
			return nil, helpers.ErrorWrap(err, s.name, "Order/GetAllStock", helpers.InternalServerError,
				http.StatusInternalServerError)
		}

		var warehouses []models.WarehouseModel
		for _, stock := range stocks {

			warehouse, err := models.GetOneWarehouse(ctx, s.db, stock.WarehouseID)

			if err != nil {
				return nil, helpers.ErrorWrap(err, s.name, "Order/GetOneWarehouse", helpers.InternalServerError,
					http.StatusInternalServerError)
			}

			warehouses = append(warehouses, warehouse)

		}

		var distanceArray []int
		for _, warehouse := range warehouses {

			distanceMatrix, err := maps.GetDistanceBetweenTwoLocations(ctx, warehouse.Address, order.DeliveryAddress)
			if err != nil {
				return nil, helpers.ErrorWrap(err, s.name, "Order/maps.GetDistanceBetweenTwoLocations",
					helpers.InternalServerError,
					http.StatusInternalServerError)
			}

			distance := distanceMatrix.Rows[0].Elements[0].Distance.Value
			distanceArray = append(distanceArray, distance)
		}

		minimum := util.GetMinDistance(distanceArray) / 1e3

		configuration, err := models.GetConfiguration(ctx, s.db)
		if err != nil {
			return nil, helpers.ErrorWrap(err, s.name, "Order/GetConfiguration", helpers.InternalServerError,
				http.StatusInternalServerError)
		}

		deliveryFee := configuration.DeliveryFee.Mul(decimal.NewFromInt(int64(minimum)))

		subtotal := decimal.Sum(product.Price.Mul(decimal.NewFromInt(int64(orderProduct.Quantity))),
			deliveryFee)

		orderProduct := models.OrderProductModel{
			OrderID:   order.ID,
			ProductID: orderProduct.ID,
			Quantity:  orderProduct.Quantity,
			SubTotal:  subtotal,
			CreatedBy: uuid.FromStringOrNil(ctx.Value("user_id").(string)),
		}

		err = orderProduct.Insert(ctx, s.db)
		if err != nil {
			return nil, helpers.ErrorWrap(err, s.name, "Order/orderProduct.Insert", helpers.InternalServerError,
				http.StatusInternalServerError)
		}

		totalStock := product.Stock - orderProduct.Quantity

		productStock := models.ProductModel{
			ID:    product.ID,
			Stock: totalStock,
			UpdatedBy: uuid.NullUUID{
				UUID:  uuid.FromStringOrNil(ctx.Value("user_id").(string)),
				Valid: true,
			},
		}

		err = productStock.StockUpdate(ctx, s.db)
		if err != nil {
			return nil, helpers.ErrorWrap(err, s.name, "Order/productStock.StockUpdate", helpers.InternalServerError,
				http.StatusInternalServerError)
		}

	}

	orderProducts, err := models.GetAllOrderProduct(ctx, s.db, helpers.Filter{
		FilterOption: helpers.FilterOption{
			Limit:  999,
			Offset: 0,
		},
		OrderID: order.ID,
	})

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Order/GetAllOrderProduct", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	var totalPrice decimal.Decimal
	for _, orderProduct := range orderProducts {
		totalPrice = decimal.Sum(totalPrice, orderProduct.SubTotal)
	}

	orderUpdate := models.OrderModel{
		ID:         order.ID,
		TotalPrice: totalPrice,
		UpdatedBy: uuid.NullUUID{
			UUID:  uuid.FromStringOrNil(ctx.Value("user_id").(string)),
			Valid: true,
		},
	}
	err = orderUpdate.UpdatePrice(ctx, s.db)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Order/UpdatePrice", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	order, err = models.GetOneOrder(ctx, s.db, order.ID)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Order/GetOneOrder", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	payment := models.PaymentModel{
		OrderID:   order.ID,
		Status:    0,
		CreatedBy: uuid.FromStringOrNil(ctx.Value("user_id").(string)),
	}

	err = payment.Insert(ctx, s.db)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Order/payment.Insert", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	response, err := order.Response(ctx, s.db, s.logger)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Order/Response", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return response, nil

}

func (s OrderModule) Delete(ctx context.Context, param OrderDeleteParam) (interface{}, *helpers.Error) {

	order := models.OrderModel{
		ID: param.ID,
		UpdatedBy: uuid.NullUUID{
			UUID:  uuid.FromStringOrNil(ctx.Value("user_id").(string)),
			Valid: true,
		},
	}

	err := order.Delete(ctx, s.db)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Delete/Delete", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return nil, nil

}
