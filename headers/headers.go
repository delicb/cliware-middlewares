// Package headers contains middlewares for manipulating headers on request.
package headers

import (
	"context"
	"net/http"

	"errors"

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

// Header holds information about what headers to add from context for FromContext middleware.
type Header struct {
	Key   string
	Value []string
}

// FromContext adds header to request that is defined in context with provided key.
func FromContext(key interface{}) c.Middleware {
	return c.MiddlewareFunc(func(next c.Handler) c.Handler {
		return c.HandlerFunc(func(ctx context.Context, req *http.Request) (resp *http.Response, err error) {
			value := ctx.Value(key)
			switch header := value.(type) {
			case Header:
				for _, v := range header.Value {
					req.Header.Set(header.Key, v)
				}
			case []Header:
				for _, hh := range header {
					for _, v := range hh.Value {
						req.Header.Set(hh.Key, v)
					}
				}
			default:
				return nil, errors.New("headers.FromContext middleware: value in unsupported format")
			}
			return next.Handle(ctx, req)
		})
	})
}

// ToContext adds headers to context that can be used with FromContext middleware.
// This is intended to be used for single header (but possibly multiple values).
// If you need to set multiple headers, use ToContextList.
func ToContext(ctx context.Context, key interface{}, header string, values ...string) context.Context {
	return context.WithValue(ctx, key, Header{
		Key:   header,
		Value: values,
	})
}

// ToContextList adds list of headers to context that can be used with FromContext
// middleware.
func ToContextList(ctx context.Context, key interface{}, h []Header) context.Context {
	return context.WithValue(ctx, key, h)
}
