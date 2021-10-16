package goredis

import (
	"context"
	localLog "github.com/daqnext/LocalLog/log"
	"strconv"

	"github.com/go-redis/redis/v8"
)

const SPR_REDIS_DB = 0

var Ctx = context.Background()

var RedisClient *redis.ClusterClient
var lg *localLog.LocalLog

func InitRedisClient(addr string, port int, userName string, password string, llog *localLog.LocalLog) error {
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    []string{addr + ":" + strconv.Itoa(port)},
		Username: userName,
		Password: password,
	})
	lg = llog

	_, err := rdb.Ping(Ctx).Result()
	if err != nil {
		if lg != nil {
			lg.Errorln("SPR-go Redis connect failed")
		}
		return err
	}
	if lg != nil {
		lg.Println("SPR-go Redis connect success")
	}
	RedisClient = rdb
	return nil
}
