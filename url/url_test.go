package url_test

import (
	"bytes"
	"fmt"
	"net/http"
	neturl "net/url"
	"testing"

	"reflect"

	"github.com/delicb/cliware"
	"github.com/delicb/cliware-middlewares/url"
)

type testData struct {
	Input string
	URL   *neturl.URL
}

func diffURL(first *neturl.URL, second *neturl.URL) string {
	var b bytes.Buffer
	if first.Scheme != second.Scheme {
		b.WriteString(fmt.Sprintf("Scheme: %q != %q\n", first.Scheme, second.Scheme))
	}
	if first.Host != second.Host {
		b.WriteString(fmt.Sprintf("Host: %q != %q\n", first.Host, second.Host))
	}
	if first.Path != second.Path {
		b.WriteString(fmt.Sprintf("Path: %q != %q\n", first.Path, second.Path))
	}
	return b.String()
}

func testURLMiddleware(t *testing.T, data []testData, factory func(r *http.Request, input testData) cliware.Middleware) {
	for _, d := range data {
		req := cliware.EmptyRequest()
		handler := createHandler()
		factory(req, d).Exec(handler).Handle(req)
		if !reflect.DeepEqual(req.URL, d.URL) {
			t.Errorf("URL did not match. Diff: \n\t%s", diffURL(d.URL, req.URL))
			// t.Errorf("URL did not match. Got: \"%s\", expected: \"%s\".", d.URL, req.URL)
		}
	}
}

func createHandler() cliware.Handler {
	return cliware.HandlerFunc(func(req *http.Request) (resp *http.Response, err error) {
		return nil, nil
	})
}

func TestURL(t *testing.T) {
	data := []testData{
		{"https://bojan.delic.rs:8443/path?query=value", &neturl.URL{
			Scheme:   "https",
			Host:     "bojan.delic.rs:8443",
			Path:     "/path",
			RawQuery: "query=value",
		}},
		{"bojan.delic.rs", &neturl.URL{Scheme: "https", Path: "bojan.delic.rs"}},
		{"/path", &neturl.URL{Scheme: "https", Path: "/path"}},
	}
	testURLMiddleware(t, data, func(r *http.Request, d testData) cliware.Middleware {
		return url.URL(d.Input)
	})
}

func TestBaseURL(t *testing.T) {
	data := []testData{
		{"https://bojan.delic.rs:1234/path?query=value", &neturl.URL{Scheme: "https", Host: "bojan.delic.rs:1234"}},
		{"http://localhost/path?query=1", &neturl.URL{Scheme: "http", Host: "localhost"}},
	}
	testURLMiddleware(t, data, func(r *http.Request, d testData) cliware.Middleware {
		return url.BaseURL(d.Input)
	})
}

func TestPath(t *testing.T) {
	data := []testData{
		{"/foobar", &neturl.URL{Path: "/foobar"}},
	}
	testURLMiddleware(t, data, func(r *http.Request, d testData) cliware.Middleware {
		return url.Path(d.Input)
	})
}

func TestAddPath(t *testing.T) {
	data := []testData{
		{"/additional", &neturl.URL{Path: "/base/additional"}},
		{"/", &neturl.URL{Path: "/base"}},
		{"", &neturl.URL{Path: "/base"}},
	}
	testURLMiddleware(t, data, func(r *http.Request, d testData) cliware.Middleware {
		r.URL.Path = "/base"
		return url.AddPath(d.Input)
	})
}

func TestPathPrefix(t *testing.T) {
	data := []testData{
		{"/prefix", &neturl.URL{Path: "/prefix/rest"}},
		{"/", &neturl.URL{Path: "/rest"}},
		{"", &neturl.URL{Path: "/rest"}},
	}
	testURLMiddleware(t, data, func(r *http.Request, d testData) cliware.Middleware {
		r.URL.Path = "/rest"
		return url.PathPrefix(d.Input)
	})
}

func TestParam(t *testing.T) {
	for _, data := range []struct {
		InitialPath string
		ParamKey    string
		ParamValue  string
		ResultPath  string
	}{
		{
			InitialPath: "/:parameter",
			ParamKey:    "parameter",
			ParamValue:  "value",
			ResultPath:  "/value",
		},
		{
			InitialPath: "/:parameter",
			ParamKey:    "parameter",
			ParamValue:  "",
			ResultPath:  "/",
		},
		{
			InitialPath: ":parameter",
			ParamKey:    "missing_key",
			ParamValue:  "",
			ResultPath:  ":parameter",
		},
	} {
		req := cliware.EmptyRequest()
		req.URL.Path = data.InitialPath
		handler := createHandler()
		url.Param(data.ParamKey, data.ParamValue).Exec(handler).Handle(req)

		if req.URL.Path != data.ResultPath {
			t.Errorf("Got wrong path. Got: %s, expected: %s.", req.URL.Path, data.ResultPath)
		}
	}
}

func TestParams(t *testing.T) {
	for _, data := range []struct {
		InitialPath string
		Params      map[string]string
		ResultPath  string
	}{
		{
			InitialPath: "/:parameter",
			Params: map[string]string{
				"parameter": "value",
			},
			ResultPath: "/value",
		},
		{
			InitialPath: "/:param1/:param2",
			Params: map[string]string{
				"param1": "value1",
				"param2": "value2",
			},
			ResultPath: "/value1/value2",
		},
		{
			InitialPath: "/:param1/:param2",
			Params: map[string]string{
				"param1": "value1",
			},
			ResultPath: "/value1/:param2",
		},
	} {
		req := cliware.EmptyRequest()
		req.URL.Path = data.InitialPath
		handler := createHandler()
		url.Params(data.Params).Exec(handler).Handle(req)

		if req.URL.Path != data.ResultPath {
			t.Errorf("Got wrong path. Got: %s, expected: %s.", req.URL.Path, data.ResultPath)
		}
	}
}
