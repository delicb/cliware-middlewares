package retry

import (
	"context"
	"time"
)

type retryConfigKey string

var (
	retryTimesKey retryConfigKey = "retry-times"
	maxDurationKey retryConfigKey = "max-time"
	backoffKey    retryConfigKey = "backoff"
	classifierKey retryConfigKey = "classifier"
)

func setRetryTimes(ctx context.Context, times int) context.Context {
	return context.WithValue(ctx, retryTimesKey, times)
}

func getRetryTimes(ctx context.Context) int {
	times := ctx.Value(retryTimesKey)
	if times == nil {
		return 1 // default number of retries
	}
	return times.(int)
}

func setMaxDuration(ctx context.Context, maxTime time.Duration) context.Context {
	return context.WithValue(ctx, maxDurationKey, maxTime)
}

func geMaxDuration(ctx context.Context) time.Duration {
	maxTime := ctx.Value(maxDurationKey)
	if maxTime == nil {
		return time.Duration(0)
	}
	return maxTime.(time.Duration)
}

func setBackoff(ctx context.Context, backoff BackoffStrategy) context.Context {
	return context.WithValue(ctx, backoffKey, backoff)
}

func getBackoff(ctx context.Context) BackoffStrategy {
	backoff := ctx.Value(backoffKey)
	if backoff == nil {
		return nil
	}
	return backoff.(BackoffStrategy)
}

func setClassifier(ctx context.Context, classifier Classifier) context.Context {
	return context.WithValue(ctx, classifierKey, classifier)
}

func getClassifier(ctx context.Context) Classifier {
	classifier := ctx.Value(classifierKey)
	if classifier == nil {
		return nil
	}
	return classifier.(Classifier)
}
