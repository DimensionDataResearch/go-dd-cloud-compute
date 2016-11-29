package compute

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	// VIPNodeStatusEnabled represents a VIP node that is currently enabled.
	VIPNodeStatusEnabled = "ENABLED"

	// VIPNodeStatusDisabled represents a VIP node that is currently disabled.
	VIPNodeStatusDisabled = "DISABLED"

	// VIPNodeStatusForcedOffline represents a VIP node that has been forced offline.
	VIPNodeStatusForcedOffline = "FORCED_OFFLINE"
)

// VIPNodeReference represents a reference to a VIP node.
type VIPNodeReference struct {
	EntityReference

	IPAddress string `json:"ipAddress"`
	Status    string `json:"status"`
}

// VIPNodeHealthMonitor represents a health Monitor to a VIP node.
type VIPNodeHealthMonitor struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// VIPNode represents a VIP node.
type VIPNode struct {
	// The node Id.
	ID string `json:"id"`

	// The node name.
	Name string `json:"name"`

	// The node description.
	Description string `json:"description"`

	// VIPNode's IPv4 address (either IPv4 or IPv6 address must be specified).
	IPv4Address string `json:"ipv4Address,omitempty"`

	// VIPNode's IPv6 address (either IPv4 or IPv6 address must be specified).
	IPv6Address string `json:"ipv6Address,omitempty"`

	// The node status (VIPNodeStatusEnabled, VIPNodeStatusDisabled, or VIPNodeStatusForcedOffline).
	Status string `json:"status"`

	// The Id of the node's associated health monitor (if any).
	HealthMonitor VIPNodeHealthMonitor `json:"healthMonitor,omitempty"`

	// The node's connection limit (must be greater than 0).
	ConnectionLimit int `json:"connectionLimit"`

	// The node's connection rate limit (must be greater than 0).
	ConnectionRateLimit int `json:"connectionRateLimit"`

	// The Id of the network domain where the node is located.
	NetworkDomainID string `json:"networkDomainId"`

	// The Id of the data centre where the node is located.
	DataCenterID string `json:"datacenterId"`

	// The node's creation timestamp.
	CreateTime string `json:"createTime"`

	// The node's current state.
	State string `json:"state"`

	// The node's current progress (if any).
	Progress string `json:"progress"`
}

// GetID returns the node's Id.
func (node *VIPNode) GetID() string {
	return node.ID
}

// GetResourceType returns the node's resource type.
func (node *VIPNode) GetResourceType() ResourceType {
	return ResourceTypeVIPNode
}

// GetName returns the node's name.
func (node *VIPNode) GetName() string {
	return node.Name
}

// GetState returns the node's current state.
func (node *VIPNode) GetState() string {
	return node.State
}

// IsDeleted determines whether the node has been deleted (is nil).
func (node *VIPNode) IsDeleted() bool {
	return node == nil
}

var _ Resource = &VIPNode{}

// ToEntityReference creates an EntityReference representing the VIPNode.
func (node *VIPNode) ToEntityReference() EntityReference {
	return EntityReference{
		ID:   node.ID,
		Name: node.Name,
	}
}

var _ NamedEntity = &VIPNode{}

// VIPNodes represents a page of VIPNode results.
type VIPNodes struct {
	// The current page of node results.
	Items []VIPNode `json:"node"`

	PagedResult
}

// NewVIPNodeConfiguration represents the configuration for a new VIP node.
type NewVIPNodeConfiguration struct {
	// The VIP node name.
	Name string `json:"name"`

	// The VIP node description.
	Description string `json:"description"`

	// The node's IPv4 address (either IPv4 or IPv6 address must be specified).
	IPv4Address string `json:"ipv4Address,omitempty"`

	// The node's IPv6 address (either IPv4 or IPv6 address must be specified).
	IPv6Address string `json:"ipv6Address,omitempty"`

	// The node status (VIPNodeStatusEnabled, VIPNodeStatusDisabled, or VIPNodeStatusForcedOffline).
	Status string `json:"status"`

	// The Id of the node's associated health monitor (if any).
	HealthMonitorID string `json:"healthMonitorId,omitempty"`

	// The node's connection limit (must be greater than 0).
	ConnectionLimit int `json:"connectionLimit"`

	// The node's connection rate limit (must be greater than 0).
	ConnectionRateLimit int `json:"connectionRateLimit"`

	// The Id of the network domain where the node is located.
	NetworkDomainID string `json:"networkDomainId"`
}

// EditVIPNodeConfiguration represents the request body when editing a VIP node.
type EditVIPNodeConfiguration struct {
	// The VIP node Id.
	ID string `json:"id"`

	// The VIP node description.
	Description *string `json:"description,omitempty"`

	// The node status (VIPNodeStatusEnabled, VIPNodeStatusDisabled, or VIPNodeStatusForcedOffline).
	Status *string `json:"status,omitempty"`

	// The Id of the node's associated health monitor (if any).
	HealthMonitorID *string `json:"healthMonitorId,omitempty"`

	// The node's connection limit (must be greater than 0).
	ConnectionLimit *int `json:"connectionLimit,omitempty"`

	// The node's connection rate limit (must be greater than 0).
	ConnectionRateLimit *int `json:"connectionRateLimit,omitempty"`
}

