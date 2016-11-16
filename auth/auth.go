package auth

import (
	"net/http"

	c "go.delic.rs/cliware"
)

// Basic sets basic authentication to request with provided username and apssword.
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
