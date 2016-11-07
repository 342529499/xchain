package container

import (
	"errors"
	"fmt"
)

var (
	InterfaceAssertError = func(tp string) error {
		return fmt.Errorf("assert interface to %s error", tp)
	}

	ErrUnknownJobType           error = errors.New("job type unkown error")
	ErrDeployWorkDuplicated     error = errors.New("deploy work duplicated")
	ErrWorkerNil                error = errors.New("job worker nil error")
	ErrWorkerActionNotAllow     error = errors.New("job worker action not allowed")
	ErrWorkerIDNotFound         error = errors.New("job worker id not found")
	ErrWorkerLanguageNotAllowed error = errors.New("job worker id not allowed")
	ErrWorkerMetadataNotFound   error = errors.New("job worker metadata not found")
	ErrWorkerOptionsNotFound    error = errors.New("job worker options not found")

	ErrJobDeployTimeout error = errors.New("job deployed timeout")
)
