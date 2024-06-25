package network

import "net/http"

type Adapter func(http.Handler) http.Handler

func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
	for i := len(adapters) - 1; i >= 0; i-- {
		h = adapters[i](h)
	}
	return h
}

// maxBodyConfig limits the amount of bytes which can be read from the request body to size (in number of MiB).
func maxBodyConfig(size int64) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, 1024*1024*size)
			h.ServeHTTP(w, r)
		})
	}
}
