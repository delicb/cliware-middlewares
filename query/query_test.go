package query_test

import (
	"net/http"
	"testing"

	neturl "net/url"

	"reflect"

	"github.com/delicb/cliware"
	"github.com/delicb/cliware-middlewares/query"
)

func TestSet(t *testing.T) {
	for _, data := range []struct {
		InitialURL     string
		ParamName      string
		ParamValue     string
		ResultingQuery string
	}{
		{
			InitialURL:     "https://bojan.delic.rs:8443/path?query=value",
			ParamName:      "query",
			ParamValue:     "other_value",
			ResultingQuery: "query=other_value",
		},
		{
			InitialURL:     "https://bojan.delic.rs:8443/path?query=value",
			ParamName:      "query2",
			ParamValue:     "value2",
			ResultingQuery: "query=value&query2=value2",
		},
	} {
		m := query.Set(data.ParamName, data.ParamValue)
		req := cliware.EmptyRequest()
		parsedURL, err := neturl.Parse(data.InitialURL)
		if err != nil {
			t.Fatal("Input data not valid. URL parsing failed: ", err)
		}
		req.URL = parsedURL
		handler := createHandler()
		m.Exec(handler).Handle(req)

		q, err := neturl.ParseQuery(data.ResultingQuery)
		if err != nil {
			t.Fatal("Invalid test data, failed to parse query params: ", data.ResultingQuery)
		}

		if !reflect.DeepEqual(req.URL.Query(), q) {
			t.Errorf("Wrong query parameters. Got: %s, expected: %s.", req.URL.RawQuery, data.ResultingQuery)
		}
	}
}

func TestAdd(t *testing.T) {
	for _, data := range []struct {
		InitialURL     string
		ParamName      string
		ParamValue     string
		ResultingQuery string
	}{
		{
			InitialURL:     "https://bojan.delic.rs:8443/path?query=value",
			ParamName:      "query",
			ParamValue:     "other_value",
			ResultingQuery: "query=value&query=other_value",
		},
		{
			InitialURL:     "https://bojan.delic.rs:8443/path?query=value",
			ParamName:      "query2",
			ParamValue:     "value2",
			ResultingQuery: "query=value&query2=value2",
		},
	} {
		m := query.Add(data.ParamName, data.ParamValue)
		req := cliware.EmptyRequest()
		parsedURL, err := neturl.Parse(data.InitialURL)
		if err != nil {
			t.Fatal("Input data not valid. URL parsing failed: ", err)
		}
		req.URL = parsedURL
		handler := createHandler()
		m.Exec(handler).Handle(req)

		q, err := neturl.ParseQuery(data.ResultingQuery)
		if err != nil {
			t.Fatal("Invalid test data, failed to parse query params: ", data.ResultingQuery)
		}

		if !reflect.DeepEqual(req.URL.Query(), q) {
			t.Errorf("Wrong query parameters. Got: %s, expected: %s.", req.URL.RawQuery, data.ResultingQuery)
		}
	}
}

func TestDel(t *testing.T) {
	for _, data := range []struct {
		InitialURL     string
		ParamName      string
		ResultingQuery string
	}{
		{
			InitialURL:     "https://bojan.delic.rs:8443/path?query=value",
			ParamName:      "query",
			ResultingQuery: "",
		},
		{
			InitialURL:     "https://bojan.delic.rs:8443/path?query=value&other_query=other_value",
			ParamName:      "query",
			ResultingQuery: "other_query=other_value",
		},
	} {
		m := query.Del(data.ParamName)
		req := cliware.EmptyRequest()
		parsedURL, err := neturl.Parse(data.InitialURL)
		if err != nil {
			t.Fatal("Input data not valid. URL parsing failed: ", err)
		}
		req.URL = parsedURL
		handler := createHandler()
		m.Exec(handler).Handle(req)

		q, err := neturl.ParseQuery(data.ResultingQuery)
		if err != nil {
			t.Fatal("Invalid test data, failed to parse query params: ", data.ResultingQuery)
		}

		if !reflect.DeepEqual(req.URL.Query(), q) {
			t.Errorf("Wrong query parameters. Got: %s, expected: %s.", req.URL.RawQuery, data.ResultingQuery)
		}
	}
}

func TestDelAll(t *testing.T) {
	for _, rawURL := range []string{
		"https://bojan.delic.rs:8443/path?query=value",
		"https://bojan.delic.rs:8443/path?query=value&other_query=other_value",
	} {
		m := query.DelAll()
		req := cliware.EmptyRequest()
		parsedURL, err := neturl.Parse(rawURL)
		if err != nil {
			t.Fatal("Input data not valid. URL parsing failed: ", err)
		}
		req.URL = parsedURL
		handler := createHandler()
		m.Exec(handler).Handle(req)

		q := neturl.Values{}

		if !reflect.DeepEqual(req.URL.Query(), q) {
			t.Errorf("Wrong query parameters. Got: %s, expected empty.", req.URL.RawQuery)
		}
	}
}

func TestSetMap(t *testing.T) {
	for _, data := range []struct {
		InitialURL     string
		QueryMap       map[string]string
		ResultingQuery string
	}{
		{
			InitialURL: "https://bojan.delic.rs:8443/path?query=value",
			QueryMap: map[string]string{
				"query": "other_value",
			},
			ResultingQuery: "query=other_value",
		},
		{
			InitialURL: "https://bojan.delic.rs:8443/path?query=value&other_query=other_value",
			QueryMap: map[string]string{
				"query":    "other_value",
				"my_query": "my_value",
			},
			ResultingQuery: "query=other_value&other_query=other_value&my_query=my_value",
		},
	} {
		m := query.SetMap(data.QueryMap)
		req := cliware.EmptyRequest()
		parsedURL, err := neturl.Parse(data.InitialURL)
		if err != nil {
			t.Fatal("Input data not valid. URL parsing failed: ", err)
		}
		req.URL = parsedURL
		handler := createHandler()
		m.Exec(handler).Handle(req)

		q, err := neturl.ParseQuery(data.ResultingQuery)
		if err != nil {
			t.Fatal("Invalid test data, failed to parse query params: ", data.ResultingQuery)
		}

		if !reflect.DeepEqual(req.URL.Query(), q) {
			t.Errorf("Wrong query parameters. Got: %s, expected: %s.", req.URL.RawQuery, data.ResultingQuery)
		}
	}
}

func createHandler() cliware.Handler {
	return cliware.HandlerFunc(func(req *http.Request) (resp *http.Response, err error) {
		return nil, nil
	})
}
