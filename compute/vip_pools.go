package compute

import (
	"fmt"
	"net/http"
)

const (
	// ServiceDownActionNone indicates no action will be taken when a pool service is down.
	ServiceDownActionNone = "NONE"

	// ServiceDownActionDrop indicates that a pool service will be dropped when it is down.
	ServiceDownActionDrop = "DROP"

	// ServiceDownActionReselect indicates that a pool service will be reselected when it is down.
	ServiceDownActionReselect = "RESELECT"

	// LoadBalanceMethodRoundRobin indicates that requests will be directed to pool nodes in round-robin fashion.
	LoadBalanceMethodRoundRobin = "ROUND_ROBIN"

	// LoadBalanceMethodLeastConnectionsNode indicates that requests will be directed to the pool node that has the smallest number of active connections at the moment of connection.
	// All connections to the node are considered.
	LoadBalanceMethodLeastConnectionsNode = "LEAST_CONNECTIONS_NODE"

	// LoadBalanceMethodLeastConnectionsMember indicates that requests will be directed to the pool node that has the smallest number of active connections at the moment of connection.
	// Only connections to the node as a member of the current pool are considered.
	LoadBalanceMethodLeastConnectionsMember = "LEAST_CONNECTIONS_MEMBER"

	// LoadBalanceMethodObservedNode indicates that requests will be directed to the pool node that has the smallest number of active connections over time.
	// All connections to the node are considered.
	LoadBalanceMethodObservedNode = "OBSERVED_NODE"

	// LoadBalanceMethodObservedMember indicates that requests will be directed to the pool node that has the smallest number of active connections over time.
	// Only connections to the node as a member of the current pool are considered.
	LoadBalanceMethodObservedMember = "OBSERVED_MEMBER"

	// LoadBalanceMethodPredictiveNode indicates that requests will be directed to the pool node that is predicted to have the smallest number of active connections.
	// All connections to the node are considered.
	LoadBalanceMethodPredictiveNode = "PREDICTIVE_NODE"

	// LoadBalanceMethodPredictiveMember indicates that requests will be directed to the pool node that is predicted to have the smallest number of active connections over time.
	// Only connections to the node as a member of the current pool are considered.
	LoadBalanceMethodPredictiveMember = "PREDICTIVE_MEMBER"
)

// NewVIPPoolConfiguration represents the configuration for a new VIP pool.
type NewVIPPoolConfiguration struct {
	// The VIP node name.
	Name string `json:"name"`

	// The VIP node description.
	Description string `json:"description"`

	// The load-balancing method used for nodes in the pool.
	LoadBalanceMethod string `json:"loadBalanceMethod"`

	// The Id of the node's associated health monitors (if any).
	// Up to 2 health monitors can be specified per pool.
	HealthMonitorIDs []string `json:"healthMonitorId,omitempty"`

	// The action performed when a node in the pool is down.
	ServiceDownAction string `json:"serviceDownAction"`

	// The time, in seconds, over which the the pool will ramp new nodes up to their full request rate.
	SlowRampTime int `json:"slowRampTime"`

	// The Id of the network domain where the node is located.
	NetworkDomainID string `json:"networkDomainId"`

	// The Id of the data centre where the node is located.
	DatacenterID string `json:"datacenterId"`
}

// CreateVIPPool creates a new VIP pool.
// Returns the Id of the new pool.
func (client *Client) CreateVIPPool(poolConfiguration NewVIPPoolConfiguration) (poolID string, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return "", err
	}

	requestURI := fmt.Sprintf("%s/network/networkDomainVip/createPool", organizationID)
	request, err := client.newRequestV22(requestURI, http.MethodPost, &poolConfiguration)
	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return "", err
	}

	apiResponse, err := readAPIResponseAsJSON(responseBody, statusCode)
	if err != nil {
		return "", err
	}

	if apiResponse.ResponseCode != ResponseCodeOK {
		return "", apiResponse.ToError("Request to create VIP pool '%s' failed with status code %d (%s): %s", poolConfiguration.Name, statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	// Expected: "info" { "name": "poolId", "value": "the-Id-of-the-new-pool" }
	poolIDMessage := apiResponse.GetFieldMessage("poolId")
	if poolIDMessage == nil {
		return "", apiResponse.ToError("Received an unexpected response (missing 'poolId') with status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	return *poolIDMessage, nil
}
