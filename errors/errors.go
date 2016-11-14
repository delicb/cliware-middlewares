// Package errors contains middlewares for converting HTTP response codes to GoLang errors.
package errors // import "go.delic.rs/cliware-middlewares/errors"

import (
	"fmt"
	"net/http"

	"io/ioutil"

	c "go.delic.rs/cliware"
)

type HTTPError struct {
	Name       string
	StatusCode int
	Body       string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTP error: %s (%d)", e.Name, e.StatusCode)
}

func createError(resp *http.Response) error {
	if resp.StatusCode < 400 {
		return nil
	}
	defer resp.Body.Close()
	rawData, _ := ioutil.ReadAll(resp.Body)
	return &HTTPError{
		Name:       resp.Status,
		StatusCode: resp.StatusCode,
		Body:       string(rawData),
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
