package cookies_test

import (
	"context"
	"net/http"
	"testing"

	"fmt"
	"reflect"

	"go.delic.rs/cliware"
	"go.delic.rs/cliware-middlewares/cookies"
)

func cookieInSlice(a *http.Cookie, list []*http.Cookie) bool {
	for _, b := range list {
		if reflect.DeepEqual(a, b) {
			return true
		}
	}
	return false
}

func TestAdd(t *testing.T) {
	for _, data := range []struct {
		Cookie   *http.Cookie
		Existing []*http.Cookie
	}{
		{
			Cookie: &http.Cookie{
				Name:  "mycookie",
				Value: "value",
			},
			Existing: []*http.Cookie{},
		},
		{
			Cookie: &http.Cookie{
				Name:  "some-cookie",
				Value: "golang",
			},
			Existing: []*http.Cookie{
				{
					Name:  "existing",
					Value: "should be there...",
				},
			},
		},
	} {
		m := cookies.Add(data.Cookie)
		chain := cliware.NewChain(m)
		req := cliware.EmptyRequest()
		for _, existing := range data.Existing {
			req.AddCookie(existing)
		}
		_, err := chain.Exec(createHandler()).Handle(nil, req)
		if err != nil {
			t.Error("Handle returned error: ", err)
		}

		for _, c := range data.Existing {
			if !cookieInSlice(c, req.Cookies()) {
				fmt.Printf("\tLooking for cookie: %v in %v\n", c, req.Cookies())
				t.Errorf("Not found preset cookie in request: %v.", c)
			}
		}
		cookie, err := req.Cookie(data.Cookie.Name)
		if err != nil {
			t.Fatal("Getting cookie by name failed: ", err)
		}

		if cookie.Value != data.Cookie.Value {
			t.Errorf("Wrong cookie value. Got: %s, expected: %s.", cookie.Value, data.Cookie.Value)
		}

	}
}

func TestSet(t *testing.T) {
	for _, data := range []struct {
		Cookie   *http.Cookie
		Existing []*http.Cookie
	}{
		{
			Cookie: &http.Cookie{
				Name:  "mycookie",
				Value: "value",
			},
			Existing: []*http.Cookie{},
		},
		{
			Cookie: &http.Cookie{
				Name:  "some-cookie",
				Value: "golang",
			},
			Existing: []*http.Cookie{
				{
					Name:  "existing",
					Value: "should be there...",
				},
			},
		},
	} {
		m := cookies.Set(data.Cookie.Name, data.Cookie.Value)
		chain := cliware.NewChain(m)
		req := cliware.EmptyRequest()
		for _, existing := range data.Existing {
			req.AddCookie(existing)
		}
		_, err := chain.Exec(createHandler()).Handle(nil, req)
		if err != nil {
			t.Error("Handle returned error: ", err)
		}

		for _, c := range data.Existing {
			if !cookieInSlice(c, req.Cookies()) {
				fmt.Printf("\tLooking for cookie: %v in %v\n", c, req.Cookies())
				t.Errorf("Not found preset cookie in request: %v.", c)
			}
		}
		cookie, err := req.Cookie(data.Cookie.Name)
		if err != nil {
			t.Fatal("Getting cookie by name failed: ", err)
		}

		if cookie.Value != data.Cookie.Value {
			t.Errorf("Wrong cookie value. Got: %s, expected: %s.", cookie.Value, data.Cookie.Value)
		}

	}
}

func TestDelAll(t *testing.T) {
	for _, existing := range [][]*http.Cookie{
		[]*http.Cookie{},
		[]*http.Cookie{&http.Cookie{}},
		[]*http.Cookie{{Name: "some-cookie", Value: "cookie value"}, {Name: "other-cookie", Value: "other value"}},
	} {
		m := cookies.DelAll()
		chain := cliware.NewChain(m)
		req := cliware.EmptyRequest()
		for _, c := range existing {
			req.AddCookie(c)
		}
		_, err := chain.Exec(createHandler()).Handle(nil, req)
		if err != nil {
			t.Error("Handle returned error: ", err)
		}
		if len(req.Cookies()) > 0 {
			t.Errorf("Expected no cookies on request, found: %v.", req.Cookies())
		}
	}
}

func TestSetMap(t *testing.T) {
	for _, data := range []struct {
		ToAdd    map[string]string
		Existing []*http.Cookie
	}{
		{
			ToAdd: map[string]string{
				"cookie": "value",
			},
			Existing: []*http.Cookie{},
		},
		{
			ToAdd:    map[string]string{},
			Existing: []*http.Cookie{},
		},
		{
			ToAdd: map[string]string{
				"cookie":       "value",
				"other-cookie": "other-value",
			},
			Existing: []*http.Cookie{
				{
					Name:  "existing",
					Value: "existing value",
				},
			},
		},
	} {
		m := cookies.SetMap(data.ToAdd)
		chain := cliware.NewChain(m)
		req := cliware.EmptyRequest()
		for _, c := range data.Existing {
			req.AddCookie(c)
		}
		_, err := chain.Exec(createHandler()).Handle(nil, req)
		if err != nil {
			t.Error("Handle returned error: ", err)
		}

		for _, c := range data.Existing {
			if !cookieInSlice(c, req.Cookies()) {
				t.Errorf("Expected cookie: %v to be in request, but not found:", c)
			}
		}
		for k, v := range data.ToAdd {
			c, err := req.Cookie(k)
			if err != nil {
				t.Errorf("Getting cookie %s returned error: %s", k, err)
			}
			if c.Value != v {
				t.Errorf("Wrong cookie value. Got: %s, expected: %s.", c.Value, v)
			}
		}
	}
}

func TestAddMultiple(t *testing.T) {
	for _, data := range []struct {
		ToAdd    []*http.Cookie
		Existing []*http.Cookie
	}{
		{
			ToAdd:    []*http.Cookie{},
			Existing: []*http.Cookie{},
		},
		{
			ToAdd:    []*http.Cookie{{Name: "cookie", Value: "some value"}},
			Existing: []*http.Cookie{},
		},
		{
			ToAdd:    []*http.Cookie{{Name: "cookie", Value: "some value"}},
			Existing: []*http.Cookie{{Name: "existing", Value: "existing value"}},
		},
	} {
		m := cookies.AddMultiple(data.ToAdd)
		chain := cliware.NewChain(m)
		req := cliware.EmptyRequest()
		for _, c := range data.Existing {
			req.AddCookie(c)
		}
		_, err := chain.Exec(createHandler()).Handle(nil, req)
		if err != nil {
			t.Error("Handle returned error: ", err)
		}

		for _, c := range data.Existing {
			if !cookieInSlice(c, req.Cookies()) {
				t.Errorf("Expected cookie: %v to be in request, but not found:", c)
			}
		}
		for _, c := range data.ToAdd {
			if !cookieInSlice(c, req.Cookies()) {
				t.Errorf("Expected cookie: %v to be in request, but not found:", c)
			}
		}
	}
}

func createHandler() cliware.Handler {
	return cliware.HandlerFunc(func(ctx context.Context, req *http.Request) (resp *http.Response, err error) {
		return nil, nil
	})
}
