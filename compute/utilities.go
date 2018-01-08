package compute

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

func stringToPtr(value string) *string {
	return &value
}

func intToPtr(value int) *int {
	return &value
}

// Get the request body, replacing it with a copy of the original
func getRequestBody(request *http.Request) (requestBody []byte, err error) {
	if request.Body != nil {
		defer request.Body.Close()

		requestBody, err = ioutil.ReadAll(request.Body)
		if err != nil {
			err = errors.Wrapf(err, "unexpected error reading request body")

			return
		}

		var requestBodyReader io.Reader = bytes.NewReader(requestBody)
		requestBodyReadCloser, ok := requestBodyReader.(io.ReadCloser)
		if !ok {
			requestBodyReadCloser = ioutil.NopCloser(requestBodyReader)
		}
		request.Body = requestBodyReadCloser
	}

	return
}
