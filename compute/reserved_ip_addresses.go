package compute

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

//// ReservedIPAddress represents a private IP address that has been reserved for use on a VLAN.
//type ReservedIPAddress struct {
//	IPAddress    string `json:"value"`
//	VLANID       string `json:"vlanId"`
//	DatacenterID string `json:"datacenterId"`
//}

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

// Request body when reserving an IP address.
type ReservedIPAddress struct {
	IPAddress string `json:"ipAddress"`
	VLANID    string `json:"vlanId"`
	Description  string `json:"description"`
}

// ListReservedPrivateIPv4AddressesInVLAN retrieves all private IPv4 addresses reserved in the specified VLAN.
func (client *Client) ListReservedPrivateIPv4AddressesInVLAN(vlanID string) (reservedIPAddresses *ReservedIPv4Addresses, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return nil, err
	}

	requestURI := fmt.Sprintf("%s/network/reservedPrivateIpv4Address?vlanId=%s",
		url.QueryEscape(organizationID),
		url.QueryEscape(vlanID),
	)
	request, err := client.newRequestV29(requestURI, http.MethodGet, nil)
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

		return nil, apiResponse.ToError("Request to list reserved IPv4 addresses failed with status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	reservedIPAddresses = &ReservedIPv4Addresses{}
	err = json.Unmarshal(responseBody, reservedIPAddresses)

	return reservedIPAddresses, err
}

// ReservePrivateIPv4Address creates a reservation for a private IPv4 address on a VLAN.
func (client *Client) ReservePrivateIPv4Address(vlanID string, ipAddress string, description string) error {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return err
	}

	requestURI := fmt.Sprintf("%s/network/reservePrivateIpv4Address",
		url.QueryEscape(organizationID),
	)
	request, err := client.newRequestV29(requestURI, http.MethodPost, &ReservedIPAddress{
		IPAddress: ipAddress,
		VLANID:    vlanID,
		Description: description,
	})

	if err != nil {
		return err
	}

	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return err
	}

	apiResponse, err := readAPIResponseAsJSON(responseBody, statusCode)
	if err != nil {
		return err
	}
	if apiResponse.ResponseCode != ResponseCodeOK {
		return apiResponse.ToError("Request to create reservation for private IPv4 address failed with status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	return nil
}

// UnreservePrivateIPv4Address removes the reservation (if any) for a private IPv4 address on a VLAN.
func (client *Client) UnreservePrivateIPv4Address(vlanID string, ipAddress string, description string) error {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return err
	}

	requestURI := fmt.Sprintf("%s/network/unreservePrivateIpv4Address",
		url.QueryEscape(organizationID),
	)
	request, err := client.newRequestV29(requestURI, http.MethodPost, &ReservedIPAddress{
		IPAddress: ipAddress,
		VLANID:    vlanID,
		Description: description,
	})
	if err != nil {
		return err
	}

	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return err
	}

	apiResponse, err := readAPIResponseAsJSON(responseBody, statusCode)
	if err != nil {
		return err
	}
	if apiResponse.ResponseCode != ResponseCodeOK {
		return apiResponse.ToError("Request to remove reservation for private IPv4 address failed with status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	return nil
}

// ListReservedIPv6AddressesInVLAN retrieves all IPv6 addresses reserved in the specified VLAN.
func (client *Client) ListReservedIPv6AddressesInVLAN(vlanID string) (reservedIPAddresses *ReservedIPv6Addresses, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return nil, err
	}

	requestURI := fmt.Sprintf("%s/network/reservedIpv6Address?vlanId=%s",
		url.QueryEscape(organizationID),
		url.QueryEscape(vlanID),
	)
	request, err := client.newRequestV29(requestURI, http.MethodGet, nil)
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

		return nil, apiResponse.ToError("Request to list reserved IPv6 addresses failed with status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	reservedIPAddresses = &ReservedIPv6Addresses{}
	err = json.Unmarshal(responseBody, reservedIPAddresses)

	return reservedIPAddresses, err
}

// ReserveIPv6Address creates a reservation for an IPv6 address on a VLAN.
func (client *Client) ReserveIPv6Address(vlanID string, ipAddress string, description string) error {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return err
	}

	requestURI := fmt.Sprintf("%s/network/reserveIpv6Address",
		url.QueryEscape(organizationID),
	)
	request, err := client.newRequestV29(requestURI, http.MethodPost, &ReservedIPAddress{
		IPAddress: ipAddress,
		VLANID:    vlanID,
		Description: description,
	})
	if err != nil {
		return err
	}

	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return err
	}

	apiResponse, err := readAPIResponseAsJSON(responseBody, statusCode)
	if err != nil {
		return err
	}
	if apiResponse.ResponseCode != ResponseCodeOK {
		return apiResponse.ToError("Request to reserve IPV6 address failed with status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	return nil
}

// UnreserveIPv6Address removes the reservation (if any) for an IPv6 address on a VLAN.
func (client *Client) UnreserveIPv6Address(vlanID string, ipAddress string, description string) error {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return err
	}

	requestURI := fmt.Sprintf("%s/network/unreserveIpv6Address",
		url.QueryEscape(organizationID),
	)
	request, err := client.newRequestV29(requestURI, http.MethodPost, &ReservedIPAddress{
		IPAddress: ipAddress,
		VLANID:    vlanID,
		Description: description,
	})
	if err != nil {
		return err
	}

	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return err
	}

	apiResponse, err := readAPIResponseAsJSON(responseBody, statusCode)
	if err != nil {
		return err
	}
	if apiResponse.ResponseCode != ResponseCodeOK {
		return apiResponse.ToError("Request to remove IP address reservation failed with status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	return nil
}
