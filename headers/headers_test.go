package headers_test

import (
	"context"
	"net/http"
	"testing"

	"reflect"

	"time"

	"github.com/delicb/cliware"
	"github.com/delicb/cliware-middlewares/headers"
)

func TestMethod(t *testing.T) {
	for _, method := range []string{
		"GET", "POST", "PUT", "ANYTHING",
	} {
		m := headers.Method(method)
		req := cliware.EmptyRequest()
		handler := createHandler()
		m.Exec(handler).Handle(nil, req)

		if req.Method != method {
			t.Errorf("Wrong method. Got: %s, expected: %s", req.Method, method)
		}
	}
}

func TestAdd(t *testing.T) {
	for _, data := range []struct {
		Name             string
		Value            string
		ExpectedValue    []string
		Existing         http.Header
		ExpectingHeaders int
	}{
		{
			Name:             "Content-Type",
			Value:            "application/json",
			ExpectedValue:    []string{"application/json"},
			Existing:         http.Header{},
			ExpectingHeaders: 1,
		},
		{
			Name:          "Custom-Header",
			Value:         "whatever",
			ExpectedValue: []string{"whatever"},
			Existing: http.Header{
				"Content-Type": []string{"text/html"},
				"Host":         []string{"bojan.delic.rs"},
			},
			ExpectingHeaders: 3,
		},
		{
			Name:          "Content-Type",
			Value:         "application/json",
			ExpectedValue: []string{"text/html", "application/json"},
			Existing: http.Header{
				"Content-Type": []string{"text/html"},
				"Host":         []string{"bojan.delic.rs"},
			},
			ExpectingHeaders: 2,
		},
	} {
		m := headers.Add(data.Name, data.Value)
		req := cliware.EmptyRequest()
		// assigning data.Existing to req.Header will not work since it is map
		// and same data.Existing would be modified during execution of middleware
		for k, v := range data.Existing {
			for _, vv := range v {
				req.Header.Set(k, vv)
			}
		}

		handler := createHandler()
		m.Exec(handler).Handle(nil, req)

		if len(req.Header) != data.ExpectingHeaders {
			t.Errorf("Number of headers to not match. Got: %d, expected: %d.", len(req.Header), data.ExpectingHeaders)
		}

		res, ok := req.Header[data.Name]
		if !ok {
			t.Errorf("Header \"%s\" not found.", data.Name)
		}
		if !reflect.DeepEqual(res, data.ExpectedValue) {
			t.Errorf("Wrong value for added header. Got: %s, expected: %s", res, data.ExpectedValue)
		}
	}
}

func TestSet(t *testing.T) {
	for _, data := range []struct {
		Name             string
		Value            string
		ExpectedValue    []string
		Existing         http.Header
		ExpectingHeaders int
	}{
		{
			Name:             "Content-Type",
			Value:            "application/json",
			ExpectedValue:    []string{"application/json"},
			Existing:         http.Header{},
			ExpectingHeaders: 1,
		},
		{
			Name:          "Custom-Header",
			Value:         "whatever",
			ExpectedValue: []string{"whatever"},
			Existing: http.Header{
				"Content-Type": []string{"text/html"},
				"Host":         []string{"bojan.delic.rs"},
			},
			ExpectingHeaders: 3,
		},
		{
			Name:          "Content-Type",
			Value:         "application/json",
			ExpectedValue: []string{"application/json"},
			Existing: http.Header{
				"Content-Type": []string{"text/html"},
				"Host":         []string{"bojan.delic.rs"},
			},
			ExpectingHeaders: 2,
		},
	} {
		m := headers.Set(data.Name, data.Value)
		req := cliware.EmptyRequest()
		// assigning data.Existing to req.Header will not work since it is map
		// and same data.Existing would be modified during execution of middleware
		for k, v := range data.Existing {
			for _, vv := range v {
				req.Header.Set(k, vv)
			}
		}

		handler := createHandler()
		m.Exec(handler).Handle(nil, req)

		if len(req.Header) != data.ExpectingHeaders {
			t.Errorf("Number of headers to not match. Got: %d, expected: %d.", len(req.Header), data.ExpectingHeaders)
		}

		res, ok := req.Header[data.Name]
		if !ok {
			t.Errorf("Header \"%s\" not found.", data.Name)
		}
		if !reflect.DeepEqual(res, data.ExpectedValue) {
			t.Errorf("Wrong value for added header. Got: %s, expected: %s", res, data.ExpectedValue)
		}
	}
}

func TestDel(t *testing.T) {
	for _, data := range []struct {
		Name     string
		Existing http.Header
	}{
		{
			Name:     "Content-Type",
			Existing: http.Header{},
		},
		{
			Name: "Custom-Header",
			Existing: http.Header{
				"Content-Type": []string{"text/html"},
				"Host":         []string{"bojan.delic.rs"},
			},
		},
		{
			Name: "Content-Type",
			Existing: http.Header{
				"Content-Type": []string{"text/html"},
				"Host":         []string{"bojan.delic.rs"},
			},
		},
	} {
		m := headers.Del(data.Name)
		req := cliware.EmptyRequest()
		// assigning data.Existing to req.Header will not work since it is map
		// and same data.Existing would be modified during execution of middleware
		for k, v := range data.Existing {
			for _, vv := range v {
				req.Header.Set(k, vv)
			}
		}

		handler := createHandler()
		m.Exec(handler).Handle(nil, req)

		_, ok := req.Header[data.Name]
		if ok {
			t.Errorf("Header \"%s\" even when it should be deleted.", data.Name)
		}
	}
}

