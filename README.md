# SPR-go

Make sure specific job can only be run over 1 process among different machines and different process. If some process is killed or stucked, the job will be switched to other process on some machine.

### usage
```
go get github.com/daqnext/SPR-go
```

```go
import (
    localLog "github.com/daqnext/LocalLog/log"
    SPR_go "github.com/daqnext/SPR-go"
    "log"
    "time"
)

func Test_usage(t *testing.T) {
    //new instance
    //init redis with config
    //If connect to redis failed, return err

    //type RedisConfig struct{
    //	Addr string
    //	Port int
    //	Db int
    //	UserName string
    //	Password string
    //}

	lg, err := localLog.New("logs", 10, 10, 10)
	if err != nil {
	panic(err)
	}
	
    sMgr,err:=SPR_go.New(SPR_go.RedisConfig{
        Addr:     "127.0.0.1",
        Port:     6379,
    },lg)
    if err != nil {
        log.Println(err)
        return
    }

    //add job with unique job name which used in redis
    //the process with same job name will scramble for the master token
    err = sMgr.AddJobName("testjob")
    if err != nil {
        log.Println(err)
    }
    err = sMgr.AddJobName("testjob2")
    if err != nil {
        log.Println(err)
    }

    // use function IsMaster("jobName") to check whether the process get the master token or not
    // if return true means get the master token
    go func() {
        for {
            time.Sleep(time.Second)
            log.Println("testjob is master:", sMgr.IsMaster("testjob"))
            log.Println("testjob2 is master:", sMgr.IsMaster("testjob2"))
        }
    }()

    // use function RemoveJobName("jobName") to remove the job
    // removed job always return false when use IsMaster("jobName")
    time.AfterFunc(time.Second*30, func() {
        sMgr.RemoveJobName("testjob2")
    })

    time.Sleep(1 * time.Hour)

}
```