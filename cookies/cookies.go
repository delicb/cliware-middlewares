// Package cookies contains middlewares for manipulating cookies on request.
package cookies // import "go.delic.rs/cliware-middlewares/cookies"

import (
	"net/http"

	c "go.delic.rs/cliware"
)

// Add adds provided cookie to request.
func Add(cookie *http.Cookie) c.Middleware {
	return c.RequestProcessor(func(req *http.Request) error {
		req.AddCookie(cookie)
		return nil
	})
}

// Set sets new cookie for current request.
func Set(key, value string) c.Middleware {
	return c.RequestProcessor(func(req *http.Request) error {
		cookie := &http.Cookie{Name: key, Value: value}
		req.AddCookie(cookie)
		return nil
	})
}

// DelAll removes all cookies from request.
func DelAll() c.Middleware {
	return c.RequestProcessor(func(req *http.Request) error {
		req.Header.Del("Cookie")
		return nil
	})
}

// SetMap adds all cookies defined in provided map to request.
func SetMap(cookies map[string]string) c.Middleware {
	return c.RequestProcessor(func(req *http.Request) error {
		for k, v := range cookies {
			cookie := &http.Cookie{Name: k, Value: v}
			req.AddCookie(cookie)
		}
		return nil
	})
}

// AddMultiple adds all provided cookies to request.
func AddMultiple(cookies []*http.Cookie) c.Middleware {
	return c.RequestProcessor(func(req *http.Request) error {
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		return nil
	})
}
