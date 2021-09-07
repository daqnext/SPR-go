package goredis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"strconv"
)

var Ctx = context.Background()

var RedisClient *redis.Client

func InitRedisClient(addr string, port int, db int, userName string, password string) error {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr + ":" + strconv.Itoa(port),
		Username: userName,
		Password: password, // no password set
		DB:       db,       // use default DB
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
