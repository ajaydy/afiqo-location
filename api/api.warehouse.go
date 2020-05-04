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
	WarehouseModule struct {
		db     *sql.DB
		cache  *redis.Pool
		logger *helpers.Logger
		name   string
	}

	WarehouseDetailParam struct {
		ID uuid.UUID `json:"id"`
	}

	WarehouseAddParam struct {
		Name    string `json:"name" validate:"required"`
		Address string `json:"address" validate:"required"`
		PhoneNo string `json:"phone_no" validate:"required"`
	}

	WarehouseUpdateParam struct {
		ID      uuid.UUID `json:"id"`
		Name    string    `json:"name" validate:"max=20,min=4,required"`
		Address string    `json:"address" validate:"required"`
		PhoneNo string    `json:"phone_no" validate:"required"`
	}

	WarehouseDeleteParam struct {
		ID uuid.UUID `json:"id"`
	}
)

func NewWarehouseModule(db *sql.DB, cache *redis.Pool, logger *helpers.Logger) *WarehouseModule {
	return &WarehouseModule{
		db:     db,
		cache:  cache,
		logger: logger,
		name:   "module/warehouse",
	}
}

func (s WarehouseModule) Detail(ctx context.Context, param WarehouseDetailParam) (interface{}, *helpers.Error) {
	warehouse, err := models.GetOneWarehouse(ctx, s.db, param.ID)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Detail/GetOneWarehouse", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return warehouse.Response(), nil
}

func (s WarehouseModule) List(ctx context.Context, filter helpers.Filter) (interface{}, *helpers.Error) {

	warehouses, err := models.GetAllWarehouse(ctx, s.db, filter)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "List/GetAllWarehouse", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	var warehouseResponse []models.WarehouseResponse
	for _, warehouse := range warehouses {
		warehouseResponse = append(warehouseResponse, warehouse.Response())
	}

	return warehouseResponse, nil
}

func (s WarehouseModule) Add(ctx context.Context, param WarehouseAddParam) (interface{}, *helpers.Error) {

	warehouse := models.WarehouseModel{
		Name:      param.Name,
		Address:   param.Address,
		PhoneNo:   param.PhoneNo,
		CreatedBy: uuid.FromStringOrNil(ctx.Value("user_id").(string)),
	}

	err := warehouse.Insert(ctx, s.db)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Add/Insert", helpers.InternalServerError,
			http.StatusInternalServerError)
	}
	return warehouse.Response(), nil
}

func (s WarehouseModule) Update(ctx context.Context, param WarehouseUpdateParam) (interface{}, *helpers.Error) {

	warehouse := models.WarehouseModel{
		ID:      param.ID,
		Name:    param.Name,
		Address: param.Address,
		PhoneNo: param.PhoneNo,
		UpdatedBy: uuid.NullUUID{
			UUID:  uuid.FromStringOrNil(ctx.Value("user_id").(string)),
			Valid: true,
		},
	}

	err := warehouse.Update(ctx, s.db)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Update/Update", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return warehouse.Response(), nil

}

func (s WarehouseModule) Delete(ctx context.Context, param WarehouseDeleteParam) (interface{}, *helpers.Error) {

	warehouse := models.WarehouseModel{
		ID: param.ID,
		UpdatedBy: uuid.NullUUID{
			UUID:  uuid.FromStringOrNil(ctx.Value("user_id").(string)),
			Valid: true,
		},
	}

	err := warehouse.Delete(ctx, s.db)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Delete/Delete", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return nil, nil

}
