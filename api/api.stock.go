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
	StockModule struct {
		db     *sql.DB
		cache  *redis.Pool
		logger *helpers.Logger
		name   string
	}

	StockDetailParam struct {
		ID uuid.UUID `json:"id"`
	}

	StockAddParam struct {
		WarehouseID uuid.UUID `json:"warehouse_id" validate:"required"`
		ProductID   uuid.UUID `json:"product_id" validate:"required"`
		Stock       uint      `json:"stock" validate:"required"`
	}

	StockUpdateParam struct {
		ID    uuid.UUID `json:"id"`
		Stock uint      `json:"stock" validate:"required"`
	}
)

func NewStockModule(db *sql.DB, cache *redis.Pool, logger *helpers.Logger) *StockModule {
	return &StockModule{
		db:     db,
		cache:  cache,
		logger: logger,
		name:   "module/stock",
	}
}

func (s StockModule) Detail(ctx context.Context, param StockDetailParam) (
	interface{}, *helpers.Error) {

	stock, err := models.GetOneStock(ctx, s.db, param.ID)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Detail/GetOneStock", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	response, err := stock.Response(ctx, s.db, s.logger)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Detail/Response", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return response, nil
}

func (s StockModule) List(ctx context.Context, filter helpers.Filter) (interface{}, *helpers.Error) {

	stocks, err := models.GetAllStock(ctx, s.db, filter)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "List/GetAllStock", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	var stockResponses []models.StockResponse
	for _, stock := range stocks {
		response, err := stock.Response(ctx, s.db, s.logger)
		if err != nil {
			return nil, helpers.ErrorWrap(err, s.name, "List/Response", helpers.InternalServerError,
				http.StatusInternalServerError)
		}

		stockResponses = append(stockResponses, response)
	}

	return stockResponses, nil
}

func (s StockModule) ListBySupplierID(ctx context.Context, filter helpers.Filter, param SupplierDataParam) (
	interface{}, *helpers.Error) {

	stocks, err := models.GetAllStockBySupplierID(ctx, s.db, filter, param.ID)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "ListBySupplierID/GetAllStock", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	var stockResponses []models.StockResponse
	for _, stock := range stocks {
		response, err := stock.Response(ctx, s.db, s.logger)
		if err != nil {
			return nil, helpers.ErrorWrap(err, s.name, "ListBySupplierID/Response", helpers.InternalServerError,
				http.StatusInternalServerError)
		}

		stockResponses = append(stockResponses, response)
	}

	return stockResponses, nil
}

func (s StockModule) Add(ctx context.Context, param StockAddParam) (interface{}, *helpers.Error) {

	stock := models.StockModel{
		ProductID:   param.ProductID,
		WarehouseID: param.WarehouseID,
		Stock:       param.Stock,
		CreatedBy:   uuid.FromStringOrNil(ctx.Value("user_id").(string)),
	}

	err := stock.Insert(ctx, s.db)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Add/Insert", helpers.InternalServerError,
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

	response, err := stock.Response(ctx, s.db, s.logger)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Add/Response", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return response, nil

}

func (s StockModule) Update(ctx context.Context, param StockUpdateParam) (interface{}, *helpers.Error) {

	stock := models.StockModel{
		ID:    param.ID,
		Stock: param.Stock,
		UpdatedBy: uuid.NullUUID{
			UUID:  uuid.FromStringOrNil(ctx.Value("user_id").(string)),
			Valid: true,
		},
	}

	err := stock.Update(ctx, s.db)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Update/Insert", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	stock, err = models.GetOneStock(ctx, s.db, stock.ID)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Update/GetOneStock", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	stocks, err := models.GetAllStockByProductID(ctx, s.db, stock.ProductID)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Update/GetAllStockByProductID", helpers.InternalServerError,
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
		return nil, helpers.ErrorWrap(err, s.name, "Update/StockUpdate", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	response, err := stock.Response(ctx, s.db, s.logger)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Update/Response", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return response, nil

}
