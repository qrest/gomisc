package middleware

import "net/http"

type Adapter func(http.Handler) http.Handler

// Adapt constructs a sequence of http.Handler which act as a middleware.
func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
	for i := len(adapters) - 1; i >= 0; i-- {
		h = adapters[i](h)
	}
	return h
}

// MaxBody limits the amount of bytes which can be read from the request body to size (in number of MiB).
func MaxBody(size int64) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, 1024*1024*size)
			h.ServeHTTP(w, r)
		})
	}
}

// MaxBody5MiB limits the amount of bytes which can be read from the request body to 5 MiB.
func MaxBody5MiB() Adapter {
	return MaxBody(5)
}
