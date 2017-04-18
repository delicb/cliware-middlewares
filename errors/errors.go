// Package errors contains middlewares for converting HTTP response codes to GoLang errors.
package errors

import (
	"fmt"
	"net/http"

	"io/ioutil"

	c "github.com/delicb/cliware"
)

// HTTPError holds information about failed HTTP request.
type HTTPError struct {
	Name       string
	StatusCode int
	RequestURL string
	Method     string
	Body       []byte
}

// Error is implementation of error interface for HTTPError. It returns basic
// information about error that occurred (status code, requested URL)
func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTPError: %s - %s (%s)", e.Method, e.RequestURL, e.Name)
}

func createError(resp *http.Response) error {
	if resp.StatusCode < 400 {
		return nil
	}
	var rawData []byte
	if resp.Body != nil {
		rawData, _ = ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
	}

	return &HTTPError{
		Name:       resp.Status,
		StatusCode: resp.StatusCode,
		RequestURL: resp.Request.URL.String(),
		Method:     resp.Request.Method,
		Body:       rawData,
	}
}

// Errors convert HTTP status codes that represent errors to HTTPError.
func Errors() c.Middleware {
	return c.ResponseProcessor(func(resp *http.Response, err error) error {
		// if we already got error just send it down the chain
		if err != nil {
			return err
		}
		return createError(resp)
	})
}
