package container

import (
	"fmt"
	pb "github.com/1851616111/xchain/pkg/protos"
	"github.com/fsouza/go-dockerclient"
	"log"
	"os"
	"sync"
	"time"
)

//目前只使用一个controller单goroutine 处理container的所有操作。
//TODO：当业务量增加时，添加多个worker同时分发处理

var (
	ctlOnce *sync.Once = new(sync.Once)
	ctl     *Controller
	logger  = log.New(os.Stdout, "[xcode]", log.LstdFlags)
)

func GetController() *Controller {
	if ctl == nil {
		ctlOnce.Do(
			func() {
				ctl = new(Controller)
				ctl.ping = time.Second * 60
				ctl.workerNum = 1
				ctl.deployTimeout = time.Second * 600
				ctl.jobCh = make(chan Job, 200)
				ctl.maxRefreshTimes = 10

			})
	}

	cli, err := getDefaultClient()
	if err != nil {
		logger.Fatalf("xcode create controller err %v\n", err)
		return nil
	}
	ctl.c = container{cli}

	return ctl
}

type Controller struct {
	c               container
	ping            time.Duration
	maxRefreshTimes int
	refreshTimes    int
	maxWorkCh       int

	workerNum     int
	deployTimeout time.Duration
	jobCh         chan Job
}

func (c *Controller) Start() {

	logger.Printf("xcode controller: %v\n", *c)
	go func() {
		for {
			select {
			case <-time.Tick(c.ping):

				if err := c.c.client.Ping(); err != nil {
					logger.Println("container connection break")

					if err := c.c.RefreshClient(); err != nil {
						c.refreshTimes++
						//1. 当work队列中的work数量大于等于maxworksOnErr的数量时或者
						//2. 当出错时refresh的数量大于等于maxRefreshTimes数量时.程序应该停止工作.
						//TODO：程序停止工作前，将workCh队列的work缓存的系统的文件中，启动时再加载？？？
						if c.refreshTimes >= c.maxRefreshTimes {
							logger.Fatalf("controller refreshTimes times(%d) larger than maxRefreshTimes(%d), progrem exit.", c.refreshTimes, c.maxRefreshTimes)
							os.Exit(1)
						}

						if len(c.jobCh) >= c.maxWorkCh {
							logger.Fatalf("controller work channel length larger than maxworksOnErr(%d), progrem exit.", c.maxWorkCh)
							os.Exit(1)
						}
						continue
					}
					logger.Println("refresh container connection success")
				}

			case w := <-c.jobCh:
				go func() {
					w.Report(w.Do())
				}()
			}
		}
	}()
}

func (c *Controller) Dispatch(work *Worker) error {
	if work == nil {
		return ErrWorkerNil
	}

	if err := work.Validate(); err != nil {
		return err
	}

	c.jobCh <- work
	return nil
}

func (c *Controller) Deploy(spec *pb.XCodeSpec) error {

	result := make(chan interface{}, 10)
	work := &Worker{
		act:      Job_Action_BuildImage,
		id:       spec.XcodeID.Path,
		lang:     spec.Type,
		metadata: spec,

		opts: &docker.BuildImageOptions{
			Name: fmt.Sprintf("xcode-%s-%s", spec.Type.String(), spec.XcodeID.Path),
		},
		resultCh: result,
	}
	if err := c.Dispatch(work); err != nil {
		logger.Printf("dispatch deploy work(%v) err:%v\n", *work, err)
		return err
	}

	for {
		select {
		case res, ok := <-result:
			if !ok { // res == nil, close(result) success
				return nil
			}
			if err, ok := res.(error); ok {
				return err
			}
		case <-time.Tick(c.deployTimeout):
			return ErrJobDeployTimeout
		}
	}
}
