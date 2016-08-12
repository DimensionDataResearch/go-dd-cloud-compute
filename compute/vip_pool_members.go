package compute

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// VIPPoolMember represents a combination of node and port as a member of a VIP pool.
type VIPPoolMember struct {
	ID              string           `json:"id"`
	Pool            EntityReference  `json:"pool"`
	Node            VIPNodeReference `json:"node"`
	Port            *int             `json:"port,omitempty"`
	Status          string           `json:"status"`
	State           string           `json:"state"`
	NetworkDomainID string           `json:"networkDomainId"`
	DatacenterID    string           `json:"datacenterId"`
	CreateTime      string           `json:"createTime"`
}

// VIPPoolMembers represents a page of VIPPoolMember results.
type VIPPoolMembers struct {
	Items []VIPPoolMember `json:"poolMember"`

	PagedResult
}

// Request body for adding a VIP pool member.
type addPoolMember struct {
	PoolID string `json:"poolId"`
	NodeID string `json:"nodeId"`
	Status string `json:"status"`
	Port   *int   `json:"port,omitempty"`
}

// Request body for updating a VIP pool member.
type editPoolMember struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

// Request body for removing a VIP pool member.
type removePoolMember struct {
	ID string `json:"id"`
}

// ListVIPPoolMembers retrieves a list of all members of the specified VIP pool.
func (client *Client) ListVIPPoolMembers(poolID string, paging *Paging) (members *VIPPoolMembers, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return nil, err
	}

	requestURI := fmt.Sprintf("%s/networkDomainVip/poolMember?poolId=%s&%s",
		organizationID,
		url.QueryEscape(poolID),
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

		return nil, apiResponse.ToError("Request to list members of VIP pool '%s' failed with status code %d (%s): %s", poolID, statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	members = &VIPPoolMembers{}
	err = json.Unmarshal(responseBody, members)
	if err != nil {
		return nil, err
	}

	return members, nil
}

// ListVIPPoolMembershipsInNetworkDomain retrieves a list of all VIP pool memberships of the specified network domain.
func (client *Client) ListVIPPoolMembershipsInNetworkDomain(networkDomainID string, paging *Paging) (members *VIPPoolMembers, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return nil, err
	}

	requestURI := fmt.Sprintf("%s/networkDomainVip/poolMember?networkDomainId=%s&%s",
		organizationID,
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

		return nil, apiResponse.ToError("Request to list all VIP pool memberships in network domain '%s' failed with status code %d (%s): %s", networkDomainID, statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	members = &VIPPoolMembers{}
	err = json.Unmarshal(responseBody, members)
	if err != nil {
		return nil, err
	}

	return members, nil
}

// GetVIPPoolMember retrieves the VIP pool member with the specified Id.
// Returns nil if no VIP pool member is found with the specified Id.
func (client *Client) GetVIPPoolMember(id string) (member *VIPPoolMember, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return nil, err
	}

	requestURI := fmt.Sprintf("%s/networkDomainVip/poolMember/%s", organizationID, id)
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

		return nil, apiResponse.ToError("Request to retrieve VIP pool with Id '%s' failed with status code %d (%s): %s", id, statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	member = &VIPPoolMember{}
	err = json.Unmarshal(responseBody, member)
	if err != nil {
		return nil, err
	}

	return member, nil
}

// AddVIPPoolMember adds a VIP node as a member of a VIP pool.
// State must be one of VIPNodeStatusEnabled, VIPNodeStatusDisabled, or VIPNodeStatusDisabled
// Returns the member ID (uniquely identifies this combination of node, pool, and port).
func (client *Client) AddVIPPoolMember(poolID string, nodeID string, status string, port *int) (poolMemberID string, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return "", err
	}

	requestURI := fmt.Sprintf("%s/networkDomainVip/addPoolMember", organizationID)
	request, err := client.newRequestV22(requestURI, http.MethodPost, &addPoolMember{
		PoolID: poolID,
		NodeID: nodeID,
		Status: status,
		Port:   port,
	})
	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return "", err
	}

	apiResponse, err := readAPIResponseAsJSON(responseBody, statusCode)
	if err != nil {
		return "", err
	}

	if statusCode != http.StatusOK {
		return "", apiResponse.ToError("Request to add VIP node '%s' as a member of pool '%s' failed with status code %d (%s): %s", nodeID, poolID, statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	// Expected: "info" { "name": "poolMemberId", "value": "the-Id-of-the-new-pool-member" }
	poolMemberIDMessage := apiResponse.GetFieldMessage("poolMemberId")
	if poolMemberIDMessage == nil {
		return "", apiResponse.ToError("Received an unexpected response (missing 'poolMemberId') with status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	return *poolMemberIDMessage, nil
}

// EditVIPPoolMember updates the status of an existing VIP pool member.
// status can be VIPNodeStatusEnabled, VIPNodeStatusDisabled, or VIPNodeStatusForcedOffline
func (client *Client) EditVIPPoolMember(id string, status string) error {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return err
	}

	requestURI := fmt.Sprintf("%s/networkDomainVip/editPoolMember", organizationID)
	request, err := client.newRequestV22(requestURI, http.MethodPost, &editPoolMember{
		ID:     id,
		Status: status,
	})
	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return err
	}

	apiResponse, err := readAPIResponseAsJSON(responseBody, statusCode)
	if err != nil {
		return err
	}

	if statusCode != http.StatusOK {
		return apiResponse.ToError("Request to edit VIP pool member '%s' failed with status code %d (%s): %s", id, statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	return nil
}

// RemoveVIPPoolMember removes a VIP pool member.
func (client *Client) RemoveVIPPoolMember(id string) error {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return err
	}

	requestURI := fmt.Sprintf("%s/networkDomainVip/removePoolMember", organizationID)
	request, err := client.newRequestV22(requestURI, http.MethodPost, &removePoolMember{id})
	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return err
	}

	apiResponse, err := readAPIResponseAsJSON(responseBody, statusCode)
	if err != nil {
		return err
	}

	if statusCode != http.StatusOK {
		return apiResponse.ToError("Request to remove member '%s' from its pool failed with status code %d (%s): %s", id, statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	return nil
}
