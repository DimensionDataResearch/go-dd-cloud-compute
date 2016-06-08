package compute

import (
	"fmt"
	"net/http"
)

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
