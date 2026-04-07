package rpc

import (
	"errors"
	"fmt"
)

var (
	ErrNotStarted      = errors.New("rpc client not started")
	ErrAlreadyStarted  = errors.New("rpc client already started")
	ErrClientClosed    = errors.New("rpc client closed")
	ErrProcessExited   = errors.New("rpc process exited")
	ErrResponseTimeout = errors.New("rpc response timeout")
	ErrProtocol        = errors.New("rpc protocol error")
)

type CommandError struct {
	Command string
	Message string
}

func (e *CommandError) Error() string {
	return fmt.Sprintf("%s failed: %s", e.Command, e.Message)
}
