package serror

import (
	"errors"
	"fmt"
	"log/slog"
	"runtime/debug"
)

type StackError struct {
	err   error
	stack []byte
}

func NewStackError(err error) error {
	return StackError{
		err:   err,
		stack: debug.Stack(),
	}
}

func NewStackErrorStr(errorString string) error {
	return StackError{
		err:   errors.New(errorString),
		stack: debug.Stack(),
	}
}

func NewStackErrorf(format string, v ...any) error {
	return NewStackError(fmt.Errorf(format, v...))
}

func (se StackError) Error() string {
	if se.err == nil {
		return ""
	}
	return se.err.Error()
}

func (se StackError) Unwrap() error {
	return se.err
}

func (se StackError) Stack() []byte {
	return se.stack
}

func GetStack(err error) ([]byte, bool) {
	var s interface {
		Stack() []byte
	}
	if errors.As(err, &s) {
		return s.Stack(), true
	}
	return nil, false
}

func LogError(l *slog.Logger, err error, v ...any) {
	var st StackError
	if errors.As(err, &st) {
		v = append(v, slog.String("stack", string(st.Stack())))
		l.Warn(err.Error(), v...)
	} else {
		l.Warn(err.Error(), v...)
	}
}
