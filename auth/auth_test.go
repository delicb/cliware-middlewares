package auth_test

import (
	"net/http"
	"testing"

	"strings"

	"github.com/delicb/cliware"
	"github.com/delicb/cliware-middlewares/auth"
)

func TestBasic(t *testing.T) {
	m := auth.Basic("username", "password")
	chain := cliware.NewChain(m)
	req := cliware.EmptyRequest()
	chain.Exec(createHandler()).Handle(req)
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

	username, password, ok := req.BasicAuth()
	if !ok {
		t.Error("Unable to read basic auth")
		return
	}

	if username != "username" {
		t.Errorf("Wrong username, expected 'username', got: '%s'", username)
	}

	if password != "password" {
		t.Errorf("Wrong password, expected 'password', got: '%s'", password)
	}
}

func TestBearer(t *testing.T) {
	m := auth.Bearer("token")
	chain := cliware.NewChain(m)
	req := cliware.EmptyRequest()
	chain.Exec(createHandler()).Handle(req)
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

func TestCustom(t *testing.T) {
	header := "OAuth oauth_consumer_key=\"foobar\""
	m := auth.Custom(header)
	chain := cliware.NewChain(m)
	req := cliware.EmptyRequest()
	chain.Exec(createHandler()).Handle(req)
	val, ok := req.Header["Authorization"]
	if !ok {
		t.Error("Expected request to have Authorization header, none found.")
	}
	if len(val) != 1 {
		t.Fatalf("Expected only one Authorization header, found: %d.", len(val))
	}
	authVal := val[0]
	if authVal != header {
		t.Errorf("Wrong value for Authorization header. Got: %s, expected: %s", authVal, header)
	}
}

func createHandler() cliware.Handler {
	return cliware.HandlerFunc(func(req *http.Request) (resp *http.Response, err error) {
		return nil, nil
	})
}
