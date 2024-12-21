package middleware

import (
	"bytes"
	"github.com/dgraph-io/ristretto/v2"
	"github.com/qrest/gomisc/serror"
	"io"
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

// instances of this type are saved in the cache middleware
type responseItem struct {
	buffer     []byte
	header     http.Header
	statusCode int
}

// buildKey build a key from the given arguments
func buildKey(requestURI string, body []byte) string {
	if len(body) > 0 {
		requestURI += string(body)
	}

	return requestURI
}

// NewCacheFactory returns a middleware factory. Call the factory for each route, so they share the same internal cache.
func NewCacheFactory(numMiB int64, errorCallback func(error)) (func(ttl time.Duration) Adapter, error) {
	cache, err := ristretto.NewCache[string, responseItem](&ristretto.Config[string, responseItem]{
		NumCounters: 1e7, // number of keys to track frequency of (10 M).
		MaxCost:     1024 * 1024 * numMiB,
		BufferItems: 64, // number of keys per Get buffer.
	})
	if err != nil {
		return nil, serror.New(err)
	}

	return func(ttl time.Duration) Adapter { return newCache(ttl, cache, errorCallback) }, nil
}

// newCache caches all successful requests (status code < 400) for the given duration.
// If errorCallback is set, it is called when an error occurs
func newCache(ttl time.Duration, cache *ristretto.Cache[string, responseItem], errorCallback func(error)) Adapter {

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// extract body
			var body []byte
			if r.Body != nil {
				var err error
				body, err = io.ReadAll(r.Body)
				if err != nil {
					http.Error(w, "an error occurred", http.StatusInternalServerError)
					if errorCallback != nil {
						errorCallback(err)
					}

					return
				}
				// reset body, so it can be read by the next handler
				r.Body = io.NopCloser(bytes.NewBuffer(body))
			}

			cacheKey := buildKey(r.URL.Path+"|"+r.URL.RawQuery, body)

			// try to get request from cache
			response, found := cache.Get(cacheKey)
			if !found {
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
					cache.SetWithTTL(cacheKey, response, int64(40+4*len(response.buffer)), ttl)
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

			if _, err := w.Write(response.buffer); err != nil {
				if errorCallback != nil {
					errorCallback(err)
				}
			}
		})
	}
}
