package retry

import (
	"context"
	"time"
)

// retryConfigKey is private type to be used for storing information in context
// and be sure that there will be no collision with other keys
type retryConfigKey string

var (
	retryTimesKey   retryConfigKey = "retry-times"
	maxDurationKey  retryConfigKey = "max-time"
	backoffKey      retryConfigKey = "backoff"
	classifierKey   retryConfigKey = "classifier"
	bodyStrategyKey retryConfigKey = "body-strategy"
	retryMethodsKey retryConfigKey = "retry-methods"
)

// setRetryTimes sets provided number of retry times to provided context and
// returns new context.
func setRetryTimes(ctx context.Context, times int) context.Context {
	return context.WithValue(ctx, retryTimesKey, times)
}

// getRetryTimes returns number of retry times from provided context or 1,
// if provided context does not contain value for number of retries.
func getRetryTimes(ctx context.Context) int {
	times := ctx.Value(retryTimesKey)
	if times == nil {
		return 1 // default number of retries
	}
	return times.(int)
}

// setMaxDuration sets provided maximum retry duration to provided context and
// returns new context.
func setMaxDuration(ctx context.Context, maxTime time.Duration) context.Context {
	return context.WithValue(ctx, maxDurationKey, maxTime)
}

// getMaxDuration returns maximum retry duration from provided context or
// time.Duration(0) if provided context does not contain value for maximum
// retry duration time.
func getMaxDuration(ctx context.Context) time.Duration {
	maxTime := ctx.Value(maxDurationKey)
	if maxTime == nil {
		return time.Duration(0)
	}
	return maxTime.(time.Duration)
}

// setBackoff sets provided backoff strategy to provided context and returns
// new context.
func setBackoff(ctx context.Context, backoff BackoffStrategy) context.Context {
	return context.WithValue(ctx, backoffKey, backoff)
}

// getBackoff returns backoff strategy from provided context or nil if
// provided context does not contain value for backoff strategy.
func getBackoff(ctx context.Context) BackoffStrategy {
	backoff := ctx.Value(backoffKey)
	if backoff == nil {
		return nil
	}
	return backoff.(BackoffStrategy)
}

// setClassifier sets provided classifier to provided context and returns
// new context.
func setClassifier(ctx context.Context, classifier Classifier) context.Context {
	return context.WithValue(ctx, classifierKey, classifier)
}

// getClassifier returns classifier from provided context or nil if provided
// context does not contain value for classifier.
func getClassifier(ctx context.Context) Classifier {
	classifier := ctx.Value(classifierKey)
	if classifier == nil {
		return nil
	}
	return classifier.(Classifier)
}

// setBodyStrategy sets provided body strategy to provided context and returns
// new context.
func setBodyStrategy(ctx context.Context, bodyStrategy BodyStrategy) context.Context {
	return context.WithValue(ctx, bodyStrategyKey, bodyStrategy)
}

// getBodyStrategy returns body strategy from provided context or nil if provided
// context does not contain value for body strategy.
func getBodyStrategy(ctx context.Context) BodyStrategy {
	bodyStrategy := ctx.Value(bodyStrategyKey)
	if bodyStrategy == nil {
		return nil
	}
	return bodyStrategy.(BodyStrategy)
}

// setRetryMethods sets provided list of HTTP methods to provided context as
// methods acceptable to retry on and returns new context.
func setRetryMethods(ctx context.Context, methods ...string) context.Context {
	return context.WithValue(ctx, retryMethodsKey, methods)
}

// getRetryMethods returns slice of methods that are acceptable to retry on or
// nil if provided context does not contain value for methods.
func getRetryMethods(ctx context.Context) []string {
	retryMethods := ctx.Value(retryMethodsKey)
	if retryMethods == nil {
		return nil
	}
	return retryMethods.([]string)
}
