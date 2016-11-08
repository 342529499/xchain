package container

import (
	"github.com/1851616111/go-dockerclient"
	pb "github.com/1851616111/xchain/pkg/protos"
	bm "github.com/1851616111/xchain/pkg/xcode/server/broker_manager"
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
				ctl.timeout = time.Second * 600
				ctl.jobCh = make(chan Job, 200)
				ctl.maxRefreshTimes = 10

				ctl.brokerNameToPortM = map[string]string{}

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

	workerNum int
	timeout   time.Duration
	jobCh     chan Job

	brokerNotifier    bm.Notifier
	brokerNameToPortM map[string]string
}

func (c *Controller) SetBrokerNotifier(nodeID, nodeAddress string) {
	c.brokerNotifier = bm.GetBrokerManager(nodeID, nodeAddress)
}

func (c *Controller) Run() {

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

func (c *Controller) PreDeploy(spec *pb.XCodeSpec) (err error) {
	resultCh, errCh := make(chan interface{}, 5), make(chan error, 5)

	id := genCodeID(spec)
	opt := &docker.ListImagesOptions{}
	convertLabelToFilter(&opt.Filters, map[string]string{
		"language": spec.Type.String(),
		"code":     spec.XcodeID.Path,
		"id":       id,
	})

	work := &Worker{
		act:      Job_Action_ListImage,
		id:       spec.XcodeID.Path,
		lang:     spec.Type,
		metadata: spec,

		opts:     opt,
		resultCh: resultCh,
		errCh:    errCh,
	}
	if err = c.Dispatch(work); err != nil {
		logger.Printf("dispatch prepare deploy work(%#v) err:%v\n", *work, err)
		return
	}

	var ok bool
	for {
		select {
		case err, ok = <-errCh:
			if ok {
				return
			} else if res, ok2 := <-resultCh; ok2 {
				if images, ok3 := res.([]docker.APIImages); ok3 {
					if len(images) > 0 {
						err = ErrDeployImageExists
					}
					return
				}
			}
		case <-time.Tick(c.timeout):
			err = ErrJobDeployTimeout
			return
		}
	}

}

func (c *Controller) Deploy(spec *pb.XCodeSpec) (err error) {

	errCh := make(chan error, 5)

	id := genCodeID(spec)
	work := &Worker{
		act:      Job_Action_BuildImage,
		id:       spec.XcodeID.Path,
		lang:     spec.Type,
		metadata: spec,

		opts: &docker.BuildImageOptions{
			Name: id,
			Labels: map[string]string{
				"language": spec.Type.String(),
				"code":     spec.XcodeID.Path,
				"id":       id,
			},
		},
		errCh: errCh,
	}
	if err = c.Dispatch(work); err != nil {
		logger.Printf("dispatch deploy work(%v) err:%v\n", *work, err)
		return
	}

	for {
		select {
		case e, ok := <-errCh:
			if !ok {
				err = nil
				return
			} else {
				err = e
				return
			}

		case <-time.Tick(c.timeout):
			err = ErrJobDeployTimeout
			return
		}
	}
}

func (c *Controller) PreStart(spec *pb.XCodeSpec) (err error) {
	resultCh, errCh := make(chan interface{}, 5), make(chan error, 5)

	opt := &docker.ListContainersOptions{}
	convertLabelToFilter(&opt.Filters, map[string]string{
		"language": spec.Type.String(),
		"code":     spec.XcodeID.Path,
		"id":       genDockerID(genCodeID(spec)),
	})

	work := &Worker{
		act:      Job_Action_ListContainer,
		id:       spec.XcodeID.Path,
		lang:     spec.Type,
		metadata: spec,

		opts:     opt,
		resultCh: resultCh,
		errCh:    errCh,
	}
	if err = c.Dispatch(work); err != nil {
		logger.Printf("dispatch prepare start(%#v) err:%v\n", *work, err)
		return
	}

	var ok bool
	for {
		select {
		case err, ok = <-errCh:
			if ok {
				return
			} else if res, ok2 := <-resultCh; ok2 {
				if containers, ok3 := res.([]docker.APIContainers); ok3 {
					if len(containers) > 0 {
						err = ErrDeployContainerExists
					}
					return
				}
			}
		case <-time.Tick(c.timeout):
			err = ErrJobDeployTimeout
			return
		}
	}


}

func (c *Controller) Start(spec *pb.XCodeSpec) (err error) {
	defer func() {
		if err == nil {

			//TODO:这里只使用10692. 以后需要根绝docker is 来变换10692
			c.brokerNameToPortM[spec.XcodeID.Name] = "10692"
			if c.brokerNotifier != nil {
				go func() {
					c.brokerNotifier.Notify(bm.Event{
						BrokerName: spec.XcodeID.Name,
						BrokerPort: "10692", //broker的监听地址
						Kind:       bm.EVENT_BROKER_START,
					})
				}()
			}
		}
	}()

	//TODO: 这个名字生成缺少deploy的参数部分
	id := genDockerID(genCodeID(spec))
	errCh := make(chan error, 5)
	work := &Worker{
		act:      Job_Action_CreateContainer,
		id:       spec.XcodeID.Path,
		lang:     spec.Type,
		metadata: spec,

		opts: &docker.CreateContainerOptions{
			Name:       id,
			Config:     &docker.Config{
				Image: genCodeID(spec),
				Labels: map[string]string{
					"language": spec.Type.String(),
					"code":     spec.XcodeID.Path,
					"id":       id,
				},
			},
			HostConfig: getDockerHostConfig(),
		},
		errCh: errCh,
	}

	if err = c.Dispatch(work); err != nil {
		logger.Printf("dispatch start container work(%#v) err:%v\n", *work, err)
		return
	}

	for {
		select {
		case e, ok := <-errCh:
			if !ok {
				err = nil
				return
			} else {
				err = e
				return
			}
		case <-time.Tick(c.timeout):
			err = ErrJobDeployTimeout
			return
		}
	}
}
