package api

import (
	"afiqo-location/helpers"
	"database/sql"
	"github.com/gomodule/redigo/redis"
)

type (
	OrderModule struct {
		db     *sql.DB
		cache  *redis.Pool
		logger *helpers.Logger
		name   string
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
