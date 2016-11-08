package container

import (
	"fmt"
	pb "github.com/1851616111/xchain/pkg/protos"
	"strings"
)

type action int

var (
	m map[int]string = map[int]string{
		Job_Action_ListImage:       "list images",
		Job_Action_BuildImage:      "build image",
		Job_Action_ListContainer:   "list containers",
		Job_Action_CreateContainer: "create container",
		Job_Action_RemoveContainer: "remove container",
		Job_Action_RemoveImage:     "remove image",
	}
)

const (
	Job_Action_ListImage = iota + 1
	Job_Action_BuildImage
	Job_Action_ListContainer
	Job_Action_CreateContainer
	Job_Action_RemoveContainer
	Job_Action_RemoveImage
)

type Job interface {
	Do() (interface{}, error)
	Report(interface{}, error)

	Action() string
	Details() string
	Language() string
}

type Worker struct {
	preWork *Worker

	act  action
	id   string
	lang pb.XCodeSpec_Type

	metadata interface{}
	opts     interface{}

	resultCh chan interface{}
	errCh    chan error
}

func (w *Worker) Validate() error {
	if w.act < Job_Action_ListImage || w.act > Job_Action_RemoveImage {
		return ErrWorkerActionNotAllow
	}

	if len(strings.TrimSpace(w.id)) == 0 {
		return ErrWorkerIDNotFound
	}

	if w.lang < pb.XCodeSpec_GOLANG || w.lang > pb.XCodeSpec_JAVA {
		return ErrWorkerLanguageNotAllowed
	}

	if w.metadata == nil {
		return ErrWorkerMetadataNotFound
	}

	if w.opts == nil {
		return ErrWorkerOptionsNotFound
	}

	return nil
}

func (w *Worker) preDo() (interface{}, error) {
	if w.preWork != nil {
		return w.preWork.Do()
	}

	return nil, nil
}

func (w *Worker) Do() (interface{}, error) {
	if res, err := w.preDo(); err != nil {
		w.preWork.Report(res, err)
		return nil, err
	}

	switch w.act {
	case Job_Action_ListImage:
		return w.listImage()
	case Job_Action_BuildImage:
		return nil, w.buildImage()
	case Job_Action_ListContainer:
		return w.listContainers()
	case Job_Action_CreateContainer:
		return nil, w.createContainer()
	case Job_Action_RemoveContainer:
		return nil, w.removeContainer()
	case Job_Action_RemoveImage:
		return nil, w.removeImage()
	default:
		return nil, ErrUnknownJobType
	}
}

func (w *Worker) Report(result interface{}, err error) {
	if err == nil {
		if w.resultCh != nil {
			w.resultCh <- result
		}
		close(w.errCh)
		logger.Printf("do job(%s) ok, details:%s\n", w.Action(), w.Details())
	} else {
		w.errCh <- err
		if w.resultCh != nil {
			close(w.resultCh)
		}
		logger.Printf("do job(%s) error:%v, details:%s\n", w.Action(), result, w.Details())
	}
}

func (w *Worker) Action() string {
	act, ok := m[int(w.act)]
	if !ok {
		return "unknown action of job"
	}

	return act
}

func (w *Worker) Details() string {
	return fmt.Sprintf("id=%s", w.id)
}

func (w *Worker) Language() string {
	lang, ok := m[int(w.lang)]
	if !ok {
		return "unknown language of job"
	}

	return lang

}
