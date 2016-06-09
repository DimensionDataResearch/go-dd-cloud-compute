package compute

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// VLAN represents a compute VLAN.
type VLAN struct {
	// The VLAN Id.
	ID string `json:"id"`

	// The VLAN name.
	Name string `json:"name"`

	// The VLAN description.
	Description string `json:"description"`

	// The network domain in which the VLAN is deployed associated.
	NetworkDomain EntitySummary `json:"networkDomain"`

	// The VLAN's associated IPv4 network range.
	IPv4Range IPv4Range `json:"privateIpv4Range"`

	// The VLAN's IPv4 gateway address.
	IPv4GatewayAddress string `json:"ipv4GatewayAddress"`

	// The VLAN's associated IPv6 network range.
	IPv6Range IPv6Range `json:"ipv6Range"`

	// The VLAN's IPv6 gateway address.
	IPv6GatewayAddress string `json:"ipv6GatewayAddress"`

	// The date / time that the VLAN was first created.
	CreateTime int `json:"createTime"`

	// The VLAN's current state.
	State string `json:"state"`

	// The ID of the data center in which the VLAN and its containing network domain are deployed.
	DataCenterID string `json:"datacenterId"`
}

// DeployVLAN represents the request body when deploying a cloud compute VLAN.
type DeployVLAN struct {
	// The Id of the network domain in which the VLAN will be deployed.
	NetworkDomainID string `json:"networkDomainId"`

	// The VLAN name.
	Name string `json:"name"`

	// The VLAN description.
	Description string `json:"description"`

	// The private IPv4 base address for the VLAN.
	IPv4BaseAddress string `json:"privateIpv4BaseAddress"`

	// The private IPv4 prefix size (i.e. netmask) for the VLAN.
	IPv4PrefixSize int `json:"privateIpv4PrefixSize"`
}

// GetVLAN retrieves the VLAN with the specified Id.
// id is the Id of the VLAN to retrieve.
// Returns nil if no VLAN is found with the specified Id.
func (client *Client) GetVLAN(id string) (vlan *VLAN, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return nil, err
	}

	requestURI := fmt.Sprintf("%s/network/vlan/%s", organizationID, id)
	request, err := client.newRequestV22(requestURI, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		var apiResponse *APIResponse

		apiResponse, err = readAPIResponseAsJSON(responseBody, statusCode)
		if err != nil {
			return nil, err
		}

		if apiResponse.ResponseCode == ResponseCodeResourceNotFound {
			return nil, nil // Not an error, but was not found.
		}

		return nil, fmt.Errorf("Request to retrieve VLAN failed with status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	vlan = &VLAN{}
	err = json.Unmarshal(responseBody, vlan)

	return vlan, err
}

// DeployVLAN deploys a new VLAN into a network domain.
func (client *Client) DeployVLAN(networkDomainID string, name string, description string, ipv4BaseAddress string, ipv4PrefixSize int) (vlanID string, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return "", err
	}

	requestURI := fmt.Sprintf("%s/network/deployVlan", organizationID)
	request, err := client.newRequestV22(requestURI, http.MethodPost, &DeployVLAN{
		NetworkDomainID: networkDomainID,
		Name:            name,
		Description:     description,
		IPv4BaseAddress: ipv4BaseAddress,
		IPv4PrefixSize:  ipv4PrefixSize,
	})
	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return "", err
	}

	apiResponse, err := readAPIResponseAsJSON(responseBody, statusCode)
	if err != nil {
		return "", err
	}

	if apiResponse.ResponseCode != ResponseCodeInProgress {
		return "", fmt.Errorf("Request to deploy VLAN '%s' failed with status code %d (%s): %s", name, statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	// Expected: "info" { "name": "vlanId", "value": "the-Id-of-the-new-VLAN" }
	if len(apiResponse.FieldMessages) != 1 || apiResponse.FieldMessages[0].FieldName != "vlanId" {
		return "", fmt.Errorf("Received an unexpected response (missing 'vlanId') with status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	return apiResponse.FieldMessages[0].Message, nil
}
