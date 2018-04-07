package compute

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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
	// All connections to the pool are considered.
	LoadBalanceMethodPredictiveNode = "PREDICTIVE_NODE"

	// LoadBalanceMethodPredictiveMember indicates that requests will be directed to the pool node that is predicted to have the smallest number of active connections over time.
	// Only connections to the pool as a member of the current pool are considered.
	LoadBalanceMethodPredictiveMember = "PREDICTIVE_MEMBER"
)

// VIPPool represents a VIP pool.
type VIPPool struct {
	ID                string            `json:"id"`
	Name              string            `json:"name"`
	Description       string            `json:"description"`
	LoadBalanceMethod string            `json:"loadBalanceMethod"`
	HealthMonitors    []EntityReference `json:"healthMonitor"`
	ServiceDownAction string            `json:"serviceDownAction"`
	SlowRampTime      int               `json:"slowRampTime"`
	State             string            `json:"state"`
	NetworkDomainID   string            `json:"networkDomainID"`
	DataCenterID      string            `json:"datacenterId"`
	CreateTime        string            `json:"createTime"`
}

// GetID returns the pool's Id.
func (pool *VIPPool) GetID() string {
	return pool.ID
}

// GetResourceType returns the pool's resource type.
func (pool *VIPPool) GetResourceType() ResourceType {
	return ResourceTypeVIPPool
}

// GetName returns the pool's name.
func (pool *VIPPool) GetName() string {
	return pool.Name
}

// GetState returns the pool's current state.
func (pool *VIPPool) GetState() string {
	return pool.State
}

// IsDeleted determines whether the pool has been deleted (is nil).
func (pool *VIPPool) IsDeleted() bool {
	return pool == nil
}

var _ Resource = &VIPPool{}

// ToEntityReference creates an EntityReference representing the VIPNode.
func (pool *VIPPool) ToEntityReference() EntityReference {
	return EntityReference{
		ID:   pool.ID,
		Name: pool.Name,
	}
}

var _ NamedEntity = &VIPNode{}

// VIPPools represents a page of VIPPool results.
type VIPPools struct {
	Items []VIPPool

	PagedResult
}

// NewVIPPoolConfiguration represents the configuration for a new VIP pool.
type NewVIPPoolConfiguration struct {
	// The VIP pool name.
	Name string `json:"name"`

	// The VIP pool description.
	Description string `json:"description"`

	// The load-balancing method used for pools in the pool.
	LoadBalanceMethod string `json:"loadBalanceMethod"`

	// The Id of the pool's associated health monitors (if any).
	// Up to 2 health monitors can be specified per pool.
	HealthMonitorIDs []string `json:"healthMonitorId,omitempty"`

	// The action performed when a pool in the pool is down.
	ServiceDownAction string `json:"serviceDownAction"`

	// The time, in seconds, over which the the pool will ramp new pools up to their full request rate.
	SlowRampTime int `json:"slowRampTime"`

	// The Id of the network domain where the pool is located.
	NetworkDomainID string `json:"networkDomainId"`

	// The Id of the data centre where the pool is located.
	DatacenterID string `json:"datacenterId,omitempty"`
}

// EditVIPPoolConfiguration represents the request body when editing a VIP pool.
type EditVIPPoolConfiguration struct {
	// The VIP pool Id.
	ID string `json:"id"`

	// The VIP pool description.
	Description *string `json:"description,omitempty"`

	// The load-balancing method used for pools in the pool.
	LoadBalanceMethod *string `json:"loadBalanceMethod"`

	// The Id of the pool's associated health monitors (if any).
	// Up to 2 health monitors can be specified per pool.
	HealthMonitorIDs *[]string `json:"healthMonitorId,omitempty"`

	// The action performed when a pool in the pool is down.
	ServiceDownAction *string `json:"serviceDownAction"`

	// The time, in seconds, over which the the pool will ramp new pools up to their full request rate.
	SlowRampTime *int `json:"slowRampTime"`
}

// Request body for deleting a VIP pool.
type deleteVIPPool struct {
	// The VIP pool ID.
	ID string `json:"id"`
}

