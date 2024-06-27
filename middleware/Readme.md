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
	cache, err := mw.NewCache(time.Minute*5, nil)
	if err != nil {
		// handle error
	}

	// use the same cache for multiple routes
	mux.Handle("/somePath", mw.Adapt(someHandler(), cache, mw.MaxBody5MiB()))
	mux.Handle("/somePath2", mw.Adapt(someHandler(), cache, mw.MaxBody5MiB()))
}
```
