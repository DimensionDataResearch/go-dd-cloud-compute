package compute

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

const (
	// FirewallRuleIPVersion4 indicates a firewall rule that targets IPv4
	IPVersion4 = "IPv4"

	// FirewallRuleIPVersion6 indicates a firewale rule that targets IPv6
	IPVersion6 = "IPv6"
)

// StaticRoute reporesents client static route on a network domain in an MCP2 data center.
type StaticRoute struct {
	// UUID of a Network Domain belonging to {org-id} within which the Static Route is to be created.
	NetworkDomainId string `json:"networkDomainId"`

	// Must be between 1 and 75 characters in length.
	//Cannot start with a number, a period ('.'), 'CCSYSTEM.' or 'CCDEFAULT.'.
	Name string `json:"name"`

	// Maximum length: 255 characters.
	Description string `json:"description"`

	// Type
	Type string `json:"type"`

	// One of IPV4 or IPV6
	IpVersion string `json:"ipVersion"`

	// Either a valid IPv4 address in dot-decimal notation or an IPv6 address in compressed or extended format.
	// In conjunction with the destinationPrefixSize this must represent a CIDR boundary.
	DestinationNetworkAddress string `json:"destinationNetworkAddress"`

	// Integer prefix defining the size of the network.
	// In conjunction with the destinationPrefixSize this must represent a CIDR boundary.
	DestinationPrefixSize int `json:"destinationPrefixSize"`

	// Gateway address in the form of an INET gateway, CPNC gateway or an address on an Attached VLAN in the same Network Domain.
	NextHopAddress string `json:"nextHopAddress"`

	// State
	State string `json:"state"`

	// The date/time that the Static Route was created in CloudControl
	CreateTime string `json:"createTime"`

	// Static Route ID
	ID string `json:"id"`

	// Data Center
	DataCenter string `json:"datacenterId"`

	// If an IP address on an Attached VLAN in the same Network Domain was provided as the nextHopAddress when creating
	// the Static Route, then the VLAN's UUID is included. This can be used with the Get VLAN API to retrieve full VLAN
	// details.
	// NextHopAddressVlanId string `json:"nextHopAddressVlanId"`

}

// Get ID
func (staticRoute *StaticRoute) GetID() string {
	return staticRoute.ID
}

// Get name
func (staticRoute *StaticRoute) GetName() string {
	return staticRoute.Name
}



// GetResourceType returns the static route resource type.
func (staticRoute *StaticRoute) GetResourceType() ResourceType {
	return ResourceTypeStaticRoutes
}


// GetState returns the Static Routes current state.
func (staticRoute *StaticRoute) GetState() string {
	return staticRoute.State
}

// IsDeleted determines whether the network domain has been deleted (is nil).
func (staticRoute *StaticRoute) IsDeleted() bool {
	// TODO
	return staticRoute == nil
}

// ToEntityReference creates an EntityReference representing the VLAN.
func (staticRoute *StaticRoute) ToEntityReference() EntityReference {
	return EntityReference{
		ID:   staticRoute.ID,
		Name: staticRoute.Name,
	}
}

var _ Resource = &StaticRoute{}

// Request body for creating a Static Route
type CreateStaticRoute struct {
	// UUID of a Network Domain belonging to {org-id} within which the Static Route is to be created.
	NetworkDomainId string `json:"networkDomainId"`

	// Must be between 1 and 75 characters in length.
	//Cannot start with a number, a period ('.'), 'CCSYSTEM.' or 'CCDEFAULT.'.
	Name string `json:"name"`

	// Maximum length: 255 characters.
	Description string `json:"description"`

	// One of IPV4 or IPV6
	IpVersion string `json:"ipVersion"`

	// Either a valid IPv4 address in dot-decimal notation or an IPv6 address in compressed or extended format.
	// In conjunction with the destinationPrefixSize this must represent a CIDR boundary.
	DestinationNetworkAddress string `json:"destinationNetworkAddress"`

	// Integer prefix defining the size of the network.
	// In conjunction with the destinationPrefixSize this must represent a CIDR boundary.
	DestinationPrefixSize int `json:"destinationPrefixSize"`

	// Gateway address in the form of an INET gateway, CPNC gateway or an address on an Attached VLAN in the same Network Domain.
	NextHopAddress string `json:"nextHopAddress"`
}

type StaticRoutes struct {
	// The current page of network domains.
	Routes []StaticRoute `json:"staticRoute"`

	PagedResult
}

type deleteStaticRoute struct {
	ID string `json:"id"`
}

type restoreStaticRoute struct {
	NetworkDomainId string `json:"networkDomainId"`
}


