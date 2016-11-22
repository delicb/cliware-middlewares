package retry

import "net/http"

// Classifier is function that determines if request should be retried.
// Boolean return value indicates if request should be repeated or not.
type Classifier func(resp *http.Response, err error) (shouldRetry bool)

// AnyErrorClassifier is classifier that indicates that all requests that return
// errors should be retried.
func AnyErrorClassifier(resp *http.Response, err error) bool {
	return err != nil
}

// On500PlusClassifier is classifier that indicates that all requests whose
// response code is 500 or higher should be repeated.
func On500PlusClassifier(resp *http.Response, err error) bool {
	return err != nil || resp.StatusCode >= 500
}

// OrClassifier is classifier that combines other classifiers. Returned classifier
// will indicate that request should be repeated if any or provided classifiers
// returned true.
func OrClassifier(classifiers ...Classifier) Classifier {
	return Classifier(func(resp *http.Response, err error) bool {
		// if there are not classifiers provided, default is false
		if len(classifiers) == 0 {
			return false
		}
		for _, c := range classifiers {
			if c(resp, err) {
				return true
			}
		}
		// if we got this far, all classifiers returned false
		return false

	})
}

// AndClassifier is classifier that combines other classifiers. Returned
// classifier will indicate that request should be repeated if all of provided
// classifiers returned true.
func AndClassifier(classifiers ...Classifier) Classifier {
	return Classifier(func(resp *http.Response, err error) bool {
		// if there are not classifiers provided, false is default
		if len(classifiers) == 0 {
			return false
		}

		for _, c := range classifiers {
			if !c(resp, err) {
				return false
			}
		}
		// if we got this far, all classifiers returned true
		return true
	})
}

// ErrorOr500Plus is classifier that combines AnyError and On500Plus classifiers.
// This means that classifier will classify any response that returned error or
// status code >= 500 to be retried.
var ErrorOr500Plus = OrClassifier(AnyErrorClassifier, On500PlusClassifier)
