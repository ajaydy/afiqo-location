package api

import (
	"afiqo-location/helpers"
	"afiqo-location/models"
	"context"
	"database/sql"
	"github.com/gomodule/redigo/redis"
	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"net/http"
)

type (
	ProductModule struct {
		db     *sql.DB
		cache  *redis.Pool
		logger *helpers.Logger
		name   string
	}

	ProductDetailParam struct {
		ID uuid.UUID `json:"id"`
	}

	ProductAddParam struct {
		SupplierID  uuid.UUID       `json:"supplier_id" validate:"required"`
		CategoryID  uuid.UUID       `json:"category_id" validate:"required"`
		Name        string          `json:"name" validate:"required"`
		Price       decimal.Decimal `json:"price" validate:"required"`
		Description string          `json:"description" validate:"required"`
	}

	ProductUpdateParam struct {
		ID          uuid.UUID       `json:"id"`
		Name        string          `json:"name" validate:"required"`
		Price       decimal.Decimal `json:"price" validate:"required"`
		Description string          `json:"description" validate:"required"`
	}

	ProductDeleteParam struct {
		ID uuid.UUID `json:"id"`
	}
)

func NewProductModule(db *sql.DB, cache *redis.Pool, logger *helpers.Logger) *ProductModule {
	return &ProductModule{
		db:     db,
		cache:  cache,
		logger: logger,
		name:   "module/product",
	}
}

func (s ProductModule) Detail(ctx context.Context, param ProductDetailParam) (interface{}, *helpers.Error) {
	product, err := models.GetOneProduct(ctx, s.db, param.ID)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Detail/GetOneProduct", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	response, err := product.Response(ctx, s.db, s.logger)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Detail/Response", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return response, nil
}

func (s ProductModule) List(ctx context.Context, filter helpers.Filter) (interface{}, *helpers.Error) {

	products, err := models.GetAllProduct(ctx, s.db, filter)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "List/GetAllProduct", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	var productResponse []models.ProductResponse
	for _, product := range products {
		response, err := product.Response(ctx, s.db, s.logger)
		if err != nil {
			return nil, helpers.ErrorWrap(err, s.name, "List/Response", helpers.InternalServerError,
				http.StatusInternalServerError)
		}

		productResponse = append(productResponse, response)
	}

	return productResponse, nil
}

func (s ProductModule) Add(ctx context.Context, param ProductAddParam) (interface{}, *helpers.Error) {

	product := models.ProductModel{
		SupplierID:  param.SupplierID,
		CategoryID:  param.CategoryID,
		Name:        param.Name,
		Stock:       0,
		Price:       param.Price,
		Description: param.Description,
		CreatedBy:   uuid.FromStringOrNil(ctx.Value("user_id").(string)),
	}

	err := product.Insert(ctx, s.db)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Add/Insert", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	response, err := product.Response(ctx, s.db, s.logger)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Add/Response", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return response, nil
}

func (s ProductModule) Update(ctx context.Context, param ProductUpdateParam) (interface{}, *helpers.Error) {

	product := models.ProductModel{
		ID:          param.ID,
		Name:        param.Name,
		Price:       param.Price,
		Description: param.Description,
		UpdatedBy: uuid.NullUUID{
			UUID:  uuid.FromStringOrNil(ctx.Value("user_id").(string)),
			Valid: true,
		},
	}

	err := product.Update(ctx, s.db)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Update/Update", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	response, err := product.Response(ctx, s.db, s.logger)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Update/Response", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return response, nil

}

func (s ProductModule) Delete(ctx context.Context, param ProductDeleteParam) (interface{}, *helpers.Error) {

	product := models.ProductModel{
		ID: param.ID,
		UpdatedBy: uuid.NullUUID{
			UUID:  uuid.FromStringOrNil(ctx.Value("user_id").(string)),
			Valid: true,
		},
	}

	err := product.Delete(ctx, s.db)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Delete/Delete", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return nil, nil

}