// Request body for deleting a VIP node.
type deleteVIPNode struct {
	// The VIP node ID.
	ID string `json:"id"`
}

// ListVIPNodesInNetworkDomain retrieves a list of all VIP nodes in the specified network domain.
func (client *Client) ListVIPNodesInNetworkDomain(networkDomainID string, paging *Paging) (nodes *VIPNodes, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return nil, err
	}

	requestURI := fmt.Sprintf("%s/networkDomainVip/node?networkDomainId=%s&%s",
		url.QueryEscape(organizationID),
		url.QueryEscape(networkDomainID),
		paging.EnsurePaging().toQueryParameters(),
	)
	request, err := client.newRequestV22(requestURI, http.MethodGet, nil)
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

		return nil, apiResponse.ToError("Request to list VIP nodes in network domain '%s' failed with status code %d (%s): %s", networkDomainID, statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	nodes = &VIPNodes{}
	err = json.Unmarshal(responseBody, nodes)
	if err != nil {
		return nil, err
	}

	return nodes, nil
}

// GetVIPNode retrieves the VIP node with the specified Id.
// Returns nil if no VIP node is found with the specified Id.
func (client *Client) GetVIPNode(id string) (node *VIPNode, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return nil, err
	}

	requestURI := fmt.Sprintf("%s/networkDomainVip/node/%s",
		url.QueryEscape(organizationID),
		url.QueryEscape(id),
	)
	request, err := client.newRequestV22(requestURI, http.MethodGet, nil)
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

		return nil, apiResponse.ToError("Request to retrieve VIP node with Id '%s' failed with status code %d (%s): %s", id, statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	node = &VIPNode{}
	err = json.Unmarshal(responseBody, node)
	if err != nil {
		return nil, err
	}

	return node, nil
}

// CreateVIPNode creates a new VIP node.
// Returns the Id of the new node.
func (client *Client) CreateVIPNode(nodeConfiguration NewVIPNodeConfiguration) (nodeID string, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return "", err
	}

	requestURI := fmt.Sprintf("%s/networkDomainVip/createNode",
		url.QueryEscape(organizationID),
	)
	request, err := client.newRequestV22(requestURI, http.MethodPost, &nodeConfiguration)
	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return "", err
	}

	apiResponse, err := readAPIResponseAsJSON(responseBody, statusCode)
	if err != nil {
		return "", err
	}

	if apiResponse.ResponseCode != ResponseCodeOK {
		return "", apiResponse.ToError("Request to create VIP node '%s' failed with status code %d (%s): %s", nodeConfiguration.Name, statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	// Expected: "info" { "name": "nodeId", "value": "the-Id-of-the-new-node" }
	nodeIDMessage := apiResponse.GetFieldMessage("nodeId")
	if nodeIDMessage == nil {
		return "", apiResponse.ToError("Received an unexpected response (missing 'nodeId') with status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	return *nodeIDMessage, nil
}

// EditVIPNode updates an existing VIP node.
func (client *Client) EditVIPNode(id string, nodeConfiguration EditVIPNodeConfiguration) error {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return err
	}

	editNodeConfiguration := &nodeConfiguration
	editNodeConfiguration.ID = id

	requestURI := fmt.Sprintf("%s/networkDomainVip/editNode",
		url.QueryEscape(organizationID),
	)
	request, err := client.newRequestV22(requestURI, http.MethodPost, editNodeConfiguration)
	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return err
	}

	apiResponse, err := readAPIResponseAsJSON(responseBody, statusCode)
	if err != nil {
		return err
	}

	if statusCode != http.StatusOK {
		return apiResponse.ToError("Request to edit VIP node '%s' failed with status code %d (%s): %s", nodeConfiguration.ID, statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	return nil
}

// DeleteVIPNode deletes an existing VIP node.
// Returns an error if the operation was not successful.
func (client *Client) DeleteVIPNode(id string) (err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return err
	}

	requestURI := fmt.Sprintf("%s/networkDomainVip/deleteNode",
		url.QueryEscape(organizationID),
	)
	request, err := client.newRequestV22(requestURI, http.MethodPost, &deleteVIPNode{id})
	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return err
	}

	apiResponse, err := readAPIResponseAsJSON(responseBody, statusCode)
	if err != nil {
		return err
	}

	if apiResponse.ResponseCode != ResponseCodeOK {
		return apiResponse.ToError("Request to delete VIP node '%s' failed with unexpected status code %d (%s): %s", id, statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	return nil
}
