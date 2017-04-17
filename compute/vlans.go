package compute

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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
	NetworkDomain EntityReference `json:"networkDomain"`

	// The VLAN's associated IPv4 network range.
	IPv4Range IPv4Range `json:"privateIpv4Range"`

	// The VLAN's IPv4 gateway address.
	IPv4GatewayAddress string `json:"ipv4GatewayAddress"`

	// The VLAN's associated IPv6 network range.
	IPv6Range IPv6Range `json:"ipv6Range"`

	// The VLAN's IPv6 gateway address.
	IPv6GatewayAddress string `json:"ipv6GatewayAddress"`

	// The date / time that the VLAN was first created.
	CreateTime string `json:"createTime"`

	// The VLAN's current state.
	State string `json:"state"`

	// The ID of the data center in which the VLAN and its containing network domain are deployed.
	DataCenterID string `json:"datacenterId"`
}

// GetID returns the VLAN's Id.
func (vlan *VLAN) GetID() string {
	return vlan.ID
}

// GetResourceType returns the network domain's resource type.
func (vlan *VLAN) GetResourceType() ResourceType {
	return ResourceTypeVLAN
}

// GetName returns the VLAN's name.
func (vlan *VLAN) GetName() string {
	return vlan.Name
}

// GetState returns the VLAN's current state.
func (vlan *VLAN) GetState() string {
	return vlan.State
}

// IsDeleted determines whether the VLAN has been deleted (is nil).
func (vlan *VLAN) IsDeleted() bool {
	return vlan == nil
}

var _ Resource = &VLAN{}

// ToEntityReference creates an EntityReference representing the VLAN.
func (vlan *VLAN) ToEntityReference() EntityReference {
	return EntityReference{
		ID:   vlan.ID,
		Name: vlan.Name,
	}
}

var _ NamedEntity = &VLAN{}

// VLANs represents the response to a "List VLANs" API call.
type VLANs struct {
	// The current page of network domains.
	VLANs []VLAN `json:"vlan"`

	PagedResult
}

// DeployVLAN represents the request body when deploying a cloud compute VLAN.
type DeployVLAN struct {
	// The Id of the network domain in which the VLAN will be deployed.
	VLANID string `json:"networkDomainId"`

	// The VLAN name.
	Name string `json:"name"`

	// The VLAN description.
	Description string `json:"description"`

	// The private IPv4 base address for the VLAN.
	IPv4BaseAddress string `json:"privateIpv4BaseAddress"`

	// The private IPv4 prefix size (i.e. netmask) for the VLAN.
	IPv4PrefixSize int `json:"privateIpv4PrefixSize"`
}

// EditVLAN represents the request body when editing a cloud compute VLAN.
type EditVLAN struct {
	// The ID of the VLAN to edit.
	ID string `json:"id"`

	// The VLAN name (optional).
	Name *string `json:"name,omitempty"`

	// The VLAN description (optional).
	Description *string `json:"description,omitempty"`
}

// DeleteVLAN represents a request to delete a compute VLAN.
type DeleteVLAN struct {
	// The VLAN Id.
	ID string `json:"id"`
}

// GetVLAN retrieves the VLAN with the specified Id.
// id is the Id of the VLAN to retrieve.
// Returns nil if no VLAN is found with the specified Id.
func (client *Client) GetVLAN(id string) (vlan *VLAN, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return nil, err
	}

	requestURI := fmt.Sprintf("%s/network/vlan/%s",
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

		return nil, apiResponse.ToError("Request to retrieve VLAN failed with status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	vlan = &VLAN{}
	err = json.Unmarshal(responseBody, vlan)

	return vlan, err
}

// GetVLANByName retrieves the VLAN (if any) with the specified name in the specified network domain.
func (client *Client) GetVLANByName(name string, networkDomainID string) (*VLAN, error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return nil, err
	}

	requestURI := fmt.Sprintf("%s/network/vlan?name=%s&networkDomainId=%s",
		url.QueryEscape(organizationID),
		url.QueryEscape(name),
		url.QueryEscape(networkDomainID),
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

		return nil, apiResponse.ToError("Request failed with status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	vlans := &VLANs{}
	err = json.Unmarshal(responseBody, vlans)
	if err != nil {
		return nil, err
	}
	if vlans.IsEmpty() {
		return nil, nil // No matching VLAN was found.
	}

	if len(vlans.VLANs) != 1 {
		return nil, fmt.Errorf("found multiple VLANs (%d) named '%s' in network domain '%s'", len(vlans.VLANs), name, networkDomainID)
	}

	return &vlans.VLANs[0], nil
}

// ListVLANs retrieves a list of all VLANs in the specified network domain.
// TODO: Support filtering and sorting.
func (client *Client) ListVLANs(networkDomainID string, paging *Paging) (vlans *VLANs, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return nil, err
	}

	requestURI := fmt.Sprintf("%s/network/vlan?networkDomainId=%s&%s",
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

		return nil, apiResponse.ToError("Request to list VLANs failed with status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	vlans = &VLANs{}
	err = json.Unmarshal(responseBody, vlans)

	return vlans, err
}

// DeployVLAN deploys a new VLAN into a network domain.
func (client *Client) DeployVLAN(networkDomainID string, name string, description string, ipv4BaseAddress string, ipv4PrefixSize int) (vlanID string, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return "", err
	}

	requestURI := fmt.Sprintf("%s/network/deployVlan",
		url.QueryEscape(organizationID),
	)
	request, err := client.newRequestV22(requestURI, http.MethodPost, &DeployVLAN{
		VLANID:          networkDomainID,
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
		return "", apiResponse.ToError("Request to deploy VLAN '%s' failed with status code %d (%s): %s", name, statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	// Expected: "info" { "name": "vlanId", "value": "the-Id-of-the-new-VLAN" }
	vlanIDMessage := apiResponse.GetFieldMessage("vlanId")
	if vlanIDMessage == nil {
		return "", apiResponse.ToError("Received an unexpected response (missing 'vlanId') with status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	return *vlanIDMessage, nil
}

// EditVLAN updates an existing VLAN.
// Pass an empty string for any field to retain its existing value.
// Returns an error if the operation was not successful.
func (client *Client) EditVLAN(id string, name *string, description *string) (err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return err
	}

	requestURI := fmt.Sprintf("%s/network/editVlan",
		url.QueryEscape(organizationID),
	)
	request, err := client.newRequestV22(requestURI, http.MethodPost, &EditVLAN{
		ID:          id,
		Name:        name,
		Description: description,
	})
	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return err
	}

	apiResponse, err := readAPIResponseAsJSON(responseBody, statusCode)
	if err != nil {
		return err
	}

	if apiResponse.ResponseCode != ResponseCodeOK {
		return apiResponse.ToError("Request to edit VLAN failed with unexpected status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	return nil
}

// DeleteVLAN deletes an existing VLAN.
// Returns an error if the operation was not successful.
func (client *Client) DeleteVLAN(id string) (err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return err
	}

	requestURI := fmt.Sprintf("%s/network/deleteVlan",
		url.QueryEscape(organizationID),
	)
	request, err := client.newRequestV22(requestURI, http.MethodPost, &DeleteVLAN{id})
	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return err
	}

	apiResponse, err := readAPIResponseAsJSON(responseBody, statusCode)
	if err != nil {
		return err
	}

	if apiResponse.ResponseCode != ResponseCodeInProgress {
		return apiResponse.ToError("Request to delete VLAN failed with unexpected status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	return nil
}
