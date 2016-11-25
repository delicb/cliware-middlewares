package retry

import (
	"net/http"
	"testing"
	"time"

	"errors"
	"io"

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

func TestEnable(t *testing.T) {
	for _, client := range []*http.Client{
		http.DefaultClient,
		{
			Transport: http.DefaultTransport,
		},
		{
			Transport: NewRetryTransport(http.DefaultTransport),
		},
	} {
		Enable(client)
		if _, ok := client.Transport.(*retryTransport); !ok {
			t.Errorf("Wrong transport, expected retryTransport, got: %T.", client.Transport)
		}

	}
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
		//Mock          *mockRoundTripper
		ExpectedCalls int
		Classifier    Classifier
		Backoff       BackoffStrategy
		MaxRetries    int
		MaxDuration   time.Duration
		HTTPMethods   []string
		SendMethod    string
		BodyStrategy  RetryBodyStrategy
		ExpectedError string
	}{
		{
			//Mock:          &mockRoundTripper{},
			ExpectedCalls: 2,
			Classifier:    func(resp *http.Response, err error) bool { return true },
			Backoff:       func(n int) time.Duration { return time.Second },
			MaxRetries:    1,
			MaxDuration:   time.Minute,
			HTTPMethods:   []string{"GET"},
			SendMethod:    "GET",
		},
		{
			ExpectedCalls: 1,
			Classifier:    func(resp *http.Response, err error) bool { return true },
			SendMethod:    "POST",
		},
		{
			ExpectedCalls: 0,
			BodyStrategy:  func(r *http.Request) (func() io.ReadCloser, error) { return nil, errors.New("my error") },
			ExpectedError: "my error",
		},
	} {
		mock := &mockRoundTripper{}
		transport := NewRetryTransport(mock)
		req := cliware.EmptyRequest()
		req.Method = data.SendMethod
		req = req.WithContext(context.WithValue(req.Context(), classifierKey, data.Classifier))
		req = req.WithContext(context.WithValue(req.Context(), backoffKey, data.Backoff))
		req = req.WithContext(context.WithValue(req.Context(), maxDurationKey, data.MaxDuration))
		req = req.WithContext(context.WithValue(req.Context(), retryTimesKey, data.MaxRetries))
		req = req.WithContext(context.WithValue(req.Context(), bodyStrategyKey, data.BodyStrategy))

		_, err := transport.RoundTrip(req)
		if err != nil && data.ExpectedError == "" {
			t.Error("retryTransport returned error:", err)
		}
		if err != nil && data.ExpectedError != "" {
			if err.Error() != data.ExpectedError {
				t.Errorf("Wrong error. Got: %s, expected: %s.", err.Error(), data.ExpectedError)
			}
		}
		if mock.calledCount != data.ExpectedCalls {
			t.Errorf("Wrong number of calls. Got: %d, expected: %d.", mock.calledCount, data.ExpectedCalls)
		}
	}
}
