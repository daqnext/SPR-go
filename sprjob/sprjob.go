package sprjob

import (
	"fmt"
	"github.com/daqnext/SPR-go/goredis"
	"github.com/go-redis/redis/v8"
	"log"
	"math/rand"
	"time"
)

const LoopIntervalSec = 15
const MasterKeepTime = 90

type SprJob struct {
	JobName         string
	IsMaster        bool
	JobRand         string
	LoopIntervalSec int
	StopFlag        bool
}

func New(name string) *SprJob {
	s := &SprJob{
		JobName:         "spr:" + name,
		IsMaster:        false,
		JobRand:         fmt.Sprintf("%d", rand.Intn(100000000)+1),
		LoopIntervalSec: LoopIntervalSec,
		StopFlag:        false,
	}
	return s
}

func (s *SprJob) StartLoop() {
	//s.loop()
	go func() {
		for {
			if s.StopFlag {
				return
			}
			s.loop()
			time.Sleep(time.Second * time.Duration(s.LoopIntervalSec))
		}
	}()
}

func (s *SprJob) StopLoop() {
	s.StopFlag = true
}

func (s *SprJob) loop() {
	log.Println(s.JobName, "loop job run")
	if goredis.RedisClient == nil {
		return
	}

	//check jobname in redis
	value, err := goredis.RedisClient.Get(goredis.Ctx, s.JobName).Result()
	//any err
	if err != nil && err != redis.Nil {
		//log.Println(err)
		return
	}

	if err == redis.Nil {
		//if no value
		success, err := goredis.RedisClient.SetNX(goredis.Ctx, s.JobName, s.JobRand, time.Second*time.Duration(MasterKeepTime)).Result()
		if err != nil {
			return
		}
		if !success {
			return
		}
		s.IsMaster = true

	} else {
		//value error
		if value != s.JobRand {
			s.IsMaster = false
			return
		}

		//value==jobRand
		//keep master token
		s.IsMaster = true
		goredis.RedisClient.Expire(goredis.Ctx, s.JobName, time.Second*time.Duration(MasterKeepTime))
	}
}
