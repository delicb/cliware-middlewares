package auth_test

import (
	"context"
	"net/http"
	"testing"

	"strings"

	"go.delic.rs/cliware"
	"go.delic.rs/cliware-middlewares/auth"
)

func TestBasic(t *testing.T) {
	m := auth.Basic("username", "password")
	chain := cliware.NewChain(m)
	req := cliware.EmptyRequest()
	chain.Exec(createHandler()).Handle(nil, req)
	val, ok := req.Header["Authorization"]
	if !ok {
		t.Error("Expected request to have Authorization header, none found.")
	}
	if len(val) != 1 {

		t.Fatalf("Expected only one Authorization header, found %d.", len(val))
	}
	authVal := val[0]
	if !strings.HasPrefix(authVal, "Basic ") {
		t.Errorf("Expected Authorization header to have value that starts with 'Basic', but got: %s", authVal)
	}
}

func TestBearer(t *testing.T) {
	m := auth.Bearer("token")
	chain := cliware.NewChain(m)
	req := cliware.EmptyRequest()
	chain.Exec(createHandler()).Handle(nil, req)
	val, ok := req.Header["Authorization"]
	if !ok {
		t.Error("Expected request to have Authorization header, none found.")
	}
	if len(val) != 1 {
		t.Fatalf("Expected only one Authorization header, found: %d.", len(val))
	}
	authVal := val[0]
	expected := "Bearer token"
	if authVal != expected {
		t.Errorf("Wrong value for Authorization header. Got: %s, expected: %s", authVal, expected)
	}
}

func createHandler() cliware.Handler {
	return cliware.HandlerFunc(func(ctx context.Context, req *http.Request) (resp *http.Response, err error) {
		return nil, nil
	})
}
