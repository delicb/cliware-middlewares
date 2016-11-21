package retry

import (
	"context"
	"net/http"
	"testing"

	"time"

	"go.delic.rs/cliware"
)

func TestTimes(t *testing.T) {
	for _, times := range []int{0, 1, 2, 3, 4, 5, 100} {
		m := Times(times)
		resultContext := context.Background()
		initialContext := context.Background()
		m.Exec(createHandler(&resultContext)).Handle(initialContext, nil)
		got := getRetryTimes(resultContext)
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
		resultContext := context.Background()
		initialContext := context.Background()
		m.Exec(createHandler(&resultContext)).Handle(initialContext, nil)
		got := getClassifier(resultContext)
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
		resultContext := context.Background()
		initialContext := context.Background()
		m.Exec(createHandler(&resultContext)).Handle(initialContext, nil)
		got := getBackoff(resultContext)
		// can not really compare functions, so just check if we not non-nill value
		if got == nil {
			t.Error("Wrong backoff strategy.. Got nil")
		}
	}
}

func TestMaxDuration(t *testing.T) {
	for _, duration := range []time.Duration{
		time.Second, time.Minute, 2 * time.Hour, 3 * time.Hour, 0,
	} {
		m := MaxDuration(duration)
		var resultContext context.Context
		initialContext := context.Background()
		m.Exec(createHandler(&resultContext)).Handle(initialContext, nil)
		got := geMaxDuration(resultContext)
		if got != duration {
			t.Errorf("Wrong max duration. Got: %s, expected: %s.", got, duration)
		}
	}
}

func createHandler(resultContext *context.Context) cliware.Handler {
	return cliware.HandlerFunc(func(ctx context.Context, req *http.Request) (resp *http.Response, err error) {
		*resultContext = ctx
		return nil, nil
	})
}
