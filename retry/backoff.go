package retry

import (
	"math"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// BackoffStrategy is function that calculates for how long should delay be before
// sending next request.
type BackoffStrategy func(n int) time.Duration

// ConstantBackoff returns backoff that returns always same, provided, parameter.
func ConstantBackoff(duration time.Duration) BackoffStrategy {
	return BackoffStrategy(func(_ int) time.Duration {
		return duration
	})
}

// LinearBackoff returns BackoffStrategy that multiplies provided step with repeat
// attempt number.
func LinearBackoff(step time.Duration, max time.Duration) BackoffStrategy {
	return BackoffStrategy(func(n int) time.Duration {
		return minDuration(calcLinear(step, n), max)
	})
}

// LinearJitterBackoff returns BackoffStrategy that multiplies provided step with
// attempt number and adds random jitter to result.
func LinearJitterBackoff(step time.Duration, max time.Duration) BackoffStrategy {
	return BackoffStrategy(func(n int) time.Duration {
		return minDuration(addJitter(calcLinear(step, n), step), max)
	})
}

//ExponentialBackoff returns BackoffStrategy that exponentially increases time.
func ExponentialBackoff(min, max time.Duration, factor float64) BackoffStrategy {
	return BackoffStrategy(func(n int) time.Duration {
		return calcExponential(min, max, factor, n)
	})
}

//ExponentialJitterBackoff returns BackoffStrategy that exponentially increases time and
// adds random jitter to result.
func ExponentialJitterBackoff(min, max time.Duration, factor float64) BackoffStrategy {
	return BackoffStrategy(func(n int) time.Duration {
		return addJitter(calcExponential(min, max, factor, n), min)
	})
}

func calcExponential(min, max time.Duration, factor float64, attempt int) time.Duration {
	d := time.Duration(float64(min) * math.Pow(factor, float64(attempt)))
	return minDuration(d, max)
}

func calcLinear(step time.Duration, attempt int) time.Duration {
	return time.Duration(attempt) * step
}

func minDuration(first time.Duration, second time.Duration) time.Duration {
	if first.Nanoseconds() > second.Nanoseconds() {
		return second
	}
	return first
}

func addJitter(next, min time.Duration) time.Duration {
	return time.Duration(rand.Float64()*float64(2*min) + float64(next-min))
}
