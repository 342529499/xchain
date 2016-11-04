package broker

import "github.com/pkg/errors"

var (
	ERRInstructionSourceUnauthorized = errors.New("instruction source unauthorized.")
	ERRPointerObjectNil              = errors.New("pointer object nil.")
	ERRWaitReturnInstructionTimeOut  = errors.New("wait return instruction timeout.")
	ERRUnExpectedResponsePayload     = errors.New("unexpected return instruction payload")
)
