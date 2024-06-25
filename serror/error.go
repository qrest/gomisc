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

func New(err error) error {
	return StackError{
		err:   err,
		stack: debug.Stack(),
	}
}

func FromStr(errorString string) error {
	return StackError{
		err:   errors.New(errorString),
		stack: debug.Stack(),
	}
}

func FromFormat(format string, v ...any) error {
	return New(fmt.Errorf(format, v...))
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

func Log(l *slog.Logger, err error, v ...any) {
	var st StackError
	if errors.As(err, &st) {
		v = append(v, slog.String("stack", string(st.Stack())))
		l.Warn(err.Error(), v...)
	} else {
		l.Warn(err.Error(), v...)
	}
}
