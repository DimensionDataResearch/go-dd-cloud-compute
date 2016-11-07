package compute

import (
	"log"
	"net/http"
	"time"

	"github.com/DimensionDataResearch/go-dd-cloud-compute/compute/requests"
)

// executeRequestWithRetry performs the specified request, retrying if it fails at the transport level.
func (client *Client) executeRequestWithRetry(action string, builder requests.RequestBuilder) (responseBody []byte, statusCode int, err error) {
	clientContext := &clientRequestContext{
		Action: action,
		Client: client,
	}
	clientContext.OrganizationID, err = client.getOrganizationID()
	if err != nil {
		return
	}

	for {
		var request *http.Request
		request, err = builder(clientContext)
		if err != nil {
			return
		}
		
		responseBody, statusCode, err = client.executeRequest(request)
		if err == nil {
			return
		}

		clientContext.RetryCount++
		if clientContext.GetRemainingRetryCount() == 0 {
			log.Printf("Exceeded the maximum number of retries (%d) for '%s' request.",
				client.maxRetryCount,
				clientContext.Action,
			)

			return
		}

		if client.IsExtendedLoggingEnabled() {
			log.Printf("Waiting %d seconds before retrying '%s' request to '%s' (%d retries remaining)...",
				client.retryDelay/time.Second,
				request.Method,
				request.URL.String(),
				clientContext.GetRemainingRetryCount(),
			)
		}
	}
}
