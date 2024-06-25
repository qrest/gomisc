package serror

import (
	"errors"
	"github.com/stretchr/testify/require"
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
			err:      NewStackError(errors.New("some string")),
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
	require.Error(t, NewStackError(errors.New("some string")))
}

func TestStackError_Error(t *testing.T) {
	tests := []struct {
		err  error
		want string
	}{
		{
			err:  NewStackError(errors.New("some string")),
			want: "some string",
		},
		{
			err:  NewStackError(errors.New("")),
			want: "",
		},
		{
			err:  NewStackError(nil),
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
			err: NewStackError(errors.New("some string")),
		},
		{
			err: NewStackError(errors.New("")),
		},
		{
			err: NewStackError(nil),
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
			err:     NewStackError(errors.New("some string")),
			wantErr: true,
		},
		{
			err:     NewStackError(errors.New("")),
			wantErr: true,
		},
		{
			err:     NewStackError(nil),
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
	require.Error(t, NewStackErrorStr("some string"))
}

func TestNewStackErrorf(t *testing.T) {
	require.Error(t, NewStackErrorf("new error %d", 1))
}
