package errors_test

import (
	"bytes"
	"context"
	"net/http"
	"testing"

	"net/url"

	"io/ioutil"

	sterrors "errors"

	"regexp"

	"go.delic.rs/cliware"
	"go.delic.rs/cliware-middlewares/errors"
)

func TestErrors(t *testing.T) {
	for _, data := range []struct {
		Response      *http.Response
		OriginalError error
		Error         *errors.HTTPError
	}{
		{
			Response: &http.Response{
				StatusCode: 200,
			},
			OriginalError: nil,
			Error:         nil,
		},
		{
			Response: &http.Response{
				StatusCode: 302,
			},
			OriginalError: nil,
			Error:         nil,
		},
		{
			Response: &http.Response{
				StatusCode: 400,
				Request: &http.Request{
					Method: "POST",
					URL: &url.URL{
						Host: "delic.rs",
						Path: "/foobar",
					},
				},
			},
			OriginalError: nil,
			Error: &errors.HTTPError{
				StatusCode: 400,
				Method:     "POST",
			},
		},
		{
			Response: &http.Response{
				StatusCode: 401,
				Body:       ioutil.NopCloser(bytes.NewBuffer([]byte("body"))),
				Request: &http.Request{
					Method: "GET",
					URL: &url.URL{
						Host: "golang.com",
						Path: "/somepath",
					},
				},
			},
			OriginalError: nil,
			Error: &errors.HTTPError{
				StatusCode: 401,
				Method:     "GET",
				Body:       []byte("body"),
			},
		},
		{
			Response:      nil,
			OriginalError: sterrors.New("Custom error"),
			Error:         nil,
		},
	} {
		m := errors.Errors()
		req := cliware.EmptyRequest()
		handler := createHandler(data.Response, data.OriginalError)
		_, err := m.Exec(handler).Handle(nil, req)

		if data.Error == nil && data.OriginalError == nil {
			if err != nil {
				t.Errorf("Did not expect error, got: %s", err)
			}
			continue
		}

		if data.OriginalError != nil {
			if data.OriginalError.Error() != err.Error() {
				t.Errorf("Wrong error. Got: %s, expected: %s", err.Error(), data.OriginalError.Error())
			}
			continue
		}

		if httpErr, ok := err.(*errors.HTTPError); ok {
			if httpErr.StatusCode != data.Error.StatusCode {
				t.Errorf("Wrong status code. Got: %d, expected: %d", httpErr.StatusCode, data.Error.StatusCode)
			}
			if httpErr.Method != data.Error.Method {
				t.Errorf("Wrong method. Got: %s, expected: %s", httpErr.Method, data.Error.Method)
			}
			if httpErr != nil && string(httpErr.Body) != string(data.Error.Body) {
				t.Errorf("Wrong body. Got: %s, expected: %s", string(httpErr.Body), string(data.Error.Body))
			}
		} else {
			t.Errorf("Wrong error type. Expected HTTPError, got: %T", err)
		}
	}
}

func TestHTTPError_Error(t *testing.T) {
	for _, data := range []struct {
		Error    *errors.HTTPError
		Expected string
	}{
		{
			Error: &errors.HTTPError{
				Method: "POST",
			},
			Expected: ".*POST.*",
		},
		{
			Error: &errors.HTTPError{
				Name: "401 Forbidden",
			},
			Expected: ".*401.*",
		},
		{
			Error: &errors.HTTPError{
				RequestURL: "some_url",
			},
			Expected: ".*some_url.*",
		},
	} {
		errStr := data.Error.Error()
		match, err := regexp.Match(data.Expected, []byte(errStr))
		if err != nil {
			t.Error("Regex match failed: ", err)
		}
		if !match {
			t.Errorf("Wrong error string. Got: %s, did not match regexp: %s", errStr, data.Expected)
		}
	}
}

func createHandler(wantedResponse *http.Response, wantError error) cliware.Handler {
	return cliware.HandlerFunc(func(ctx context.Context, req *http.Request) (resp *http.Response, err error) {
		return wantedResponse, wantError
	})
}
