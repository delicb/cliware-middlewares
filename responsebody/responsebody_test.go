package responsebody_test

import (
	"errors"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"bytes"

	"github.com/delicb/cliware"
	"github.com/delicb/cliware-middlewares/responsebody"
)

func TestJSON(t *testing.T) {
	for _, data := range []struct {
		RawData  string
		Expected map[string]interface{}
		Error    error
	}{
		{
			RawData: `{"foo": "bar"}`,
			Expected: map[string]interface{}{
				"foo": "bar",
			},
			Error: nil,
		},
		{
			RawData:  `{"foo": "bar"}`,
			Expected: map[string]interface{}{},
			Error:    errors.New("Some error"),
		},
	} {
		var body map[string]interface{}
		req := cliware.EmptyRequest()
		handler := func(req *http.Request) (*http.Response, error) {
			r := &http.Response{
				Body: ioutil.NopCloser(strings.NewReader(data.RawData)),
			}
			return r, data.Error
		}

		_, err := responsebody.JSON(&body).Exec(cliware.HandlerFunc(handler)).Handle(req)
		if err != nil && data.Error == nil {
			t.Error("Got unexpected error: ", err)
		}
		if err == nil && data.Error != nil {
			t.Error("Did not get error and one was expected:", data.Error)
		}
		if err != nil && data.Error != nil {
			if err != data.Error {
				t.Errorf("Wrong error. Expected: %s, got: %s.", data.Error, err)
			}
		}
		// check for body != nil, since in error cases, body will not be populated
		if !reflect.DeepEqual(body, data.Expected) && body != nil {
			t.Errorf("Wrong response data. Expected: %#v, got: %#v", data.Expected, body)
		}
	}
}

func TestString(t *testing.T) {
	for _, data := range []struct {
		Data  string
		Error error
	}{
		{
			Data:  "foo bar",
			Error: nil,
		},
		{
			Data:  "foo bar",
			Error: errors.New("custom error"),
		},
	} {
		var body string
		req := cliware.EmptyRequest()
		handler := func(req *http.Request) (*http.Response, error) {
			r := &http.Response{
				Body: ioutil.NopCloser(strings.NewReader(data.Data)),
			}
			return r, data.Error
		}
		_, err := responsebody.String(&body).Exec(cliware.HandlerFunc(handler)).Handle(req)
		if err != nil && data.Error == nil {
			t.Error("Got unexpected error: ", err)
		}
		if err == nil && data.Error != nil {
			t.Error("Did not get error and one was expected:", data.Error)
		}
		if err != nil && data.Error != nil {
			if err != data.Error {
				t.Errorf("Wrong error. Expected: %s, got: %s.", data.Error, err)
			}
		}
		// check for body != "" since in error cases body will not be populated
		if body != data.Data && body != "" {
			t.Errorf("Wrong response data. Expected: %#v, got: %#v", data.Data, body)
		}
	}
}

func TestWriter(t *testing.T) {
	for _, data := range []struct {
		Data  string
		Error error
	}{
		{
			Data:  "foo bar",
			Error: nil,
		},
		{
			Data:  "foo bar",
			Error: errors.New("my error"),
		},
	} {
		buf := &bytes.Buffer{}
		req := cliware.EmptyRequest()
		handler := func(req *http.Request) (*http.Response, error) {
			r := &http.Response{
				Body: ioutil.NopCloser(strings.NewReader(data.Data)),
			}
			return r, data.Error
		}
		_, err := responsebody.Writer(buf).Exec(cliware.HandlerFunc(handler)).Handle(req)
		if err != nil && data.Error == nil {
			t.Error("Got unexpected error: ", err)
		}
		if err == nil && data.Error != nil {
			t.Error("Did not get error and one was expected:", data.Error)
		}
		if err != nil && data.Error != nil {
			if err != data.Error {
				t.Errorf("Wrong error. Expected: %s, got: %s.", data.Error, err)
			}
		}
		got := buf.String()
		// checking for data.Error != nil since on error body will be empty.
		if got != data.Data && data.Error == nil {
			t.Errorf("Wrong response data. Expected: %s, got: %s", data.Data, got)
		}
	}
}
