package serror

import (
	"errors"
	"github.com/stretchr/testify/require"
	"log/slog"
	"testing"
)

func TestGetStack(t *testing.T) {
	tests := []struct {
		err      error
		nilBytes bool
		flag     bool
	}{
		{
			err:      errors.New("some string"),
			nilBytes: true,
			flag:     false,
		},
		{
			err:      New(errors.New("some string")),
			nilBytes: false,
			flag:     true,
		},
	}
	for _, tt := range tests {
		stack, b := GetStack(tt.err)
		require.Equal(t, tt.nilBytes, stack == nil)
		require.Equal(t, tt.flag, b)
	}
}

func TestNewStackError(t *testing.T) {
	require.Error(t, New(errors.New("some string")))
}

func TestStackError_Error(t *testing.T) {
	tests := []struct {
		err  error
		want string
	}{
		{
			err:  New(errors.New("some string")),
			want: "some string",
		},
		{
			err:  New(errors.New("")),
			want: "",
		},
		{
			err:  New(nil),
			want: "",
		},
	}
	for _, tt := range tests {
		require.Equal(t, tt.want, tt.err.Error())
	}
}

func TestStackError_Stack(t *testing.T) {
	tests := []struct {
		err error
	}{
		{
			err: New(errors.New("some string")),
		},
		{
			err: New(errors.New("")),
		},
		{
			err: New(nil),
		},
	}
	for _, tt := range tests {
		var st StackError
		if !errors.As(tt.err, &st) {
			require.Fail(t, "error has wrong type")
		}
		require.NotNil(t, st.Stack())
	}
}

func TestStackError_Unwrap(t *testing.T) {
	tests := []struct {
		err     error
		wantErr bool
	}{
		{
			err:     New(errors.New("some string")),
			wantErr: true,
		},
		{
			err:     New(errors.New("")),
			wantErr: true,
		},
		{
			err:     New(nil),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		var st StackError
		if !errors.As(tt.err, &st) {
			require.Fail(t, "error has wrong type")
		}

		require.Equal(t, tt.wantErr, st.Unwrap() != nil, tt.err)
	}
}

func TestNewStackErrorStr(t *testing.T) {
	require.Error(t, FromStr("some string"))
}

func TestNewStackErrorf(t *testing.T) {
	require.Error(t, FromFormat("new error %d", 1))
}

func TestLog(t *testing.T) {
	require.NotPanics(t, func() {
		err := FromStr("test error")
		Log(slog.Default(), err)
		Log(slog.Default(), AddContext(err, "testkey", "testval"))
		Log(slog.Default(), FromStrWithContext("test string", "testkey", "testval"))
		Log(slog.Default(), FromStrWithContext("test string",
			"testkey1", "testval1", "testkey2", "testval2"))
		Log(slog.Default(), errors.New("some string"))
	})
}

func TestNewWithContext(t *testing.T) {
	require.Error(t, NewWithContext(errors.New("some string"), "key", "val"))
}

func TestFromStrWithContext(t *testing.T) {
	require.Error(t, FromStrWithContext("some string", "key", "val"))
}
