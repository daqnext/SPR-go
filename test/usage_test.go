package test

import (
	SPR_go "github.com/daqnext/SPR-go"
	"log"
	"testing"
	"time"
)

func Test_usage(t *testing.T) {
	//init redis with config
	//If connect to redis failed, all the job will not be the master

	//type RedisConfig struct{
	//	Addr string
	//	Port int
	//	Db int
	//	UserName string
	//	Password string
	//}
	SPR_go.Smgr.InitRedis(SPR_go.RedisConfig{
		Addr:     "127.0.0.1",
		Port:     6379,
		Db:       5,
		Password: "123456",
	})

	//add job with unique job name which used in redis
	//the process with same job name will scramble for the master token
	err := SPR_go.Smgr.AddJobName("testjob")
	if err != nil {
		log.Println(err)
	}
	err = SPR_go.Smgr.AddJobName("testjob2")
	if err != nil {
		log.Println(err)
	}

	// use function IsMaster("jobName") to check whether the process get the master token or not
	// if return true means get the master token
	go func() {
		for {
			time.Sleep(time.Second)
			log.Println("testjob is master:", SPR_go.Smgr.IsMaster("testjob"))
			log.Println("testjob2 is master:", SPR_go.Smgr.IsMaster("testjob2"))
		}
	}()

	// use function RemoveJobName("jobName") to remove the job
	// removed job always return false when use IsMaster("jobName")
	time.AfterFunc(time.Second*30, func() {
		SPR_go.Smgr.RemoveJobName("testjob2")
	})

	time.Sleep(1 * time.Hour)

}
