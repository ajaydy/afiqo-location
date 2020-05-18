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
	ConfigurationModule struct {
		db     *sql.DB
		cache  *redis.Pool
		logger *helpers.Logger
		name   string
	}

	ConfigurationUpdateParam struct {
		DeliveryFee decimal.Decimal `json:"delivery_fee" validate:"required"`
	}
)

func NewConfigurationModule(db *sql.DB, cache *redis.Pool, logger *helpers.Logger) *ConfigurationModule {
	return &ConfigurationModule{
		db:     db,
		cache:  cache,
		logger: logger,
		name:   "module/configuration",
	}
}

func (s ConfigurationModule) GetFee(ctx context.Context) (decimal.Decimal, *helpers.Error) {

	fee, err := models.GetConfiguration(ctx, s.db)
	if err != nil {
		return decimal.Decimal{}, helpers.ErrorWrap(err, s.name, "GetFee/GetConfiguration", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return fee.DeliveryFee, nil

}

func (s ConfigurationModule) Update(ctx context.Context, param ConfigurationUpdateParam) (interface{}, *helpers.Error) {

	configuration := models.ConfigurationModel{
		DeliveryFee: param.DeliveryFee,
		UpdatedBy: uuid.NullUUID{
			UUID:  uuid.FromStringOrNil(ctx.Value("user_id").(string)),
			Valid: true,
		},
	}

	err := configuration.Update(ctx, s.db)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Update/Update", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return configuration.Response(), nil

}
