// Package body contains middlewares for manipulating body of a request.
package body // import "go.delic.rs/cliware-middlewares/body"

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"encoding/json"
	"encoding/xml"

	"strings"

	"io"

	c "go.delic.rs/cliware"
)

// String sets request body to provided string.
func String(data string) c.Middleware {
	return c.RequestProcessor(func(req *http.Request) error {
		req.Method = getMethod(req)
		req.Body = ioutil.NopCloser(strings.NewReader(data))
		req.ContentLength = int64(bytes.NewBufferString(data).Len())
		return nil
	})
}

// JSON sets request body to JSON obtained from provided data.
// string and byte slice will be passed as is. For anything else, JSON
// encoding will be used. Content-Type header will be set to application/json.
func JSON(data interface{}) c.Middleware {
	return c.RequestProcessor(func(req *http.Request) error {
		buff := &bytes.Buffer{}
		switch data.(type) {
		case string:
			buff.WriteString(data.(string))
		case []byte:
			buff.Write(data.([]byte))
		default:
			if err := json.NewEncoder(buff).Encode(data); err != nil {
				return err
			}
		}

		req.Method = getMethod(req)
		req.Body = ioutil.NopCloser(buff)
		req.ContentLength = int64(buff.Len())
		req.Header.Set("Content-Type", "application/json")
		return nil
	})
}

// XML sets request body to XML obtained from provided data.
// string and byte slice will be passed as is. For anything else, XML
// encoding will be used. Content-Type header will be set to application/xml.
func XML(data interface{}) c.Middleware {
	return c.RequestProcessor(func(req *http.Request) error {
		buff := &bytes.Buffer{}
		switch data.(type) {
		case string:
			buff.WriteString(data.(string))

		case []byte:
			buff.Write(data.([]byte))
		default:
			if err := xml.NewEncoder(buff).Encode(data); err != nil {
				return err
			}
		}
		req.Method = getMethod(req)
		req.Body = ioutil.NopCloser(buff)
		req.ContentLength = int64(buff.Len())
		req.Header.Set("Content-Type", "application/xml")
		return nil
	})
}

// Reader sets request body to contain content from provided reader.
// Content type header is not set by this middleware.
func Reader(body io.Reader) c.Middleware {
	return c.RequestProcessor(func(req *http.Request) error {
		rc, ok := body.(io.ReadCloser)
		if !ok && body != nil {
			rc = ioutil.NopCloser(body)
		}

		if body != nil {
			switch v := body.(type) {
			case *bytes.Buffer:
				req.ContentLength = int64(v.Len())
			case *bytes.Reader:
				req.ContentLength = int64(v.Len())
			case *strings.Reader:
				req.ContentLength = int64(v.Len())
			}
		}
		req.Body = rc
		req.Method = getMethod(req)
		return nil
	})
}

func getMethod(req *http.Request) string {
	method := req.Method
	if method == "GET" || method == "" {
		return "POST"
	}
	return method
}
