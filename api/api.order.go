package api

import (
	"afiqo-location/helpers"
	"afiqo-location/models"
	"context"
	"database/sql"
	"fmt"
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

	Distance struct {
		WarehouseID   uuid.UUID
		DistanceValue int
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
		DeliveryAddress string          `json:"delivery_address" validate:"required"`
		Longitude       decimal.Decimal `json:"longitude" validate:"required"`
		Latitude        decimal.Decimal `json:"latitude" validate:"required"`
		Product         []Product       `json:"product" validate:"required"`
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

	now := time.Now()

	deliveryDateTime := now.AddDate(0, 0, 3)

	warehouses, err := models.GetAllWarehouseWithDistance(ctx, s.db, helpers.Filter{
		FilterOption: helpers.FilterOption{
			Limit:  1,
			Offset: 0,
		},
		Longitude: param.Longitude,
		Latitude:  param.Latitude,
	})

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Order/GetAllWarehouseWithDistance", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	var warehouseID uuid.UUID
	for _, warehouse := range warehouses {
		warehouseID = warehouse.ID
	}

	order := models.OrderModel{
		CustomerID:       uuid.FromStringOrNil(ctx.Value("user_id").(string)),
		WarehouseID:      warehouseID,
		DeliveryDatetime: deliveryDateTime,
		DeliveryAddress:  param.DeliveryAddress,
		Longitude:        param.Longitude,
		Latitude:         param.Latitude,
		Status:           0,
		CreatedBy:        uuid.FromStringOrNil(ctx.Value("user_id").(string)),
	}

	err = order.Insert(ctx, s.db)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Order/order.Insert",
			helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	for _, orderProduct := range param.Product {

		stock, err := models.GetOneStockByProductAndWarehouse(ctx, s.db, warehouseID, orderProduct.ID)
		if err != nil {
			return nil, helpers.ErrorWrap(err, s.name, "Order/GetOneStockByProductAndWarehouse",
				helpers.InternalServerError,
				http.StatusInternalServerError)
		}

		subStock := stock.Stock - orderProduct.Quantity

		stockModel := models.StockModel{
			ID:    stock.ID,
			Stock: subStock,
			UpdatedBy: uuid.NullUUID{
				UUID:  uuid.FromStringOrNil(ctx.Value("user_id").(string)),
				Valid: true,
			},
		}

		err = stockModel.Update(ctx, s.db)

		if err != nil {
			return nil, helpers.ErrorWrap(err, s.name, "Order/stockModel.Update",
				helpers.InternalServerError,
				http.StatusInternalServerError)
		}

		stocks, err := models.GetAllStockByProductID(ctx, s.db, stock.ProductID)
		if err != nil {
			return nil, helpers.ErrorWrap(err, s.name, "Add/GetAllStockByProductID", helpers.InternalServerError,
				http.StatusInternalServerError)
		}

		var totalStock uint
		for _, stock = range stocks {
			totalStock = totalStock + stock.Stock
		}

		productStock := models.ProductModel{
			ID:    stock.ProductID,
			Stock: totalStock,
			UpdatedBy: uuid.NullUUID{
				UUID:  uuid.FromStringOrNil(ctx.Value("user_id").(string)),
				Valid: true,
			},
		}

		err = productStock.StockUpdate(ctx, s.db)

		if err != nil {
			return nil, helpers.ErrorWrap(err, s.name, "Add/StockUpdate", helpers.InternalServerError,
				http.StatusInternalServerError)
		}

		product, err := models.GetOneProduct(ctx, s.db, orderProduct.ID)

		if err != nil {
			return nil, helpers.ErrorWrap(err, s.name, "Order/GetOneProduct",
				helpers.InternalServerError,
				http.StatusInternalServerError)
		}

		subTotal := product.Price.Mul(decimal.NewFromInt(int64(orderProduct.Quantity)))

		orderProduct := models.OrderProductModel{
			OrderID:   order.ID,
			ProductID: orderProduct.ID,
			Quantity:  orderProduct.Quantity,
			SubTotal:  subTotal,
			CreatedBy: uuid.FromStringOrNil(ctx.Value("user_id").(string)),
		}

		err = orderProduct.Insert(ctx, s.db)
		if err != nil {
			return nil, helpers.ErrorWrap(err, s.name, "Order/orderProduct.Insert",
				helpers.InternalServerError,
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

	configuration, err := models.GetConfiguration(ctx, s.db)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Order/GetConfiguration", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	fmt.Println(totalPrice)

	orderUpdate := models.OrderModel{
		ID:         order.ID,
		TotalPrice: totalPrice.Add(configuration.DeliveryFee),
		UpdatedBy: uuid.NullUUID{
			UUID:  uuid.FromStringOrNil(ctx.Value("user_id").(string)),
			Valid: true,
		},
	}

	fmt.Println(orderUpdate.TotalPrice)
	fmt.Println(order.ID)
	err = orderUpdate.UpdatePrice(ctx, s.db)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Order/UpdatePrice", helpers.InternalServerError,
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
