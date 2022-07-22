package errors

import (
	"fmt"
	"runtime/debug"
)

func WithMsg(err error, msg string) error {
	return &errWithMsg{
		error: err,
		msg:   msg,
	}
}

func WithMsgf(err error, format string, args ...any) error {
	return &errWithMsg{
		error: err,
		msg:   fmt.Sprintf(format, args...),
	}
}

func WithStack(err error) error {
	return &errWithStack{
		error: err,
		stack: debug.Stack(),
	}
}
