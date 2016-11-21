package retry

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	defaultBackoff    = ExponentialJitterBackoff(100*time.Millisecond, 30*time.Second, 2)
	defaultClassifier = AnyErrorClassifier
	defaultMaxRetries = 10
	defaultMaxTime    = 3 * time.Minute
)

// NewRetryTransport returns RoundTripper that wraps around provided RoundTripper
// and adds retry logic around it. All retry parameters are read from http.Request
// context (with sane defaults).
func NewRetryTransport(next http.RoundTripper) http.RoundTripper {
	return &retryTransport{
		next: next,
	}
}

type retryTransport struct {
	next http.RoundTripper
}

func (t *retryTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	count := 0

	// configure retry settings
	classifier := getClassifier(r.Context())
	backoff := getBackoff(r.Context())
	maxRetries := getRetryTimes(r.Context())
	maxDuration := geMaxDuration(r.Context())
	if classifier == nil {
		classifier = defaultClassifier
	}
	if backoff == nil {
		backoff = defaultBackoff
	}
	if maxRetries == 0 {
		maxRetries = defaultMaxRetries
	}
	if maxDuration == time.Duration(0) {
		maxDuration = defaultMaxTime
	}

	start := time.Now().UTC()

	// buffer request body for potentially repeated requests
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	r.Body.Close()

	for {
		// Copy request and sets its body to appropriate value
		reqCopy := &http.Request{}
		*reqCopy = *r
		reqCopy.Body = ioutil.NopCloser(bytes.NewBuffer(buf))

		// perform actual request
		resp, err := t.next.RoundTrip(reqCopy)

		// check if we reached any of conditions for stopping retry cycle
		currentDuration := time.Now().UTC().Sub(start)
		if !classifier(resp, err) || count > maxRetries || currentDuration.Nanoseconds() > maxDuration.Nanoseconds() {
			return resp, err
		}
		// if all else failed, increase number of retries and wait for some time
		count++
		time.Sleep(backoff(count))
	}
}