// Create enterprise static route
func (client *Client) CreateStaticRoute(networkDomainId string, name string, description string, ipVersion string,
	destinationNetworkAddress string, destinationPrefixSize int, nextHopAddress string) (staticRouteID string, err error){

	organizationID, err := client.getOrganizationID()
	if err != nil {
		return "", err
	}

	requestURI := fmt.Sprintf("%s/network/createStaticRoute", url.QueryEscape(organizationID))
	request, err := client.newRequestV29(requestURI, http.MethodPost, &CreateStaticRoute{
		NetworkDomainId: 			networkDomainId,
		Name:						name,
		Description: 				description,
		IpVersion: 					ipVersion,
		DestinationNetworkAddress: 	destinationNetworkAddress,
		DestinationPrefixSize:		destinationPrefixSize,
		NextHopAddress: 			nextHopAddress,

	})

	responseBody, statusCode, err := client.executeRequest(request)

	if err != nil {
		return "", err
	}

	apiResponse, err := readAPIResponseAsJSON(responseBody, statusCode)


	if statusCode != http.StatusOK {
		if apiResponse.ResponseCode == ResponseCodeResourceNotFound {
			return "", nil // Not an error, but was not found.
		}

		log.Printf("Request to create Static Route failed with status_code:%d response_code:%s  Msg: %s",
			statusCode, apiResponse.ResponseCode, apiResponse.Message)

		return "", apiResponse.ToError(
			"Request to create Static Route failed with status code %d (%s): %s",
			statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	staticRouteIdCreated := apiResponse.GetFieldMessage("staticRouteId")

	if staticRouteIdCreated == nil {
		return "", apiResponse.ToError("Unknown error occured. Request to create Static Route failed" +
			" with status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	return *staticRouteIdCreated, nil
}

// List static route of a network domain
func (client *Client) ListStaticRoute(paging *Paging) (staticRoutes *StaticRoutes, err error){
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return nil, err
	}

	requestURI := fmt.Sprintf("%s/network/staticRoute?%s",
		url.QueryEscape(organizationID),
		paging.EnsurePaging().toQueryParameters(),
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

		return nil, apiResponse.ToError("Request failed with status code %d (%s): %s",
			statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	staticRoutes = &StaticRoutes{}
	err = json.Unmarshal(responseBody, staticRoutes)
	if err != nil {
		return nil, err
	}

	return staticRoutes, nil
}

// List static route of a network domain by Name
func (client *Client) GetStaticRouteByName(name string) (staticRoute *StaticRoute, err error){
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return nil, err
	}

	requestURI := fmt.Sprintf("%s/network/staticRoute?name=%s",
		url.QueryEscape(organizationID),
		url.QueryEscape(name),
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

		return nil, apiResponse.ToError("Request failed with status code %d (%s): %s",
			statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	staticRoutes := &StaticRoutes{}
	err = json.Unmarshal(responseBody, staticRoutes)
	if err != nil {
		return nil, err
	}

	return &staticRoutes.Routes[0], nil
}

func (client *Client) GetStaticRoute(id string) (staticRoute *StaticRoute, err error){
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return nil, err
	}

	requestURI := fmt.Sprintf("%s/network/staticRoute/%s",
		url.QueryEscape(organizationID),
		url.QueryEscape(id),
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

		return nil, apiResponse.ToError("Request failed with status code %d (%s): %s",
			statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	staticRoute = &StaticRoute{}
	err = json.Unmarshal(responseBody, staticRoute)
	if err != nil {
		return nil, err
	}

	return staticRoute, nil
}

// Delete static route
func (client *Client) DeleteStaticRoute(id string) (err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return err
	}

	requestURI := fmt.Sprintf("%s/network/deleteStaticRoute",
		url.QueryEscape(organizationID),
	)

	request, err := client.newRequestV29(requestURI, http.MethodPost, &deleteStaticRoute{id})
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
		return apiResponse.ToError("Request to delete Static Route failed with unexpected status code %d (%s): %s",
			statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	return nil
}

// Restores the Static Routes of a Network Domain (networkDomainId) belonging to the organization identified by {org-id}
// to the system default (also referred to as baseline) configuration applied when the Network Domain was first deployed
func (client *Client) RestoreStaticRoute(networkDomainId string) (err error){
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return err
	}

	requestURI := fmt.Sprintf("%s/network/restoreStaticRoutes",
		url.QueryEscape(organizationID),
	)
	request, err := client.newRequestV24(requestURI, http.MethodPost, &restoreStaticRoute{networkDomainId})
	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return err
	}

	apiResponse, err := readAPIResponseAsJSON(responseBody, statusCode)
	if err != nil {
		return err
	}

	if apiResponse.ResponseCode != ResponseCodeInProgress {
		return apiResponse.ToError("Request to restore default static route failed with "+
			"unexpected status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	return nil
}