// ListVIPPoolsInNetworkDomain retrieves a list of all VIP pools in the specified network domain.
func (client *Client) ListVIPPoolsInNetworkDomain(networkDomainID string, paging *Paging) (pools *VIPPools, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return nil, err
	}

	requestURI := fmt.Sprintf("%s/networkDomainVip/pool?networkDomainId=%s&%s",
		url.QueryEscape(organizationID),
		url.QueryEscape(networkDomainID),
		paging.EnsurePaging().toQueryParameters(),
	)
	request, err := client.newRequestV26(requestURI, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		var apiResponse *APIResponseV2

		apiResponse, err = readAPIResponseAsJSON(responseBody, statusCode)
		if err != nil {
			return nil, err
		}

		return nil, apiResponse.ToError("Request to list VIP pools in network domain '%s' failed with status code %d (%s): %s", networkDomainID, statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	pools = &VIPPools{}
	err = json.Unmarshal(responseBody, pools)
	if err != nil {
		return nil, err
	}

	return pools, nil
}

// GetVIPPool retrieves the VIP pool with the specified Id.
// Returns nil if no VIP pool is found with the specified Id.
func (client *Client) GetVIPPool(id string) (pool *VIPPool, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return nil, err
	}

	requestURI := fmt.Sprintf("%s/networkDomainVip/pool/%s",
		url.QueryEscape(organizationID),
		url.QueryEscape(id),
	)
	request, err := client.newRequestV26(requestURI, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		var apiResponse *APIResponseV2

		apiResponse, err = readAPIResponseAsJSON(responseBody, statusCode)
		if err != nil {
			return nil, err
		}

		if apiResponse.ResponseCode == ResponseCodeResourceNotFound {
			return nil, nil // Not an error, but was not found.
		}

		return nil, apiResponse.ToError("Request to retrieve VIP pool with Id '%s' failed with status code %d (%s): %s", id, statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	pool = &VIPPool{}
	err = json.Unmarshal(responseBody, pool)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

// CreateVIPPool creates a new VIP pool.
// Returns the Id of the new pool.
func (client *Client) CreateVIPPool(poolConfiguration NewVIPPoolConfiguration) (poolID string, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return "", err
	}

	requestURI := fmt.Sprintf("%s/networkDomainVip/createPool",
		url.QueryEscape(organizationID),
	)
	request, err := client.newRequestV26(requestURI, http.MethodPost, &poolConfiguration)
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

// EditVIPPool updates an existing VIP pool.
func (client *Client) EditVIPPool(id string, poolConfiguration EditVIPPoolConfiguration) error {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return err
	}

	editPoolConfiguration := &poolConfiguration
	editPoolConfiguration.ID = id

	requestURI := fmt.Sprintf("%s/networkDomainVip/editPool",
		url.QueryEscape(organizationID),
	)
	request, err := client.newRequestV26(requestURI, http.MethodPost, editPoolConfiguration)
	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return err
	}

	apiResponse, err := readAPIResponseAsJSON(responseBody, statusCode)
	if err != nil {
		return err
	}

	if statusCode != http.StatusOK {
		return apiResponse.ToError("Request to edit VIP pool '%s' failed with status code %d (%s): %s", poolConfiguration.ID, statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	return nil
}

// DeleteVIPPool deletes an existing VIP pool.
// Returns an error if the operation was not successful.
func (client *Client) DeleteVIPPool(id string) (err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return err
	}

	requestURI := fmt.Sprintf("%s/networkDomainVip/deletePool",
		url.QueryEscape(organizationID),
	)
	request, err := client.newRequestV26(requestURI, http.MethodPost, &deleteVIPPool{id})
	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return err
	}

	apiResponse, err := readAPIResponseAsJSON(responseBody, statusCode)
	if err != nil {
		return err
	}

	if apiResponse.ResponseCode != ResponseCodeOK {
		return apiResponse.ToError("Request to delete VIP pool '%s' failed with unexpected status code %d (%s): %s", id, statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	return nil
}
