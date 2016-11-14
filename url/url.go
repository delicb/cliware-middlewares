// Package url contains middlewares for manipulating request endpoint.
package url // import "go.delic.rs/cliware-middlewares/url"

import (
	"net/http"
	"net/url"
	"regexp"
	"strings"

	c "go.delic.rs/cliware"
)

// URL parses and sets entire URL to the request.
// If only path or query parameters needs to be changes, use other middlewares
// from this package.
func URL(rawURL string) c.Middleware {
	return c.RequestProcessor(func(req *http.Request) error {
		u, err := url.Parse(normalize(rawURL))
		if err != nil {
			return err
		}
		req.URL = u
		return nil
	})
}

// BaseURL parses and sets schema and host to the request.
func BaseURL(uri string) c.Middleware {
	return c.RequestProcessor(func(req *http.Request) error {
		u, err := url.Parse(normalize(uri))
		if err != nil {
			return err
		}
		req.URL.Scheme = u.Scheme
		req.URL.Host = u.Host
		return nil
	})
}

// Path sets path part of URL on the request.
func Path(path string) c.Middleware {
	return c.RequestProcessor(func(req *http.Request) error {
		req.URL.Path = normalizePath(path)
		return nil
	})
}

// AddPath appends provided path to currently existing path on a request.
func AddPath(path string) c.Middleware {
	return c.RequestProcessor(func(req *http.Request) error {
		req.URL.Path += normalizePath(path)
		return nil
	})
}

// PathPrefix adds provided path segment in front of current path of the request.
func PathPrefix(path string) c.Middleware {
	return c.RequestProcessor(func(req *http.Request) error {
		req.URL.Path = normalizePath(path) + req.URL.Path
		return nil
	})
}

// Param replaces one or multiple URL parameters with given value.
func Param(key, value string) c.Middleware {
	return c.RequestProcessor(func(req *http.Request) error {
		req.URL.Path = replace(req.URL.Path, key, value)
		return nil
	})
}

// Params replaces all provided parameters in URL with mapped values.
func Params(params map[string]string) c.Middleware {
	return c.RequestProcessor(func(req *http.Request) error {
		for k, v := range params {
			req.URL.Path = replace(req.URL.Path, k, v)
		}
		return nil
	})
}

func replace(str, key, value string) string {
	return strings.Replace(str, ":"+key, value, -1)
}

func normalizePath(path string) string {
	if path == "/" {
		return ""
	}
	return path
}

func normalize(uri string) string {
	match, _ := regexp.MatchString("^http[s]?://", uri)
	if match {
		return uri
	}
	return "http://" + uri
}
