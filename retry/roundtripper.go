package retry

import (
	"context"
	"net/http"
	"time"
)

var (
	defaultBackoff      = ExponentialJitterBackoff(100*time.Millisecond, 30*time.Second, 2)
	defaultClassifier   = AnyErrorClassifier
	defaultMaxRetries   = 10
	defaultMaxDuration  = 3 * time.Minute
	defaultBodyStrategy = CacheBodyStrategy
	defaultRetryMethods = []string{"GET"}
)

// Enable modifiers provided client so that it can support all retry mechanisms
// implemented in this package. What it does it to replace client.Transport
// with RoundTripper that knows how to read and use values set via different
// middlewares. Original transport is still used for request sending, but it might
// be used multiple times (depending on retry logic).
func Enable(client *http.Client) {
	// if retryTransport is already set, no need to do it again
	if _, ok := client.Transport.(*retryTransport); ok {
		return
	}

	// set transport to retryTransport
	var origTransport http.RoundTripper
	if client.Transport != nil {
		origTransport = client.Transport
	} else {
		origTransport = http.DefaultTransport
	}
	client.Transport = NewRetryTransport(origTransport)
	return
}

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

type retryTransportConfig struct {
	Classifier   Classifier
	Backoff      BackoffStrategy
	MaxRetries   int
	MaxDuration  time.Duration
	BodyStrategy RetryBodyStrategy
	RetryMethods []string
}

func newRetryTransportConfig(ctx context.Context) *retryTransportConfig {
	config := &retryTransportConfig{
		Classifier:   getClassifier(ctx),
		Backoff:      getBackoff(ctx),
		MaxRetries:   getRetryTimes(ctx),
		MaxDuration:  getMaxDuration(ctx),
		BodyStrategy: getBodyStrategy(ctx),
		RetryMethods: getRetryMethods(ctx),
	}
	if config.Classifier == nil {
		config.Classifier = defaultClassifier
	}
	if config.Backoff == nil {
		config.Backoff = defaultBackoff
	}
	if config.MaxRetries == 0 {
		config.MaxRetries = defaultMaxRetries
	}
	if config.MaxDuration == time.Duration(0) {
		config.MaxDuration = defaultMaxDuration
	}
	if config.BodyStrategy == nil {
		config.BodyStrategy = defaultBodyStrategy
	}
	if config.RetryMethods == nil || len(config.RetryMethods) == 0 {
		config.RetryMethods = defaultRetryMethods
	}
	return config
}

func (t *retryTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	count := 0

	config := newRetryTransportConfig(r.Context())

	getBody, err := config.BodyStrategy(r)
	if err != nil {
		return nil, err
	}

	start := time.Now().UTC()
	for {
		// Copy request and sets its body to appropriate value
		reqCopy := &http.Request{}
		*reqCopy = *r
		reqCopy.Body = getBody()

		// perform actual request
		resp, err := t.next.RoundTrip(reqCopy)

		// check if we reached any of conditions for stopping retry cycle
		classifier := !config.Classifier(resp, err)
		maxRetries := count >= config.MaxRetries
		supportedMethod := !stringInSlice(r.Method, config.RetryMethods)

		currentDuration := time.Now().UTC().Sub(start)
		maxDuration := currentDuration.Nanoseconds() > config.MaxDuration.Nanoseconds()

		if classifier || maxRetries || supportedMethod || maxDuration {
			return resp, err
		}

		// if all else failed, increase number of retries and wait for some time
		count++
		time.Sleep(config.Backoff(count))
	}
}

func stringInSlice(s string, in []string) bool {
	for _, ss := range in {
		if s == ss {
			return true
		}
	}
	return false
}
