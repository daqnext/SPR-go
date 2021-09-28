package goredis

import (
	"context"
	"log"
	"strconv"

	"github.com/go-redis/redis/v8"
)

const SPR_REDIS_DB = 0

var Ctx = context.Background()

var RedisClient *redis.Client

func InitRedisClient(addr string, port int, userName string, password string) error {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr + ":" + strconv.Itoa(port),
		Username: userName,
		Password: password,     // no password set
		DB:       SPR_REDIS_DB, // use default DB
	})

	_, err := rdb.Ping(Ctx).Result()
	if err != nil {
		log.Println("Redis connect failed")
		return err
	}
	log.Println("Redis connect success")
	RedisClient = rdb
	return nil
}
