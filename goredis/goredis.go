package goredis

import (
	"context"
	"log"
	"strconv"

	"github.com/go-redis/redis/v8"
)

const SPR_REDIS_DB = 0

var Ctx = context.Background()

var RedisClient *redis.ClusterClient

func InitRedisClient(addr string, port int, userName string, password string) error {

	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    []string{addr + ":" + strconv.Itoa(port)},
		Username: userName,
		Password: password,
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
