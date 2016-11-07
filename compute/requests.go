package compute

import (
	"bytes"
	"io"
	"net/http"
)

// Build an httpRequestCloner for the specified request.
func requestCloner(request *http.Request) (*httpRequestCloner, error) {
	requestBody, err := getRequestBody(request)
	if err != nil {
		return nil, err
	}

	cloner := &httpRequestCloner{
		request:           request,
		cachedRequestBody: requestBody,
	}

	return cloner, nil
}

type httpRequestCloner struct {
	request           *http.Request
	cachedRequestBody []byte
}

func (cloner *httpRequestCloner) Clone() (clonedRequest *http.Request, err error) {
	clonedRequest, err = http.NewRequest(
		cloner.request.Method,
		cloner.request.URL.String(),
		cloner.CloneRequestBody(),
	)
	if err != nil {
		return
	}

	// Copy headers
	for headerName, headerValue := range cloner.request.Header {
		clonedRequest.Header[headerName] = headerValue
	}

	return
}
func (cloner *httpRequestCloner) CloneRequestBody() io.Reader {
	if cloner.cachedRequestBody == nil {
		return nil
	}

	return bytes.NewReader(cloner.cachedRequestBody)
}
