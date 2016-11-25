package retry

import (
	"context"
	"net/http"

	"time"

	"io"

	c "go.delic.rs/cliware"
)

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

// BodyStrategy sets strategy of how to handle request body for retries requests.
func BodyStrategy(strategy func(r *http.Request) (func() io.ReadCloser, error)) c.Middleware {
	return c.ContextProcessor(func(ctx context.Context) context.Context {
		return setBodyStrategy(ctx, RetryBodyStrategy(strategy))
	})
}

// Methods sets list of HTTP methods for which it is valid to retry failed requests.
// For example, if only GET is defined as valid retry HTTP methods, no POST or
// PUT (or any other) request will be retried. This might be useful if you are not
// sure about server idempotence.
func Methods(methods ...string) c.Middleware {
	return c.ContextProcessor((func(ctx context.Context) context.Context {
		return setRetryMethods(ctx, methods...)
	}))
}
