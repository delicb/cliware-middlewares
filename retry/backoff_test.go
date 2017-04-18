package retry_test

import (
	"testing"
	"time"

	"github.com/delicb/cliware-middlewares/retry"
)

func TestConstantBackoff(t *testing.T) {
	for _, constant := range []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10} {
		expect := time.Duration(constant)
		constantBackoff := retry.ConstantBackoff(expect)
		for i := 0; i < 10; i++ {
			got := constantBackoff(i)
			if constantBackoff(i) != expect {
				t.Errorf("Wrong backoff. Got: %s, expected: %s.", got, expect)
			}
		}
	}
}

func TestLinearBackoff(t *testing.T) {
	for _, data := range []struct {
		Step    time.Duration
		Max     time.Duration
		Attempt int
		Result  time.Duration
	}{
		{
			Step:    time.Second,
			Max:     time.Minute,
			Attempt: 3,
			Result:  3 * time.Second,
		},
		{
			Step:    time.Second,
			Max:     3 * time.Second,
			Attempt: 5,
			Result:  3 * time.Second,
		},
	} {
		linearBackoff := retry.LinearBackoff(data.Step, data.Max)
		got := linearBackoff(data.Attempt)
		if got != data.Result {
			t.Errorf("Got wrong backof. Got: %s, expected: %s.", got, data.Result)
		}
	}
}

func TestLinearJitterBackoff(t *testing.T) {
	for _, data := range []struct {
		Step      time.Duration
		Max       time.Duration
		Attempt   int
		ResultMin time.Duration
		ResultMax time.Duration
	}{
		{
			Step:      time.Second,
			Max:       time.Minute,
			Attempt:   3,
			ResultMin: 2 * time.Second,
			ResultMax: 4 * time.Second,
		},
		{
			Step:      time.Second,
			Max:       3 * time.Second,
			Attempt:   5,
			ResultMin: time.Second,
			ResultMax: 3 * time.Second,
		},
	} {
		linearBackoff := retry.LinearJitterBackoff(data.Step, data.Max)
		got := linearBackoff(data.Attempt)
		if got > data.ResultMax || got < data.ResultMin {
			t.Errorf("Got wrong backoff. Got: %s, expected result between %s and %s.",
				got, data.ResultMin, data.ResultMax)
		}
	}
}

func TestExponentialBackoff(t *testing.T) {
	for _, data := range []struct {
		Min     time.Duration
		Max     time.Duration
		Factor  float64
		Attempt int
		Result  time.Duration
	}{
		{
			Min:     time.Second,
			Max:     time.Minute,
			Factor:  2,
			Attempt: 3,
			Result:  8 * time.Second,
		},
		{
			Min:     time.Second,
			Max:     3 * time.Second,
			Factor:  2,
			Attempt: 3,
			Result:  3 * time.Second,
		},
		{
			Min:     time.Second,
			Max:     time.Minute,
			Factor:  1.5,
			Attempt: 3,
			Result:  time.Duration(3375) * time.Millisecond,
		},
	} {
		exponentialBackoff := retry.ExponentialBackoff(data.Min, data.Max, data.Factor)
		got := exponentialBackoff(data.Attempt)
		if got != data.Result {
			t.Errorf("Got wrong backoff. Got: %s, expected: %s.", got, data.Result)
		}
	}
}

func TestExponentialJitterBackoff(t *testing.T) {
	for _, data := range []struct {
		Min       time.Duration
		Max       time.Duration
		Factor    float64
		Attempt   int
		ResultMin time.Duration
		ResultMax time.Duration
	}{
		{
			Min:       time.Second,
			Max:       time.Minute,
			Factor:    2,
			Attempt:   3,
			ResultMin: 7 * time.Second,
			ResultMax: 9 * time.Second,
		},
		{
			Min:       time.Second,
			Max:       3 * time.Second,
			Factor:    2,
			Attempt:   3,
			ResultMin: 2 * time.Second,
			ResultMax: 4 * time.Second,
		},
	} {
		exponentialBackoff := retry.ExponentialJitterBackoff(data.Min, data.Max, data.Factor)
		got := exponentialBackoff(data.Attempt)
		if got > data.ResultMax || got < data.ResultMin {
			t.Errorf("Got wrong backoff. Got: %s, expected result between: %s and %s.",
				got, data.ResultMin, data.ResultMax)
		}
	}
}
