package sprjob

import (
	"fmt"
	"github.com/daqnext/SPR-go/goredis"
	"github.com/go-redis/redis/v8"
	"log"
	"math/rand"
	"time"
)

type SprJob struct {
	JobName          string
	IsMaster         bool
	JobRand          string
	KeepIntervalSec  int
	ApplyIntervalSec int
	KeepTicker       *time.Ticker
	ApplyTicker      *time.Ticker
	KeepTickerStop   chan struct{}
	ApplyTickerStop  chan struct{}
}

func New(name string, keepIntervalSec int, applyIntervalSec int) *SprJob {
	s := &SprJob{
		JobName:          "spr:" + name,
		IsMaster:         false,
		JobRand:          fmt.Sprintf("%d", rand.Intn(100000000)+1),
		KeepIntervalSec:  keepIntervalSec,
		ApplyIntervalSec: applyIntervalSec,
		KeepTickerStop:   make(chan struct{}, 1),
		ApplyTickerStop:  make(chan struct{}, 1),
	}
	return s
}

func (s *SprJob) StartLoop() {
	s.keepMaster()
	s.KeepTicker = time.NewTicker(time.Second * time.Duration(s.KeepIntervalSec))
	go func() {
		for {
			select {
			case <-s.KeepTicker.C:
				//loop
				s.keepMaster()
			case <-s.KeepTickerStop:
				return
			}
		}
	}()

	s.applyMaster()
	s.ApplyTicker = time.NewTicker(time.Second * time.Duration(s.ApplyIntervalSec))
	go func() {
		for {
			select {
			case <-s.ApplyTicker.C:
				//loop
				s.applyMaster()
			case <-s.ApplyTickerStop:
				return
			}
		}
	}()
}

func (s *SprJob) StopLoop() {
	if s.KeepTicker != nil {
		s.KeepTicker.Stop()
	}
	if s.ApplyTicker != nil {
		s.ApplyTicker.Stop()
	}
	s.KeepTickerStop <- struct{}{}
	s.ApplyTickerStop <- struct{}{}
}

func (s *SprJob) applyMaster() {
	log.Println(s.JobName, "ApplyMaster job run")

	// already master
	if s.IsMaster {
		return
	}

	if goredis.RedisClient == nil {
		return
	}

	_, err := goredis.RedisClient.Get(goredis.Ctx, s.JobName).Result()
	//not nil, there is another master
	if err != redis.Nil {
		return
	}

	success, err := goredis.RedisClient.SetNX(goredis.Ctx, s.JobName, s.JobRand, time.Second*90).Result()
	if err != nil {
		return
	}
	if !success {
		return
	}
	s.IsMaster = true
}

func (s *SprJob) keepMaster() {
	log.Println(s.JobName, "KeepMaster job run")

	//not master
	if !s.IsMaster {
		return
	}

	if goredis.RedisClient == nil {
		return
	}

	value, err := goredis.RedisClient.Get(goredis.Ctx, s.JobName).Result()
	//get value error
	if err != nil {
		s.IsMaster = false
		return
	}

	//value error
	if value != s.JobRand {
		s.IsMaster = false
		return
	}

	//value==jobRand
	//keep master token
	goredis.RedisClient.Expire(goredis.Ctx, s.JobName, time.Second*90)
}
