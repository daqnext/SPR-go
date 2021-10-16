# SPR-go

Make sure specific job can only be run over 1 process among different machines and different process. If some process is killed or stucked, the job will be switched to other process on some machine.

### usage
```
go get github.com/daqnext/SPR-go
```

```go


package SPR_go

import (
	"errors"
	"math/rand"
	"sync"
	"time"

	localLog "github.com/daqnext/LocalLog/log"
	"github.com/daqnext/SPR-go/goredis"
	"github.com/daqnext/SPR-go/sprjob"
)

type SprJobMgr struct {
	jobMap sync.Map
	llog   *localLog.LocalLog
}

type RedisConfig struct {
	Addr     string
	Port     int
	UserName string
	Password string
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func New(config RedisConfig, localLogger *localLog.LocalLog) (*SprJobMgr, error) {
	err := goredis.InitRedisClient(config.Addr, config.Port, config.UserName, config.Password, localLogger)
	if err != nil {
		return nil, errors.New("redis connect error")
	}
	sMgr := &SprJobMgr{
		llog: localLogger,
	}
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
	job.StartLoop(smgr.llog)
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



```