package chachingservice

import (
	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

func SetRedies(rdb *redis.Client) {
	RDB = rdb
}
