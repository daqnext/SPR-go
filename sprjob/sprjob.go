package sprjob

import (
	"fmt"
	"github.com/daqnext/SPR-go/goredis"
	"github.com/daqnext/go-smart-routine/sr"
	"github.com/go-redis/redis/v8"
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
	LastRuntime     int64
}

func New(name string) *SprJob {
	s := &SprJob{
		JobName:         "spr:" + name,
		IsMaster:        false,
		JobRand:         fmt.Sprintf("%d", rand.Intn(100000000)+1),
		LoopIntervalSec: LoopIntervalSec,
		StopFlag:        false,
		LastRuntime:     0,
	}
	return s
}

func (s *SprJob) StartLoop() {
	//s.loop()
	sr.New_Panic_Redo(func() {
		for {
			if s.StopFlag {
				return
			}
			nowUnixTime := time.Now().Unix()
			toSleepSecs := s.LastRuntime + int64(s.LoopIntervalSec) - nowUnixTime
			if toSleepSecs <= 0 {
				s.LastRuntime = nowUnixTime
				s.loop()
			} else {
				time.Sleep(time.Duration(toSleepSecs) * time.Second)
			}
		}
	}).Start()
}

func (s *SprJob) StopLoop() {
	s.StopFlag = true
	s.IsMaster = false
}

func (s *SprJob) loop() {
	//log.Println(s.JobName, "loop job run")
	if goredis.RedisClient == nil {
		s.IsMaster = false
		return
	}

	//check jobname in redis
	value, err := goredis.RedisClient.Get(goredis.Ctx, s.JobName).Result()

	//get value
	if err == nil {
		//value error
		if value != s.JobRand {
			s.IsMaster = false
			return
		}

		//value==jobRand
		//keep master token
		s.IsMaster = true
		goredis.RedisClient.Expire(goredis.Ctx, s.JobName, time.Second*time.Duration(MasterKeepTime))

	} else if err == redis.Nil {
		//if no value
		success, err := goredis.RedisClient.SetNX(goredis.Ctx, s.JobName, s.JobRand, time.Second*time.Duration(MasterKeepTime)).Result()
		if err != nil || !success {
			s.IsMaster = false
			return
		}
		s.IsMaster = true

	} else {
		//other err
		s.IsMaster = false
		return

	}
}
