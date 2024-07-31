package serror

import (
	"errors"
	"fmt"
	"log/slog"
	"runtime/debug"
	"slices"
	"strings"
)

type StackError struct {
	err     error
	stack   []byte
	context []any
}

// getReducedStackTrace returns debug.Stack() with irrelevant bits removed
func getReducedStackTrace() []byte {
	var stackLines = strings.Split(string(debug.Stack()), "\n")

	cutStart := 0
	for i := len(stackLines) - 1; i >= 0; i-- {
		if strings.Contains(stackLines[i], "serror/error.go") {
			cutStart = i + 1
			break
		}
	}

	// cutoff everything before the last mention of error.go
	stackLines = stackLines[cutStart:]
	// remove empty lines
	stackLines = slices.DeleteFunc(stackLines, func(d string) bool { return strings.TrimSpace(d) == "" })

	return []byte(strings.Join(stackLines, "\n"))
}

func NewWithContext(err error, v ...any) error {
	return StackError{
		err:     err,
		stack:   getReducedStackTrace(),
		context: v,
	}
}

func New(err error) error {
	return StackError{
		err:   err,
		stack: getReducedStackTrace(),
	}
}

func FromStrWithContext(errorString string, v ...any) error {
	return StackError{
		err:     errors.New(errorString),
		stack:   getReducedStackTrace(),
		context: v,
	}
}

func FromStr(errorString string) error {
	return StackError{
		err:   errors.New(errorString),
		stack: getReducedStackTrace(),
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

func AddContext(err error, v ...any) error {
	var st StackError
	if errors.As(err, &st) {
		st.context = append(st.context, v...)
		return st
	}

	return err
}

func GetStack(err error) ([]byte, bool) {
	var st StackError
	if errors.As(err, &st) {
		return st.Stack(), true
	}
	return nil, false
}

func Log(l *slog.Logger, err error, v ...any) {
	var st StackError
	if errors.As(err, &st) {
		l.Warn(err.Error(), slices.Concat(v, st.context, []any{"stack", string(st.Stack())})...)
	} else {
		l.Warn(err.Error(), v...)
	}
}
