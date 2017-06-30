package retry

import (
	"context"
	"net/http"
	"testing"

	"time"

	"io"
	"reflect"

	"github.com/delicb/cliware"
)

func TestTimes(t *testing.T) {
	for _, times := range []int{0, 1, 2, 3, 4, 5, 100} {
		m := Times(times)
		initialContext := context.Background()
		req := cliware.EmptyRequest().WithContext(initialContext)
		resp, err := m.Exec(createHandler()).Handle(req)
		if err != nil {
			t.Error("Handle returned error:", err)
		}
		got := getRetryTimes(resp.Request.Context())
		if got != times {
			t.Errorf("Wrong number of retries. Got %d, expected: %d.", got, times)
		}
	}
}

func TestSetClassifier(t *testing.T) {
	for _, classifier := range []Classifier{
		func(resp *http.Response, err error) bool { return true },
	} {
		m := SetClassifier(classifier)
		initialContext := context.Background()
		req := cliware.EmptyRequest().WithContext(initialContext)
		resp, err := m.Exec(createHandler()).Handle(req)
		if err != nil {
			t.Error("Handle returned error:", err)
		}
		got := getClassifier(resp.Request.Context())
		// can not really compare functions, so just check if we not non-nill value
		if got == nil {
			t.Error("Wrong classifier. Got nil")
		}
	}
}

func TestSetBackoffStrategy(t *testing.T) {
	for _, backoff := range []BackoffStrategy{
		func(n int) time.Duration { return time.Second },
	} {
		m := SetBackoffStrategy(backoff)
		initialContext := context.Background()
		req := cliware.EmptyRequest().WithContext(initialContext)
		resp, err := m.Exec(createHandler()).Handle(req)
		if err != nil {
			t.Error("Handle returned error:", err)
		}
		got := getBackoff(resp.Request.Context())
		// can not really compare functions, so just check if we not non-nil value
		if got == nil {
			t.Error("Wrong backoff strategy. Got nil")
		}
	}
}

func TestMaxDuration(t *testing.T) {
	for _, duration := range []time.Duration{
		time.Second, time.Minute, 2 * time.Hour, 3 * time.Hour, 0,
	} {
		m := MaxDuration(duration)
		initialContext := context.Background()
		req := cliware.EmptyRequest().WithContext(initialContext)
		resp, err := m.Exec(createHandler()).Handle(req)
		if err != nil {
			t.Error("Handle returned error:", err)
		}
		got := getMaxDuration(resp.Request.Context())
		if got != duration {
			t.Errorf("Wrong max duration. Got: %s, expected: %s.", got, duration)
		}
	}
}

func TestBodyStrategy(t *testing.T) {
	for _, strategy := range []BodyStrategy{
		BodyStrategy(func(r *http.Request) (func() io.ReadCloser, error) { return nil, nil }),
	} {
		m := SetBodyStrategy(strategy)
		initialContext := context.Background()
		req := cliware.EmptyRequest().WithContext(initialContext)
		resp, err := m.Exec(createHandler()).Handle(req)
		if err != nil {
			t.Error("Handle returned error:", err)
		}
		got := getBodyStrategy(resp.Request.Context())
		// can not really compare functions, so just check if we got non-nil value
		if got == nil {
			t.Error("Wrong body strategy. Got nil.")
		}
	}
}

func TestMethods(t *testing.T) {
	for _, methods := range [][]string{
		{},
		{"GET"},
		{"GET", "POST", "PUT"},
	} {
		m := Methods(methods...)
		initialContext := context.Background()
		req := cliware.EmptyRequest().WithContext(initialContext)
		resp, err := m.Exec(createHandler()).Handle(req)
		if err != nil {
			t.Error("Handle returned error:", err)
		}
		got := getRetryMethods(resp.Request.Context())
		if !reflect.DeepEqual(got, methods) {
			t.Errorf("Wrong HTTP methods. Got: %s, expected: %s.", got, methods)
		}
	}
}

func createHandler() cliware.Handler {
	return cliware.HandlerFunc(func(req *http.Request) (resp *http.Response, err error) {
		return &http.Response{
			Request: req,
		}, nil
	})
}
