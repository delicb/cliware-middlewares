package responsebody // import "go.delic.rs/cliware-middlewares/responsebody"

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	c "go.delic.rs/cliware"
)

// JSON decodes response body from JSON format into provided interface.
func JSON(data interface{}) c.Middleware {
	return c.ResponseProcessor(func(resp *http.Response, err error) error {
		// TODO: Should we check for Content-Type header here?
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		rawData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return json.Unmarshal(rawData, data)
	})
}

// String reads response body, converts it to string and writes it to provided
// string pointer.
func String(data *string) c.Middleware {
	return c.ResponseProcessor(func(resp *http.Response, err error) error {
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		rawData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		*data = string(rawData)
		return nil
	})
}