func TestSetMap(t *testing.T) {
	for _, data := range []struct {
		ToSet    map[string]string
		Existing http.Header
		Expected http.Header
	}{
		{
			ToSet:    map[string]string{},
			Existing: http.Header{},
			Expected: http.Header{},
		},
		{
			ToSet: map[string]string{
				"Content-Type": "application/gwc",
			},
			Existing: http.Header{},
			Expected: http.Header{
				"Content-Type": []string{"application/gwc"},
			},
		},
		{
			ToSet: map[string]string{
				"Content-Type": "application/gwc",
			},
			Existing: http.Header{
				"Host": []string{"bojan.delic.rs"},
			},
			Expected: http.Header{
				"Content-Type": []string{"application/gwc"},
				"Host":         []string{"bojan.delic.rs"},
			},
		},
		{
			ToSet: map[string]string{
				"Host": "delic.rs",
			},
			Existing: http.Header{
				"Host": []string{"bojan.delic.rs"},
			},
			Expected: http.Header{
				"Host": []string{"delic.rs"},
			},
		},
	} {
		m := headers.SetMap(data.ToSet)
		req := cliware.EmptyRequest()
		// assigning data.Existing to req.Header will not work since it is map
		// and same data.Existing would be modified during execution of middleware
		for k, v := range data.Existing {
			for _, vv := range v {
				req.Header.Set(k, vv)
			}
		}

		handler := createHandler()
		m.Exec(handler).Handle(nil, req)

		if !reflect.DeepEqual(data.Expected, req.Header) {
			t.Errorf("Wrong headers. Got: %s, expected: %s.", req.Header, data.Expected)
		}
	}
}

func TestFromContext_HeaderList(t *testing.T) {
	for _, data := range []struct {
		Key           string
		Value         interface{}
		ExpectedError string
	}{
		{
			Key: "some-id",
			Value: []headers.Header{
				{
					Key:   "Some header",
					Value: []string{"My value"},
				},
			},
			ExpectedError: "",
		},
		{
			Key: "some-id",
			Value: []headers.Header{
				{
					Key:   "Some header",
					Value: []string{"My value"},
				},
				{
					Key:   "My header",
					Value: []string{"other value"},
				},
			},
			ExpectedError: "",
		},
		{
			Key:           "some-id",
			Value:         "string",
			ExpectedError: "headers.FromContext middleware: value in unsupported format",
		},
	} {
		m := headers.FromContext(data.Key)
		req := cliware.EmptyRequest()
		var ctx context.Context
		if headers, ok := data.Value.([]headers.Header); ok {
			if len(headers) == 1 {
				ctx = context.WithValue(context.Background(), data.Key, headers[0])
			} else {
				ctx = context.WithValue(context.Background(), data.Key, headers)
			}
		} else {
			ctx = context.WithValue(context.Background(), data.Key, data.Value)
		}

		_, err := m.Exec(createHandler()).Handle(ctx, req)
		if err != nil {
			if data.ExpectedError == "" {
				t.Errorf("Did not expect error, got: %s.", data.ExpectedError)
			} else if err.Error() != data.ExpectedError {
				t.Errorf("Got wrong error. Got: %s, expected: %s.", err, data.ExpectedError)
			}
		} else {
			for _, h := range data.Value.([]headers.Header) {
				rawHeader, ok := req.Header[h.Key]
				if !ok {
					t.Fatalf("Header %s not found in request.", h.Key)
				}
				if !reflect.DeepEqual(rawHeader, h.Value) {
					t.Errorf("Wrong header value set. Got: %v, expected: %v.", rawHeader, h.Value)
				}
			}
		}
	}
}

func TestToContext(t *testing.T) {
	for _, data := range []struct {
		Key    interface{}
		Header string
		Value  []string
	}{
		{
			Key:    "some-id",
			Header: "Header",
			Value:  []string{"value"},
		},
		{
			Key:    1,
			Header: "",
			Value:  []string{},
		},
		{
			Key:    time.Now(),
			Header: "My-Header",
			Value:  []string{"val1", "val2"},
		},
	} {
		ctx := headers.ToContext(context.Background(), data.Key, data.Header, data.Value...)
		val := ctx.Value(data.Key)
		if header, ok := val.(headers.Header); ok {
			if !reflect.DeepEqual(header.Value, data.Value) {
				t.Errorf("Got wrong value for header. Got: %v, expected: %v.", header.Value, data.Value)
			}
		} else {
			t.Errorf("Got wrong type from context for key. Got: %T, expected %T", val, headers.Header{})
		}
	}
}

func TestToContextList(t *testing.T) {
	for _, data := range []struct {
		Key     interface{}
		Headers []headers.Header
	}{
		{
			Key: "some-id",
			Headers: []headers.Header{
				{
					Key:   "myheader",
					Value: []string{"value"},
				},
			},
		},
		{
			Key: "",
			Headers: []headers.Header{
				{
					Key:   "",
					Value: []string{},
				},
			},
		},
		{
			Key: 1,
			Headers: []headers.Header{
				{
					Key:   "My-Header",
					Value: []string{"val1", "val2", "val3"},
				},
			},
		},
	} {
		ctx := headers.ToContextList(context.Background(), data.Key, data.Headers)
		val := ctx.Value(data.Key)
		if header, ok := val.([]headers.Header); ok {
			if !reflect.DeepEqual(header, data.Headers) {
				t.Errorf("Got wrong value for header. Got: %v, expected: %v.", header, data.Headers)
			}
		} else {
			t.Errorf("Got wrong type from context for key. Got %T, expected: %T.", val, []headers.Header{})
		}
	}
}

func createHandler() cliware.Handler {
	return cliware.HandlerFunc(func(ctx context.Context, req *http.Request) (resp *http.Response, err error) {
		return nil, nil
	})
}
