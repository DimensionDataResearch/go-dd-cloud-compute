package requests

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// CreateSnapshot creates a snapshot of the specified request.
func CreateSnapshot(request *http.Request) (*Snapshot, error) {
	requestBody, err := CacheBody(request)
	if err != nil {
		return nil, err
	}

	snapshot := &Snapshot{
		requestURL:        request.URL.String(),
		requestMethod:     request.Method,
		requestHeaders:    make(map[string][]string, len(request.Header)),
		cachedRequestBody: requestBody,
	}
	for headerName, headerValue := range request.Header {
		snapshot.requestHeaders[headerName] = headerValue
	}

	return snapshot, nil
}

// CreateSnapshotAndClose creates a snapshot of the specified request, then destroys the original request.
func CreateSnapshotAndClose(request *http.Request) (snapshot *Snapshot, err error) {
	snapshot, err = CreateSnapshot(request)
	if err != nil {
		return
	}

	if request.Body != nil {
		request.Body.Close()
	}

	return
}

// Snapshot represents a snapshot of request state, and uses that snapshot to create clones of the request as required.
type Snapshot struct {
	requestURL        string
	requestMethod     string
	requestHeaders    map[string][]string
	cachedRequestBody []byte
}

// Copy creates a copy of the request represented by the snapshot.
func (snapshot *Snapshot) Copy() (clonedRequest *http.Request, err error) {
	clonedRequest, err = http.NewRequest(
		snapshot.requestMethod,
		snapshot.requestURL,
		snapshot.GetCachedRequestBodyReader(),
	)
	if err != nil {
		return
	}

	// Copy headers
	for headerName, headerValue := range snapshot.requestHeaders {
		clonedRequest.Header[headerName] = headerValue
	}

	return
}

// GetCachedRequestBody retrieves a copy of the cached request body from the snapshot.
//
// Returns an empty array if the request has no body.
func (snapshot *Snapshot) GetCachedRequestBody() []byte {
	if snapshot.cachedRequestBody == nil {
		return make([]byte, 0)
	}

	return snapshot.cachedRequestBody[:]
}

// GetCachedRequestBodyReader creates a reader over the cached request body from the snapshot.
//
// Returns nil if the request has no body.
func (snapshot *Snapshot) GetCachedRequestBodyReader() io.Reader {
	if snapshot.cachedRequestBody == nil {
		return nil
	}

	return bytes.NewReader(snapshot.cachedRequestBody)
}

// CacheBody retrieves the request body, replacing it with a copy of the original
func CacheBody(request *http.Request) (requestBody []byte, err error) {
	if request.Body != nil {
		defer request.Body.Close()

		requestBody, err = ioutil.ReadAll(request.Body)
		if err != nil {
			err = fmt.Errorf("Unexpected error reading request body: %s", err.Error())
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
