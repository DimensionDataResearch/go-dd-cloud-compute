package compute

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// ReservedIPAddress represents a private IP address that has been reserved for use on a VLAN.
type ReservedIPAddress struct {
	IPAddress    string `json:"value"`
	VLANID       string `json:"vlanId"`
	DatacenterID string `json:"datacenterId"`
}

// ReservedIPv4Addresses represents a page of ReservedIPAddress results for reserved IPv4 addresses.
type ReservedIPv4Addresses struct {
	PagedResult

	Items []ReservedIPAddress `json:"ipv4"`
}

// ReservedIPv6Addresses represents a page of ReservedIPAddress results for reserved IPv6 addresses.
type ReservedIPv6Addresses struct {
	PagedResult

	Items []ReservedIPAddress `json:"reservedIpv6Address"`
}

// ListReservedIPv4AddressesInVLAN retrieves all port lists associated with the specified VLAN.
func (client *Client) ListReservedIPv4AddressesInVLAN(vlanID string) (reservedIPAddresses *ReservedIPv4Addresses, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return nil, err
	}

	requestURI := fmt.Sprintf("%s/network/portList?vlanId=%s",
		url.QueryEscape(organizationID),
		url.QueryEscape(vlanID),
	)
	request, err := client.newRequestV24(requestURI, http.MethodGet, nil)
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

		return nil, apiResponse.ToError("Request to list port lists failed with status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	reservedIPAddresses = &ReservedIPv4Addresses{}
	err = json.Unmarshal(responseBody, reservedIPAddresses)

	return reservedIPAddresses, err
}

// ListReservedIPv6AddressesInVLAN retrieves all port lists associated with the specified VLAN.
func (client *Client) ListReservedIPv6AddressesInVLAN(vlanID string) (reservedIPAddresses *ReservedIPv6Addresses, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return nil, err
	}

	requestURI := fmt.Sprintf("%s/network/portList?vlanId=%s",
		url.QueryEscape(organizationID),
		url.QueryEscape(vlanID),
	)
	request, err := client.newRequestV24(requestURI, http.MethodGet, nil)
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

		return nil, apiResponse.ToError("Request to list port lists failed with status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	reservedIPAddresses = &ReservedIPv6Addresses{}
	err = json.Unmarshal(responseBody, reservedIPAddresses)

	return reservedIPAddresses, err
}
