package middleware

import (
	"bytes"
	"github.com/dgraph-io/ristretto"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strconv"
	"time"
)

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

// setCacheHeader sets the client side caching to a third of the server side cache
func setCacheHeader(w http.ResponseWriter, duration time.Duration) {
	if duration == time.Duration(0) {
		duration = time.Hour * 24
	}
	w.Header().Set("Cache-Control", "max-age="+strconv.FormatInt(int64(duration/time.Second/3), 10))
}

// buildKey build a key from the given arguments
func buildKey(requestURI string, body []byte) string {
	if len(body) > 0 {
		requestURI += string(body)
	}

	return requestURI
}

// NewCache returns a middleware which caches all successful requests (status code < 400) for the given duration
func NewCache(ttl time.Duration, logger *slog.Logger) (Adapter, error) {
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10 M).
		MaxCost:     1 << 30, // maximum cost of cache (1 GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	if err != nil {
		return nil, err
	}

	return newCache(ttl, cache, logger), nil
}

// newCache caches all successful requests (status code < 400) for the given duration
func newCache(ttl time.Duration, cache *ristretto.Cache, logger *slog.Logger) Adapter {
	type cacheElement struct {
		buffer     []byte
		header     http.Header
		statusCode int
	}
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// extract body
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "an error occurred", http.StatusInternalServerError)
				if logger != nil {
					logger.Error(err.Error())
				}

				return
			}
			// reset body, so it can be read by the next handler
			r.Body = io.NopCloser(bytes.NewBuffer(body))
			cacheKey := buildKey(r.RequestURI, body)

			// try to get request from cache
			value, found := cache.Get(cacheKey)
			var response cacheElement
			if found {
				response = value.(cacheElement)
			} else {
				// record the writes of the next handler, so the response can be saved in the cache.
				recorder := httptest.NewRecorder()
				// call next handler
				h.ServeHTTP(recorder, r)

				// get recorded values
				resp := recorder.Result() //nolint:bodyclose
				defer func(Body io.ReadCloser) {
					_ = Body.Close()
				}(resp.Body)

				response.statusCode = resp.StatusCode
				response.buffer = recorder.Body.Bytes()
				response.header = resp.Header.Clone()

				// only insert in cache if no error occurred
				if response.statusCode < http.StatusBadRequest {
					cache.SetWithTTL(cacheKey, response, 1, ttl)
				}
			}

			// write original header
			for k := range response.header {
				for _, v := range response.header.Values(k) {
					w.Header().Add(k, v)
				}
			}
			// write cache header
			setCacheHeader(w, ttl)
			// write status code
			w.WriteHeader(response.statusCode)

			if _, err = w.Write(response.buffer); err != nil {
				if logger != nil {
					logger.Error(err.Error())
				}
			}
		})
	}
}