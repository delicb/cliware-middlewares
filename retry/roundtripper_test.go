package retry

import (
	"net/http"
	"testing"
	"time"

	"go.delic.rs/cliware"
	"golang.org/x/net/context"
)

type mockRoundTripper struct {
	calledCount int
	response    *http.Response
	err         error
}

func (rt *mockRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	rt.calledCount++
	return rt.response, rt.err
}

func TestNewRetryTransport(t *testing.T) {
	mock := &mockRoundTripper{}
	transport := NewRetryTransport(mock)
	transport.RoundTrip(cliware.EmptyRequest())
	if mock.calledCount == 0 {
		t.Error("Underlying RoundTripper not called.")
	}
}

func TestRetryTransport_RoundTrip(t *testing.T) {
	for _, data := range []struct {
		Mock          *mockRoundTripper
		ExpectedCalls int
		Classifier    Classifier
		Backoff       BackoffStrategy
		MaxRetries    int
		MaxDuration   time.Duration
	}{
		{
			Mock:          &mockRoundTripper{},
			ExpectedCalls: 12,
			Classifier:    func(resp *http.Response, err error) bool { return true },
			Backoff:       func(n int) time.Duration { return time.Second },
			MaxRetries:    0,
			MaxDuration:   time.Minute,
		},
	} {
		transport := NewRetryTransport(data.Mock)
		req := cliware.EmptyRequest()
		req = req.WithContext(context.WithValue(req.Context(), classifierKey, data.Classifier))
		req = req.WithContext(context.WithValue(req.Context(), backoffKey, data.Backoff))
		req = req.WithContext(context.WithValue(req.Context(), maxDurationKey, data.MaxDuration))
		req = req.WithContext(context.WithValue(req.Context(), retryTimesKey, data.MaxRetries))

		transport.RoundTrip(req)
		if data.Mock.calledCount != data.ExpectedCalls {
			t.Errorf("Wrong number of calls. Got: %d, expected: %d.", data.Mock.calledCount, data.ExpectedCalls)
		}
	}
}
