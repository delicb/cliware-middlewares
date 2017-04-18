// Package auth contains middlewares for managing authentication during HTTP request.
package auth

import (
	"net/http"

	c "github.com/delicb/cliware"
)

// Basic sets basic authentication to request with provided username and password.
func Basic(username, password string) c.Middleware {
	return c.RequestProcessor(func(req *http.Request) error {
		req.SetBasicAuth(username, password)
		return nil
	})
}

// Bearer sets bearer authentication to request with provided token.
func Bearer(token string) c.Middleware {
	return c.RequestProcessor(func(req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+token)
		return nil
	})
}

// Custom sets a custom authorization header.
func Custom(authorization string) c.Middleware {
	return c.RequestProcessor(func(req *http.Request) error {
		req.Header.Set("Authorization", authorization)
		return nil
	})
}
