package SPR_go

import (
	"errors"
	"math/rand"
	"sync"
	"time"

	"github.com/daqnext/SPR-go/goredis"
	"github.com/daqnext/SPR-go/sprjob"
)

type SprJobMgr struct {
	jobMap sync.Map
}

type RedisConfig struct {
	Addr     string
	Port     int
	UserName string
	Password string
}

func New(config RedisConfig) (*SprJobMgr, error) {
	err := goredis.InitRedisClient(config.Addr, config.Port, config.UserName, config.Password)
	if err != nil {
		return nil, errors.New("redis connect error")
	}
	rand.Seed(time.Now().UnixNano())
	sMgr := &SprJobMgr{}
	return sMgr, nil
}

func (smgr *SprJobMgr) AddJobName(jobName string) error {
	_, exist := smgr.jobMap.Load(jobName)
	if exist {
		return errors.New("job already exist")
	}
	//new job
	job := sprjob.New(jobName)
	smgr.jobMap.Store(jobName, job)
	//start loop
	job.StartLoop()
	return nil
}

func (smgr *SprJobMgr) RemoveJobName(jobName string) {
	job, exist := smgr.jobMap.Load(jobName)
	if !exist {
		return
	}
	//stop
	job.(*sprjob.SprJob).StopLoop()
	//delete
	smgr.jobMap.Delete(jobName)
}

func (smgr *SprJobMgr) IsMaster(jobName string) bool {
	job, exist := smgr.jobMap.Load(jobName)
	if !exist {
		return false
	}
	return job.(*sprjob.SprJob).IsMaster
}
