package body_test

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"encoding/json"
	"encoding/xml"
	"reflect"

	"bytes"

	"strings"

	"go.delic.rs/cliware"
	"go.delic.rs/cliware-middlewares/body"
)

func TestString(t *testing.T) {
	for _, data := range []struct {
		Data           string
		Method         string
		ExpectedMethod string
	}{
		{"my content", "GET", "POST"},
		{"", "POST", ""},
		{"š", "GET", "POST"},
	} {
		middleware := body.String(data.Data)
		ctx := context.Background()
		req := cliware.EmptyRequest()
		req.Method = data.Method
		h := createHandler()
		_, err := middleware.Exec(h).Handle(ctx, req)
		if err != nil {
			t.Error("Got error processing request: ", err)
		}

		expectedMethod := data.ExpectedMethod
		if expectedMethod == "" {
			expectedMethod = data.Method
		}
		if req.Method != expectedMethod {
			t.Errorf("Wrong method on request. Expected: \"%s\", got \"%s\"", expectedMethod, req.Method)
		}
		contentLength := int64(len([]byte(data.Data)))
		if req.ContentLength != contentLength {
			t.Errorf("Wrong content length. Expected: %d, got %d.", contentLength, req.ContentLength)
		}
		bodyBytes, _ := ioutil.ReadAll(req.Body)
		body := string(bodyBytes)
		if body != data.Data {
			t.Errorf("Wrong body. Expected: \"%s\", got: \"%s\".", data.Data, body)
		}
	}
}

func TestJSON(t *testing.T) {
	for _, data := range []struct {
		Data          interface{}
		ResponseData  interface{}
		ExpectedError error
	}{
		{`{"foo": "bar"}`, &map[string]string{"foo": "bar"}, nil},
		{[]byte(`{"foo": "bar"}`), &map[string]string{"foo": "bar"}, nil},
		{
			struct {
				Foo string `json:"foo"`
			}{
				Foo: "bar",
			}, &map[string]string{"foo": "bar"}, nil},
	} {
		req := cliware.EmptyRequest()
		handler := createHandler()
		_, err := body.JSON(data.Data).Exec(handler).Handle(nil, req)
		if data.ExpectedError != nil {
			if data.ExpectedError != err {
				t.Errorf("Wrong error. Expected: %s, got: %s", data.ExpectedError, err)
			}
		} else {
			if err != nil {
				t.Error("Got unexpected error processing request: ", err)
			}
		}

		if req.Method != "POST" {
			t.Error("Wrong request method. Expected: POST, got: ", req.Method)
		}
		if req.Header.Get("Content-Type") != "application/json" {
			t.Error("Wrong content-type. Expected application/json, got: ", req.Header.Get("Content-Type"))
		}
		bodyBytes, _ := ioutil.ReadAll(req.Body)
		body := new(map[string]string)
		err = json.Unmarshal(bodyBytes, body)
		if err != nil {
			t.Error("Failed to unmarshal response json: ", err)
		}
		if !reflect.DeepEqual(body, data.ResponseData) {
			t.Errorf("Wrong body. Expected: %v, got: %v.", data.ResponseData, body)
		}
	}
}

func TestXML(t *testing.T) {
	type person struct {
		Name    string
		DOB     string
		XMLName string `xml:"Person"`
	}

	type testCase struct {
		Data          interface{}
		ResponseData  *person
		ExpectedError error
	}

	tests := []testCase{
		{
			Data:          `<Person><Name>Foobar Barfoo</Name><DOB>11-12-1984</DOB></Person>`,
			ResponseData:  &person{Name: "Foobar Barfoo", DOB: "11-12-1984"},
			ExpectedError: nil,
		},
		{
			Data:          []byte(`<Person><Name>Foobar Barfoo</Name><DOB>11-12-1984</DOB></Person>`),
			ResponseData:  &person{Name: "Foobar Barfoo", DOB: "11-12-1984"},
			ExpectedError: nil,
		},
		{
			Data:          &person{Name: "Foobar Barfoo", DOB: "11-12-1984"},
			ResponseData:  &person{Name: "Foobar Barfoo", DOB: "11-12-1984"},
			ExpectedError: nil,
		},
	}

	for _, data := range tests {
		req := cliware.EmptyRequest()
		handler := createHandler()
		_, err := body.XML(data.Data).Exec(handler).Handle(nil, req)
		if data.ExpectedError != nil {
			if data.ExpectedError != err {
				t.Errorf("Wrong error. Expected: %s, got: %s", data.ExpectedError, err)
			}
		} else {
			if err != nil {
				t.Error("Got unexpected error processing request: ", err)
			}
		}

		if req.Method != "POST" {
			t.Error("Wrong request method. Expected: POST, got: ", req.Method)
		}
		if req.Header.Get("Content-Type") != "application/xml" {
			t.Error("Wrong content-type. Expected application/xml, got: ", req.Header.Get("Content-Type"))
		}
		bodyBytes, _ := ioutil.ReadAll(req.Body)
		body := &person{}
		err = xml.Unmarshal(bodyBytes, body)
		if err != nil {
			t.Error("Failed to unmarshal body xml: ", err)
		}
		if !reflect.DeepEqual(body, data.ResponseData) {
			t.Errorf("Wrong body. Expected: %v, got: %v.", data.ResponseData, body)
		}
	}
}

func TestReader(t *testing.T) {
	for _, data := range []struct {
		Reader       io.Reader
		ResponseData interface{}
		ExpectError  error
	}{
		{
			Reader:       bytes.NewBuffer([]byte("data")),
			ResponseData: "data",
			ExpectError:  nil,
		},
		{
			Reader:       bytes.NewReader([]byte("data")),
			ResponseData: "data",
			ExpectError:  nil,
		},
		{
			Reader:       strings.NewReader("data"),
			ResponseData: "data",
			ExpectError:  nil,
		},
	} {
		req := cliware.EmptyRequest()
		handler := createHandler()
		_, err := body.Reader(data.Reader).Exec(handler).Handle(nil, req)
		if data.ExpectError != nil {
			if data.ExpectError.Error() != err.Error() {
				t.Errorf("Wrong error. Expected: %s, got: %s", data.ExpectError.Error(), err.Error())
			}
		}
		bodyBytes, _ := ioutil.ReadAll(req.Body)
		body := string(bodyBytes)
		if !reflect.DeepEqual(body, data.ResponseData) {
			t.Errorf("Wrong body. Expected: %s, got: %s", data.ResponseData, body)
		}
	}
}

func createHandler() cliware.Handler {
	return cliware.HandlerFunc(func(ctx context.Context, req *http.Request) (resp *http.Response, err error) {
		return nil, nil
	})
}
