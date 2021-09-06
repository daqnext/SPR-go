package SPR_go

import (
	"errors"
	"github.com/daqnext/SPR-go/goredis"
	"github.com/daqnext/SPR-go/sprjob"
	"math/rand"
	"sync"
	"time"
)

type SprJobMgr struct {
	jobMap sync.Map
}

type RedisConfig struct {
	Addr     string
	Port     int
	Db       int
	UserName string
	Password string
}

func New() *SprJobMgr {
	rand.Seed(time.Now().UnixNano())
	sMgr := &SprJobMgr{}
	return sMgr
}

func (smgr *SprJobMgr) InitRedis(config RedisConfig) {
	goredis.InitRedisClient(config.Addr, config.Port, config.Db, config.UserName, config.Password)
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
