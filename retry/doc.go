// Package retry contains set of cliware middlewares and other utilities used
// for retry logic. Any of it to actually be applied, client has to use
// RoundTripper returned via retry.NewRoundTripper function as it implements
// all logic and applies all config values set by middlewares.
//
// This RoundTripper wraps your actual round tripper and uses it to send
// requests, so you can still apply you custom settings.
//
// Setting this RoundTripper is something that client library should do
// (like github.com/delicb/gwc), so, in general, for writing simple endpoints
// there should be no need to use it.
package retry
