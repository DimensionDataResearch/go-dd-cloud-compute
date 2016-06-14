package compute

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// PublicIPBlock represents an allocated block of public IPv4 addresses.
type PublicIPBlock struct {
	ID              string `json:"id"`
	NetworkDomainID string `json:"networkDomainId"`
	DataCenterID    string `json:"datacenterId"`
	BaseIP          string `json:"baseIp"`
	Size            int    `json:"size"`
	CreateTime      string `json:"createTime"`
	State           string `json:"state"`
}

// GetID returns the public IPv4 address block's Id.
func (block *PublicIPBlock) GetID() string {
	return block.ID
}

// GetName returns the public IPv4 address block's name.
func (block *PublicIPBlock) GetName() string {
	return fmt.Sprintf("%s+%d", block.BaseIP, block.Size)
}

// GetState returns the network block's current state.
func (block *PublicIPBlock) GetState() string {
	return block.State
}

// IsDeleted determines whether the public IPv4 address block has been deleted (is nil).
func (block *PublicIPBlock) IsDeleted() bool {
	return block == nil
}

var _ Resource = &PublicIPBlock{}

// PublicIPBlocks represents a page of PublicIPBlock results.
type PublicIPBlocks struct {
	Blocks []PublicIPBlock `json:"publicIpBlock"`

	PagedResult
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

// GetPublicIPBlock retrieves the public IPv4 address block with the specified Id.
// Returns nil if no IPv4 address block is found with the specified Id.
func (client *Client) GetPublicIPBlock(id string) (block *PublicIPBlock, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return nil, err
	}

	requestURI := fmt.Sprintf("%s/network/publicIpBlock/%s", organizationID, id)
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

		return nil, fmt.Errorf("Request to retrieve public IPv4 address block failed with status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	block = &PublicIPBlock{}
	err = json.Unmarshal(responseBody, block)
	if err != nil {
		return nil, err
	}

	return block, nil
}

// ListPublicIPBlocks retrieves all blocks of public IPv4 addresses that have been allocated to the specified network domain.
func (client *Client) ListPublicIPBlocks(networkDomainID string) (blocks *PublicIPBlocks, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return nil, err
	}

	requestURI := fmt.Sprintf("%s/network/publicIpBlock?networkDomainId=%s", organizationID, networkDomainID)
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

		return nil, fmt.Errorf("Request to list public IPv4 address blocks failed with status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	blocks = &PublicIPBlocks{}
	err = json.Unmarshal(responseBody, blocks)

	return blocks, err
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

		return nil, fmt.Errorf("Request to list reserved public IPv4 addresses failed with status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	reservedPublicIPs = &ReservedPublicIPs{}
	err = json.Unmarshal(responseBody, reservedPublicIPs)

	return reservedPublicIPs, err
}
