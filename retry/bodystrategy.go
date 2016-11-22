package retry

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
)

// BodyStrategy defines what to do with request body in event when request
// has to be send again (because Classifier determined so). Strategy accepts
// request (this will be original request always) and returns function that
// should provide io.ReadCloser that will be set to every subsequent request
// body. This function will be called multiple times (as many times as request
// is retried).
type RetryBodyStrategy func(r *http.Request) (func() io.ReadCloser, error)

// CacheBodyStrategy caches initial request body in buffer and returns it
// every time it is needed.
func CacheBodyStrategy(r *http.Request) (func() io.ReadCloser, error) {
	var buf []byte
	var err error
	if r.Body != nil {
		buf, err = ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		r.Body.Close()
	}
	return func() io.ReadCloser {
		return ioutil.NopCloser(bytes.NewBuffer(buf))
	}, nil
}
