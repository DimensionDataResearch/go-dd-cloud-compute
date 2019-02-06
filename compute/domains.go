package compute

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// NetworkDomain represents a compute network domain.
type NetworkDomain struct {
	// The network domain Id.
	ID string `json:"id"`

	// The network domain name.
	Name string `json:"name"`

	// The network domain description.
	Description string `json:"description"`

	// The network domain type.
	Type string `json:"type"`

	// The network domain's NAT IPv4 address.
	NatIPv4Address string `json:"snatIpv4Address"`

	// The network domain's outside transit IPv4 subnet.
	OutsideTransitVLANIPv4Subnet IPv4Range `json:"outsideTransitVlanIpv4Subnet"`

	// The network domain's creation timestamp.
	CreateTime string `json:"createTime"`

	// The network domain's current state.
	State string `json:"state"`

	// The network domain's current progress (if any).
	Progress string `json:"progress"`

	// The Id of the data centre in which the network domain is located.
	DatacenterID string `json:"datacenterId"`
}

// GetID returns the network domain's Id.
func (domain *NetworkDomain) GetID() string {
	return domain.ID
}

// GetResourceType returns the network domain's resource type.
func (domain *NetworkDomain) GetResourceType() ResourceType {
	return ResourceTypeNetworkDomain
}

// GetName returns the network domain's name.
func (domain *NetworkDomain) GetName() string {
	return domain.Name
}

// GetState returns the network domain's current state.
func (domain *NetworkDomain) GetState() string {
	return domain.State
}

// IsDeleted determines whether the network domain has been deleted (is nil).
func (domain *NetworkDomain) IsDeleted() bool {
	return domain == nil
}

var _ Resource = &NetworkDomain{}

// ToEntityReference creates an EntityReference representing the NetworkDomain.
func (domain *NetworkDomain) ToEntityReference() EntityReference {
	return EntityReference{
		ID:   domain.ID,
		Name: domain.Name,
	}
}

var _ NamedEntity = &NetworkDomain{}

// NetworkDomains represents the response to a "List Network Domains" API call.
type NetworkDomains struct {
	// The current page of network domains.
	Domains []NetworkDomain `json:"networkDomain"`

	PagedResult
}

// Request body for deploying a compute network domain.
type deployNetworkDomain struct {
	// The network domain name.
	Name string `json:"name"`

	// The network domain description.
	Description string `json:"description"`

	// The network domain type.
	Type string `json:"type"`

	// The Id of the data centre in which the network domain is located.
	DatacenterID string `json:"datacenterId"`
}

// Request body for editing a compute network domain.
type editNetworkDomain struct {
	// The network domain ID.
	ID string `json:"id"`

	// The network domain name (optional).
	Name *string `json:"name,omitempty"`

	// The network domain description (optional).
	Description *string `json:"description,omitempty"`

	// The network domain type (optional).
	Type *string `json:"type,omitempty"`
}

// Request body for deleting a compute network domain.
type deleteNetworkDomain struct {
	// The network domain ID.
	ID string `json:"id"`
}

// ListNetworkDomains retrieves a list of all network domains.
// TODO: Support filtering and sorting.
func (client *Client) ListNetworkDomains(paging *Paging) (domains *NetworkDomains, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return nil, err
	}

	requestURI := fmt.Sprintf("%s/network/networkDomain?%s",
		url.QueryEscape(organizationID),
		paging.EnsurePaging().toQueryParameters(),
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

		return nil, apiResponse.ToError("Request failed with status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	domains = &NetworkDomains{}
	err = json.Unmarshal(responseBody, domains)
	if err != nil {
		return nil, err
	}

	return domains, nil
}

// GetNetworkDomain retrieves the network domain with the specified Id.
// id is the Id of the network domain to retrieve.
// Returns nil if no network domain is found with the specified Id.
func (client *Client) GetNetworkDomain(id string) (domain *NetworkDomain, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return nil, err
	}

	requestURI := fmt.Sprintf("%s/network/networkDomain/%s",
		url.QueryEscape(organizationID),
		url.QueryEscape(id),
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

		if apiResponse.ResponseCode == ResponseCodeResourceNotFound {
			return nil, nil // Not an error, but was not found.
		}

		return nil, apiResponse.ToError("Request to retrieve network domain failed with status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	domain = &NetworkDomain{}
	err = json.Unmarshal(responseBody, domain)
	if err != nil {
		return nil, err
	}

	return domain, nil
}

// GetNetworkDomainByName retrieves the network domain (if any) with the specified name in the specified data centre.
func (client *Client) GetNetworkDomainByName(name string, dataCenterID string) (domain *NetworkDomain, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return nil, err
	}

	requestURI := fmt.Sprintf("%s/network/networkDomain?name=%s&datacenterId=%s",
		url.QueryEscape(organizationID),
		url.QueryEscape(name),
		url.QueryEscape(dataCenterID),
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

		return nil, apiResponse.ToError("Request failed with status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	domains := &NetworkDomains{}
	err = json.Unmarshal(responseBody, domains)
	if err != nil {
		return nil, err
	}
	if domains.IsEmpty() {
		return nil, nil // No matching network domain was found.
	}

	if len(domains.Domains) != 1 {
		return nil, fmt.Errorf("found multiple network domains (%d) named '%s' in data centre '%s'", len(domains.Domains), name, dataCenterID)
	}

	return &domains.Domains[0], nil
}

// DeployNetworkDomain deploys a new network domain.
// Returns the Id of the new network domain.
func (client *Client) DeployNetworkDomain(name string, description string, plan string, datacenter string) (networkDomainID string, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return "", err
	}

	requestURI := fmt.Sprintf("%s/network/deployNetworkDomain",
		url.QueryEscape(organizationID),
	)
	request, err := client.newRequestV24(requestURI, http.MethodPost, &deployNetworkDomain{
		Name:         name,
		Description:  description,
		Type:         plan,
		DatacenterID: datacenter,
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
		return "", apiResponse.ToError("Request to deploy network domain '%s' failed with status code %d (%s): %s", name, statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	// Expected: "info" { "name": "networkDomainId", "value": "the-Id-of-the-new-network-domain" }
	networkDomainIDMessage := apiResponse.GetFieldMessage("networkDomainId")
	if networkDomainIDMessage == nil {
		return "", apiResponse.ToError("Received an unexpected response (missing 'networkDomainId') with status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	return *networkDomainIDMessage, nil
}

// EditNetworkDomain updates an existing network domain.
// Pass an empty string for any field to retain its existing value.
// Returns an error if the operation was not successful.
func (client *Client) EditNetworkDomain(id string, name *string, description *string, plan *string) (err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return err
	}

	requestURI := fmt.Sprintf("%s/network/editNetworkDomain",
		url.QueryEscape(organizationID),
	)
	request, err := client.newRequestV24(requestURI, http.MethodPost, &editNetworkDomain{
		ID:          id,
		Name:        name,
		Description: description,
		Type:        plan,
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

// DeleteNetworkDomain deletes an existing network domain.
// Returns an error if the operation was not successful.
func (client *Client) DeleteNetworkDomain(id string) (err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return err
	}

	requestURI := fmt.Sprintf("%s/network/deleteNetworkDomain",
		url.QueryEscape(organizationID),
	)
	request, err := client.newRequestV24(requestURI, http.MethodPost, &deleteNetworkDomain{id})
	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return err
	}

	apiResponse, err := readAPIResponseAsJSON(responseBody, statusCode)
	if err != nil {
		return err
	}

	if apiResponse.ResponseCode != ResponseCodeInProgress {
		return apiResponse.ToError("Request to delete network domain failed with unexpected status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	return nil
}
