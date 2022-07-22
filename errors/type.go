package errors

import (
	"errors"
	"fmt"
)

type ErrMsg struct {
	Code int
	Msg  string
}

func (em *ErrMsg) Error() string {
	return fmt.Sprintf("code<%d>, msg<%s>", em.Code, em.Msg)
}

func (em *ErrMsg) Is(err error) bool {
	e, ok := err.(*ErrMsg)
	if ok {
		return e.Code == em.Code
	}
	return false
}

type errWithMsg struct {
	error
	msg string
}

func (em *errWithMsg) Error() string {
	return fmt.Sprintf("original<%s>, msg<%s>", em.error, em.msg)
}

func (em *errWithMsg) Is(err error) bool {
	return errors.Is(em.error, err)
}

type errWithStack struct {
	error
	stack []byte
}

func (em *errWithStack) Error() string {
	return fmt.Sprintf("original<%s>\n%s", em.error, em.stack)
}

func (em *errWithStack) Is(err error) bool {
	return errors.Is(em.error, err)
}
