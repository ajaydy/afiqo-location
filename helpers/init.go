package helpers

import "github.com/gomodule/redigo/redis"

var (
	logger    *Logger
	cachePool *redis.Pool
)

func Init(log *Logger, cache *redis.Pool) {
	logger = log
	cachePool = cache
}
