package errors

import "runtime/debug"

func WithMsg(err error, msg string) error {
	return &errWithMsg{
		error: err,
		msg:   msg,
	}
}

func WithStack(err error) error {
	return &errWithStack{
		error: err,
		stack: debug.Stack(),
	}
}
