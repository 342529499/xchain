package container

import (
	"fmt"
	pb "github.com/1851616111/xchain/pkg/protos"
	"strings"
)

type action int

var (
	m map[int]string = map[int]string{
		Job_Action_BuildImage:      "build image",
		Job_Action_CreateContainer: "create container",
		Job_Action_RemoveContainer: "remove container",
		Job_Action_RemoveImage:     "remove image",
	}
)

const (
	Job_Action_BuildImage = iota + 1
	Job_Action_CreateContainer
	Job_Action_RemoveContainer
	Job_Action_RemoveImage
)

type Job interface {
	Do() error
	Report(result interface{})

	Action() string
	Details() string
	Language() string
}

type Worker struct {
	act  action
	id   string
	lang pb.XCodeSpec_Type

	metadata interface{}
	opts     interface{}
	resultCh chan interface{}
}

func (w *Worker) Validate() error {
	if w.act < Job_Action_BuildImage || w.act > Job_Action_RemoveImage {
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

func (w *Worker) Do() error {
	switch w.act {
	case Job_Action_BuildImage:
		return w.buildImage()
	case Job_Action_CreateContainer:
		return w.createContainer()
	case Job_Action_RemoveContainer:
		return w.removeContainer()
	case Job_Action_RemoveImage:
		return w.removeImage()
	default:
		return ErrUnknownJobType
	}
}

func (w *Worker) Report(result interface{}) {
	if result == nil {
		close(w.resultCh)

		logger.Printf("%s job ok, details:%s\n", w.Action(), w.Details())
		return
	}

	logger.Printf("%s job error:%v, details:%s\n", w.Action(), result, w.Details())
	w.resultCh <- result
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
