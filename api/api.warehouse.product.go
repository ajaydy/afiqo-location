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
	WarehouseProductModule struct {
		db     *sql.DB
		cache  *redis.Pool
		logger *helpers.Logger
		name   string
	}

	WarehouseProductDetailParam struct {
		ID uuid.UUID `json:"id"`
	}

	StockAddParam struct {
		WarehouseID uuid.UUID `json:"warehouse_id"`
		ProductID   uuid.UUID `json:"product_id"`
		Stock       uint      `json:"stock"`
	}
)

func NewWarehouseProductModule(db *sql.DB, cache *redis.Pool, logger *helpers.Logger) *WarehouseProductModule {
	return &WarehouseProductModule{
		db:     db,
		cache:  cache,
		logger: logger,
		name:   "module/warehouseProduct",
	}
}

func (s WarehouseProductModule) Detail(ctx context.Context, param WarehouseProductDetailParam) (
	interface{}, *helpers.Error) {

	warehouseProduct, err := models.GetOneWarehouseProduct(ctx, s.db, param.ID)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Detail/GetOneWarehouseProduct", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	response, err := warehouseProduct.Response(ctx, s.db, s.logger)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Detail/Response", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return response, nil
}

func (s WarehouseProductModule) List(ctx context.Context, filter helpers.Filter) (interface{}, *helpers.Error) {

	warehouseProducts, err := models.GetAllWarehouseProduct(ctx, s.db, filter)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "List/GetAllWarehouseProduct", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	var warehouseProductResponse []models.WarehouseProductResponse
	for _, warehouseProduct := range warehouseProducts {
		response, err := warehouseProduct.Response(ctx, s.db, s.logger)
		if err != nil {
			return nil, helpers.ErrorWrap(err, s.name, "List/Response", helpers.InternalServerError,
				http.StatusInternalServerError)
		}

		warehouseProductResponse = append(warehouseProductResponse, response)
	}

	return warehouseProductResponse, nil
}

func (s WarehouseProductModule) Add(ctx context.Context, param StockAddParam) (interface{}, *helpers.Error) {

	warehouseProduct := models.WarehouseProductModel{
		ProductID:   param.ProductID,
		WarehouseID: param.WarehouseID,
		Stock:       param.Stock,
		CreatedBy:   uuid.FromStringOrNil(ctx.Value("user_id").(string)),
	}

	err := warehouseProduct.Insert(ctx, s.db)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Add/Insert", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	product, err := models.GetOneProduct(ctx, s.db, warehouseProduct.ProductID)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Add/GetOneProduct", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	totalStock := product.Stock + param.Stock

	productStock := models.ProductModel{
		ID:    warehouseProduct.ProductID,
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

	response, err := warehouseProduct.Response(ctx, s.db, s.logger)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Add/Response", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return response, nil

}
