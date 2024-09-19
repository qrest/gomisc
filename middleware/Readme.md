# Middleware

This module supports adding middlewares to a route adapter chain.

## Using this module

Use `Adapt()` to construct a chain of middlewares.

## Example

```go
package main

import (
	mw "github.com/qrest/gomisc/middleware"
	"net/http"
	"time"
)

func someHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// do stuff
	})
}

func main() {
	mux := http.NewServeMux()
	// create a cache with a size of 50 MiB 
	cacheFactory, err := mw.NewCacheFactory(50, nil)
	if err != nil {
		// handle error
	}

	// call the cache factory to generate a new middleware for each request, so it shares its internal cache
	mux.Handle("/somePath", mw.Adapt(someHandler(), cacheFactory(time.Minute*5), mw.MaxBody5MiB()))
	mux.Handle("/somePath2", mw.Adapt(someHandler(), cacheFactory(time.Minute*10), mw.MaxBody5MiB()))
}
```
