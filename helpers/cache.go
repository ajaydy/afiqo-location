package helpers

import (
	"context"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"strconv"
	"time"
)

type CacheOptions struct {
	Host        string
	Port        int
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout int
	Enabled     bool
}

var pool *redis.Pool

func ConnectToCache(cacheOptions CacheOptions) *redis.Pool {
	if pool == nil {
		pool = &redis.Pool{
			MaxIdle:     cacheOptions.MaxIdle,
			MaxActive:   cacheOptions.MaxActive,
			IdleTimeout: time.Duration(cacheOptions.IdleTimeout) * time.Second,
			Dial: func() (redis.Conn, error) {
				address := fmt.Sprintf("%s:%d", cacheOptions.Host, cacheOptions.Port)
				c, err := redis.Dial("tcp", address)
				if err != nil {
					return nil, err
				}
				// Do authentication process if password not empty.
				if cacheOptions.Password != "" {
					if _, err := c.Do("AUTH", cacheOptions.Password); err != nil {
						c.Close()
						return nil, err
					}
				}
				return c, nil
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				if time.Since(t) < time.Minute {
					return nil
				}
				_, err := c.Do("PING")
				return err
			},
			Wait:            true,
			MaxConnLifetime: 15 * time.Minute,
		}
		return pool
	}
	return pool
}

func SetDataToCache(ctx context.Context, id, value string) error {
	conn, err := cachePool.Dial()

	if err != nil {
		return err
	}

	_, err = conn.Do("SET", id, value)
	if err != nil {
		return err
	}

	return nil
}

func SetDataToCacheWithExpiry(ctx context.Context, id, value string, expiryTime int) error {
	conn, err := cachePool.Dial()

	if err != nil {
		return err
	}

	_, err = conn.Do("SETEX", id, strconv.Itoa(expiryTime), value)

	if err != nil {
		return err
	}

	return nil
}

func GetDataFromCache(ctx context.Context, id string) (string, error) {
	conn, err := cachePool.Dial()

	if err != nil {
		return "", err
	}

	data, err := redis.String(conn.Do("GET", id))
	if err != nil {
		return "", err
	}

	return data, nil

}

func GetKeysFromCache(ctx context.Context) ([]string, error) {
	conn, err := cachePool.Dial()

	if err != nil {
		return nil, err
	}

	data, err := redis.Strings(conn.Do("KEYS", "*"))
	if err != nil {
		return nil, err
	}

	return data, nil

}

func GetKeysFromCacheWithPrefix(ctx context.Context, prefix string) ([]string, error) {
	conn, err := cachePool.Dial()

	if err != nil {
		return nil, err
	}

	data, err := redis.Strings(conn.Do("KEYS", fmt.Sprintf(`%s*`, prefix)))
	if err != nil {
		return nil, err
	}

	return data, nil

}

func SetCacheExpiry(ctx context.Context, id string, expiryTime int) error {
	conn, err := cachePool.Dial()

	if err != nil {
		return err
	}
	_, err = conn.Do("EXPIRE", id, strconv.Itoa(expiryTime))
	if err != nil {
		return err
	}

	return nil
}

func DeleteCache(ctx context.Context, id string) error {
	conn, err := cachePool.Dial()

	if err != nil {
		return err
	}
	_, err = conn.Do("DEL", id)
	if err != nil {
		return err
	}

	return nil
}
