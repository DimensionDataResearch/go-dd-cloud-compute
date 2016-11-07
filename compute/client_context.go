package compute

import (
	"github.com/DimensionDataResearch/go-dd-cloud-compute/compute/requests"
)

// The context for an outgoing CloudControl API client request.
type clientRequestContext struct {
	Action     string
	OrganizationID string
	RetryCount int
	LastError  error

	Client *Client
}

var _ requests.RequestContext = &clientRequestContext{}
func (context *clientRequestContext) GetAction() string {
	return context.Action
}
func (context *clientRequestContext) GetOrganizationID() string {
	return context.OrganizationID
}
func (context *clientRequestContext) GetRetryCount() int {
	return context.RetryCount
}
func (context *clientRequestContext) GetRemainingRetryCount() int {
	return context.Client.maxRetryCount - context.RetryCount
}
func (context *clientRequestContext) GetLastError() error {
	return context.LastError
}
