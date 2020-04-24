package session

import (
	"afiqo-location/helpers"
	"database/sql"
	"github.com/gomodule/redigo/redis"
)

var (
	dbPool    *sql.DB
	cachePool *redis.Pool
	logger    *helpers.Logger
)

func Init(db *sql.DB, cache *redis.Pool, log *helpers.Logger) {
	dbPool = db
	cachePool = cache
	logger = log
}
