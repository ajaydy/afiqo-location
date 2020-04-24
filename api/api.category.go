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
	CategoryModule struct {
		db     *sql.DB
		cache  *redis.Pool
		logger *helpers.Logger
		name   string
	}

	CategoryDetailParam struct {
		ID uuid.UUID `json:"id"`
	}

	CategoryAddParam struct {
		Name        string `json:"name" validator:"required"`
		Description string `json:"description" validator:"required"`
	}

	CategoryUpdateParam struct {
		ID          uuid.UUID `json:"id"`
		Name        string    `json:"name" validator:"required"`
		Description string    `json:"description" validator:"required"`
	}

	CategoryDeleteParam struct {
		ID uuid.UUID `json:"id"`
	}
)

func NewCategoryModule(db *sql.DB, cache *redis.Pool, logger *helpers.Logger) *CategoryModule {
	return &CategoryModule{
		db:     db,
		cache:  cache,
		logger: logger,
		name:   "module/category",
	}
}

func (s CategoryModule) Detail(ctx context.Context, param CategoryDetailParam) (interface{}, *helpers.Error) {
	category, err := models.GetOneCategory(ctx, s.db, param.ID)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Detail/GetOneCategory", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return category.Response(), nil
}

func (s CategoryModule) List(ctx context.Context, filter helpers.Filter) (interface{}, *helpers.Error) {

	categories, err := models.GetAllCategory(ctx, s.db, filter)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "List/GetAllCategory", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	var categoryResponse []models.CategoryResponse
	for _, category := range categories {
		categoryResponse = append(categoryResponse, category.Response())
	}

	return categoryResponse, nil
}

func (s CategoryModule) Add(ctx context.Context, param CategoryAddParam) (interface{}, *helpers.Error) {

	category := models.CategoryModel{
		Name:        param.Name,
		Description: param.Description,
		CreatedBy:   uuid.FromStringOrNil(ctx.Value("user_id").(string)),
	}

	err := category.Insert(ctx, s.db)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Add/Insert", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return category.Response(), nil
}

func (s CategoryModule) Update(ctx context.Context, param CategoryUpdateParam) (interface{}, *helpers.Error) {

	category := models.CategoryModel{
		ID:          param.ID,
		Name:        param.Name,
		Description: param.Description,
		UpdatedBy: uuid.NullUUID{
			UUID:  uuid.FromStringOrNil(ctx.Value("user_id").(string)),
			Valid: true,
		},
	}

	err := category.Update(ctx, s.db)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Update/Update", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return category.Response(), nil

}

func (s CategoryModule) Delete(ctx context.Context, param CategoryDeleteParam) (interface{}, *helpers.Error) {

	category := models.CategoryModel{
		ID: param.ID,
		UpdatedBy: uuid.NullUUID{
			UUID:  uuid.FromStringOrNil(ctx.Value("user_id").(string)),
			Valid: true,
		},
	}

	err := category.Delete(ctx, s.db)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Delete/Delete", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return nil, nil

}
