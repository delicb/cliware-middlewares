package cookies_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/delicb/cliware"
	"github.com/delicb/cliware-middlewares/cookies"
)

func TestAdd(t *testing.T) {

	m := cookies.Add(&http.Cookie{
		Name:   "mycookie",
		Domain: "loclahost",
		Value:  "cookie value",
	})
	m2 := cookies.Add(&http.Cookie{
		Name:   "mycookie",
		Domain: "loclahost",
		Value:  "cookie value",
	})
	chain := cliware.NewChain(m, m2)
	req := cliware.EmptyRequest()
	handler := createHandler()
	chain.Exec(handler).Handle(nil, req)
	t.Log(req.Cookies())
	t.Log(req.Cookie("mycookie"))
}

func createHandler() cliware.Handler {
	return cliware.HandlerFunc(func(ctx context.Context, req *http.Request) (resp *http.Response, err error) {
		return nil, nil
	})
}
