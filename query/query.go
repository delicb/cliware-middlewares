// Package query contains middlewares for manipulating query string on request.
package query // import "go.delic.rs/cliware-middlewares/query"

import (
	"net/http"

	c "go.delic.rs/cliware"
)

// Set sets value as query parameter with provided key to URL.
func Set(key, value string) c.Middleware {
	return c.RequestProcessor(func(req *http.Request) error {
		query := req.URL.Query()
		query.Set(key, value)
		req.URL.RawQuery = query.Encode()
		return nil
	})
}

// Add adds query parameter to any existing query parameter or adds new
// if there are not existing query parameters to URL.
func Add(key, value string) c.Middleware {
	return c.RequestProcessor(func(req *http.Request) error {
		query := req.URL.Query()
		query.Add(key, value)
		req.URL.RawQuery = query.Encode()
		return nil
	})
}

// Del removes query parameter with provided key from URL.
func Del(key string) c.Middleware {
	return c.RequestProcessor(func(req *http.Request) error {
		query := req.URL.Query()
		query.Del(key)
		req.URL.RawQuery = query.Encode()
		return nil
	})
}

// DelAll removes all query parameters from URL.
func DelAll() c.Middleware {
	return c.RequestProcessor(func(req *http.Request) error {
		req.URL.RawQuery = ""
		return nil
	})
}

// SetMap adds all query parameters provided in map to URL.
func SetMap(queryMap map[string]string) c.Middleware {
	return c.RequestProcessor(func(req *http.Request) error {
		query := req.URL.Query()
		for k, v := range queryMap {
			query.Set(k, v)
		}
		req.URL.RawQuery = query.Encode()
		return nil
	})
}
