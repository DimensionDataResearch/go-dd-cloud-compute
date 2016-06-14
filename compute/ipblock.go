package compute

import (
	"encoding/json"
	"fmt"
	"net/http"
)

/*
	{
		"networkDomainId": "802abc9f-45a7-4efb-9d5a-810082368708",
		"baseIp": "165.180.12.12",
		"size": 2,
		"createTime": "2014-12-15T16:35:07.000Z",
		"state": "NORMAL",
		"id": "cacc028a-7f12-11e4-a91c-0030487e0302", "datacenterId": "NA9"
	}
*/

// PublicIPBlock represents an allocated block of public IPv4 addresses.
type PublicIPBlock struct {
	ID              string `json:"id"`
	NetworkDomainID string `json:"networkDomainId"`
	BaseIP          string `json:"baseIp"`
	Size            int    `json:"size"`
	CreateTime      string `json:"createTime"`
	State           string `json:"state"`
}

// ReservedPublicIP represents a public IPv4 address reserved for NAT or a VIP.
type ReservedPublicIP struct {
	IPBlockID       string `json:"ipBlockId"`
	DataCenterID    string `json:"datacenterId"`
	NetworkDomainID string `json:"networkDomainId"`
	Address         string `json:"value"`
}

// ReservedPublicIPs represents a page of ReservedPublicIP results.
type ReservedPublicIPs struct {
	IPs []ReservedPublicIP `json:"ip"`

	PagedResult
}

// ListReservedPublicIPAddresses retrieves all public IPv4 addresses in the specified network domain that have been reserved in the specified network domain.
func (client *Client) ListReservedPublicIPAddresses(networkDomainID string) (reservedPublicIPs *ReservedPublicIPs, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return nil, err
	}

	requestURI := fmt.Sprintf("%s/network/reservedPublicIpv4Address?networkDomainId=%s", organizationID, networkDomainID)
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

		return nil, fmt.Errorf("Request to list ReservedPublicIPs failed with status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	reservedPublicIPs = &ReservedPublicIPs{}
	err = json.Unmarshal(responseBody, reservedPublicIPs)

	return reservedPublicIPs, err
}
