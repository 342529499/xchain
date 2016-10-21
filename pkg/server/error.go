package server

import "errors"

var (
	ExistErr   error = errors.New("exist")
	NoExistErr error = errors.New("no exist")
)
