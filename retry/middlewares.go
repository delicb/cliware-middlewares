package retry

import (
	"context"
	"net/http"

	"time"

	"fmt"

	c "go.delic.rs/cliware"
)

// Enable modifies transport for provided client to support retry mechanism.
// All this does is to set Transport to retryTransport defined here. Transport
// that is already set to provided client will be used to send actual request,
// so any config in it will be used.
func Enable(client *http.Client) c.Middleware {
	if client == nil {
		panic("EnableRetry: nil client")
	}
	return c.MiddlewareFunc(func(next c.Handler) c.Handler {
		return c.HandlerFunc(func(ctx context.Context, req *http.Request) (resp *http.Response, err error) {
			fmt.Println("In retry middleware")

			// if retryTransport is already set no need to do it again
			if _, ok := client.Transport.(*retryTransport); ok {
				return next.Handle(ctx, req)
			}

			// set transport to retryTransport
			var originalTransport http.RoundTripper
			if client.Transport != nil {
				originalTransport = client.Transport
			} else {
				originalTransport = http.DefaultTransport
			}
			client.Transport = NewRetryTransport(originalTransport)
			// send request
			return next.Handle(ctx, req)
		})
	})
}

// Times sets maximum number of times for HTTP request sending retry.
func Times(times int) c.Middleware {
	return c.ContextProcessor(func(ctx context.Context) context.Context {
		return setRetryTimes(ctx, times)
	})
}

// SetClassifier sets provided classifier to request.
func SetClassifier(classifier func(*http.Response, error) bool) c.Middleware {
	return c.ContextProcessor(func(ctx context.Context) context.Context {
		return setClassifier(ctx, Classifier(classifier))
	})
}

// SetBackoffStrategy sets backoff strategy for requests that classifier marked
// to be retried.
func SetBackoffStrategy(backoff func(n int) time.Duration) c.Middleware {
	return c.ContextProcessor(func(ctx context.Context) context.Context {
		return setBackoff(ctx, BackoffStrategy(backoff))
	})
}

// MaxDuration sets maximum amount of time to retry. This time includes everything,
// all requests, etc...
func MaxDuration(maxTime time.Duration) c.Middleware {
	return c.ContextProcessor(func(ctx context.Context) context.Context {
		return setMaxDuration(ctx, maxTime)
	})
}
