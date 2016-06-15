package compute

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// IPAddressList represents an IP address list.
type IPAddressList struct {
	ID          string               `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	IPVersion   string               `json:"ipVersion"`
	State       string               `json:"state"`
	CreateTime  string               `json:"createTime"`
	Addresses   []IPAddressListEntry `json:"ipAddress"`
	ChildLists  []EntitySummary      `json:"childIpAddressList"`
}

// IPAddressLists represents a page of IPAddressList results.
type IPAddressLists struct {
	AddressLists []IPAddressList `json:"ipAddressList"`

	PagedResult
}

// IPAddressListEntry represents an entry in an IP address list.
type IPAddressListEntry struct {
	Begin      string  `json:"begin"`
	End        *string `json:"end,omitempty"`
	PrefixSize *int    `json:"prefixSize,omitempty"`
}

// GetIPAddressList retrieves the IP address list with the specified Id.
// id is the Id of the IP address list to retrieve.
// Returns nil if no addressList is found with the specified Id.
func (client *Client) GetIPAddressList(id string) (addressList *IPAddressList, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return nil, err
	}

	requestURI := fmt.Sprintf("%s/network/ipAddressList/%s", organizationID, id)
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

		return nil, apiResponse.ToError("Request to retrieve IP address list failed with status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	addressList = &IPAddressList{}
	err = json.Unmarshal(responseBody, addressList)

	return addressList, err
}

// ListIPAddressLists retrieves all IP address lists associated with the specified network domain.
func (client *Client) ListIPAddressLists(networkDomainID string) (addressLists *IPAddressLists, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return nil, err
	}

	requestURI := fmt.Sprintf("%s/network/ipAddressList?networkDomainId=%s", organizationID, networkDomainID)
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

		return nil, apiResponse.ToError("Request to list IP address lists failed with status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	addressLists = &IPAddressLists{}
	err = json.Unmarshal(responseBody, addressLists)

	return addressLists, err
}
