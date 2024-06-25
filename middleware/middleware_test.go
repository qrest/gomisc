package middleware

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAdapt(t *testing.T) {
	mux := http.NewServeMux()
	mux.Handle("GET /someRoute", Adapt(http.NotFoundHandler(), MaxBody5MiB()))

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/someRoute", nil)
	require.NoError(t, err)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	mux.ServeHTTP(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)
}

func TestMaxBody(t *testing.T) {
	testHandler := func(shouldFail bool) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := io.ReadAll(r.Body)
			if shouldFail {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}

	mux := http.NewServeMux()
	mux.Handle("GET /someRoute", Adapt(testHandler(true), MaxBody(1)))
	mux.Handle("GET /someRoute2", Adapt(testHandler(false), MaxBody(3)))

	// request with a 2 MiB body, which should trigger an error for /someRoute
	req, err := http.NewRequest("GET", "/someRoute", bytes.NewReader(bytes.Repeat([]byte("1"), 2*1024*1024)))
	require.NoError(t, err)

	// request with a 2 MiB body, which should not trigger an error for /someRoute2
	req2, err := http.NewRequest("GET", "/someRoute2", bytes.NewReader(bytes.Repeat([]byte("1"), 2*1024*1024)))
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	mux.ServeHTTP(recorder, req)
	mux.ServeHTTP(recorder, req2)
}

func TestMaxBody5MiB(t *testing.T) {
	testHandler := func(shouldFail bool) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := io.ReadAll(r.Body)
			if shouldFail {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}

	mux := http.NewServeMux()
	mux.Handle("GET /someRoute", Adapt(testHandler(true), MaxBody5MiB()))
	mux.Handle("GET /someRoute2", Adapt(testHandler(false), MaxBody5MiB()))

	// request with a 6 MiB body, which should trigger an error for /someRoute
	req, err := http.NewRequest("GET", "/someRoute", bytes.NewReader(bytes.Repeat([]byte("1"), 6*1024*1024)))
	require.NoError(t, err)

	// request with a 4 MiB body, which should not trigger an error for /someRoute2
	req2, err := http.NewRequest("GET", "/someRoute2", bytes.NewReader(bytes.Repeat([]byte("1"), 4*1024*1024)))
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	mux.ServeHTTP(recorder, req)
	mux.ServeHTTP(recorder, req2)
}
