package compute

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// Datacenter represents an MCP datacenter.
type Datacenter struct {
	// The datacenter Id.
	ID string `json:"id"`

	// The datacenter type (for display purposes only).
	Type string `json:"type"`

	// The datacenter display name.
	DisplayName string `json:"displayName"`

	// The name of the city the datacenter is located in.
	City string `json:"city"`

	// The name of the state the datacenter is located in.
	State string `json:"state"`

	// The name of the country the datacenter is located in.
	Country string `json:"country"`

	// The URL of the datacenter's administrative SSL VPN.
	VPNURL string `json:"vpnUrl"`

	// The name of the FTPS host used to upload / download OVF packages to / from the datacenter.
	FTPSHost string `json:"ftpsHost"`

	// The datacenter's network configuration.
	Networking DatacenterNetworking `json:"networking"`
}

// DatacenterNetworking represents the networking configuration for an MCP datacenter.
type DatacenterNetworking struct {
	// The networking infrastructure type of the data center for programmatic use.
	//
	// "1" means MCP 1.0
	// "2" means MCP 2.0
	Type string `json:"type"`

	// Indicates whether the networking infrastructure is under maintenance.
	MaintenanceStatus string `json:"maintenanceStatus"`
}

// Datacenters represents the response to a "List Datacenters" API call.
type Datacenters struct {
	// The current page of datacenters.
	Items []Datacenter `json:"datacenter"`

	PagedResult
}

// ListDatacenters retrieves a list of all datacenters.
// TODO: Support filtering and sorting.
func (client *Client) ListDatacenters(paging *Paging) (datacenters *Datacenters, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return nil, err
	}

	requestURI := fmt.Sprintf("%s/infrastructure/datacenter?%s",
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

	datacenters = &Datacenters{}
	err = json.Unmarshal(responseBody, datacenters)
	if err != nil {
		return nil, err
	}

	return datacenters, nil
}

// GetDatacenter retrieves the datacenter with the specified Id.
// id is the Id of the datacenter to retrieve.
// Returns nil if no datacenter is found with the specified Id.
func (client *Client) GetDatacenter(id string) (datacenter *Datacenter, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return nil, err
	}

	requestURI := fmt.Sprintf("%s/infrastructure/datacenter?id=%s",
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

		return nil, apiResponse.ToError("Request to retrieve datacenter '%s' failed with status code %d (%s): %s", id, statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	datacenters := &Datacenters{}
	err = json.Unmarshal(responseBody, datacenters)
	if err != nil {
		return nil, err
	}

	if datacenters.IsEmpty() {
		return nil, nil
	}

	return &datacenters.Items[0], nil
}
