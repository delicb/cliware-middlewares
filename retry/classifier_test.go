package retry_test

import (
	"testing"

	"net/http"

	"github.com/pkg/errors"
	"go.delic.rs/cliware-middlewares/retry"
)

func TestAnyError(t *testing.T) {
	if retry.AnyErrorClassifier(nil, nil) {
		t.Error("AnyErrorClassifier returned true for nil error.")
	}
	if !retry.AnyErrorClassifier(nil, errors.New("some error")) {
		t.Error("AnyErrorClassifier returned false for non-nil error.")
	}
}

func TestOn500PlusClassifier(t *testing.T) {
	for resp, expected := range map[*http.Response]bool{
		&http.Response{StatusCode: 200}: false,
		&http.Response{StatusCode: 300}: false,
		&http.Response{StatusCode: 400}: false,
		&http.Response{StatusCode: 404}: false,
		&http.Response{StatusCode: 500}: true,
		&http.Response{StatusCode: 501}: true,
		&http.Response{StatusCode: 502}: true,
		&http.Response{StatusCode: 503}: true,
	} {
		if retry.On500PlusClassifier(resp, nil) != expected {
			t.Errorf("On500PlusClassifier wrong value. Expected: %t", expected)
		}
	}
}

var (
	trueClassifier = retry.Classifier(func(resp *http.Response, err error) bool {
		return true
	})
	falseClassifier = retry.Classifier(func(resp *http.Response, err error) bool {
		return false
	})
)

func TestOrClassifier(t *testing.T) {
	for _, data := range []struct {
		Classifiers []retry.Classifier
		Result      bool
	}{
		{
			Classifiers: []retry.Classifier{trueClassifier, trueClassifier},
			Result:      true,
		},
		{
			Classifiers: []retry.Classifier{trueClassifier, falseClassifier},
			Result:      true,
		},
		{
			Classifiers: []retry.Classifier{falseClassifier, falseClassifier},
			Result:      false,
		},
		{
			Classifiers: []retry.Classifier{},
			Result: false,
		},
	} {
		res := retry.OrClassifier(data.Classifiers...)(nil, nil)
		if res != data.Result {
			t.Errorf("OrClassifier returned wrong value. Got: %t, expected: %t.", res, data.Result)
		}
	}
}

func TestAndClassifier(t *testing.T) {
	for _, data := range []struct{
		Classifiers []retry.Classifier
		Result bool
	}{
		{
			Classifiers: []retry.Classifier{trueClassifier, trueClassifier},
			Result: true,
		},
		{
			Classifiers: []retry.Classifier{trueClassifier, falseClassifier},
			Result: false,
		},
		{
			Classifiers: []retry.Classifier{falseClassifier, falseClassifier},
			Result: false,
		},
		{
			Classifiers: []retry.Classifier{},
			Result: false,
		},
	}{
		res := retry.AndClassifier(data.Classifiers...)(nil, nil)
		if res != data.Result {
			t.Errorf("AndClassifier returned wrong value. Got: %t, expected: %t.", res, data.Result)
		}
	}
}
