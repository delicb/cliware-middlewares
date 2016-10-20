// Package headers contains middlewares for manipulating headers on request.
package headers

import (
	"net/http"

	c "github.com/delicb/cliware"
)

// Method sets request method to ongoing request.
func Method(method string) c.Middleware {
	return c.RequestProcessor(func(req *http.Request) error {
		req.Method = method
		return nil
	})
}

// Add adds provided header to ongoing request.
func Add(header, value string) c.Middleware {
	return c.RequestProcessor(func(req *http.Request) error {
		req.Header.Add(header, value)
		return nil
	})
}

// Set sets provided header to ongoing request.
func Set(header, value string) c.Middleware {
	return c.RequestProcessor(func(req *http.Request) error {
		req.Header.Set(header, value)
		return nil
	})
}

// Del removes provided header from ongoing request.
func Del(header string) c.Middleware {
	return c.RequestProcessor(func(req *http.Request) error {
		req.Header.Del(header)
		return nil
	})
}

// SetMap sets multiple headers provided in a map.
func SetMap(headers map[string]string) c.Middleware {
	return c.RequestProcessor(func(req *http.Request) error {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
		return nil
	})
}
