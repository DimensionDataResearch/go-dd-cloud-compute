package compute

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// PortList represents a port list.
type PortList struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Ports       []PortListEntry   `json:"port"`
	ChildLists  []EntityReference `json:"childPortList"`
	State       string            `json:"state"`
	CreateTime  string            `json:"createTime"`
}

// BuildEditRequest creates an EditPortList using the existing ports and child list references in the port list.
func (portList *PortList) BuildEditRequest() EditPortList {
	edit := &EditPortList{
		Description:  portList.Description,
		Ports:        portList.Ports,
		ChildListIDs: make([]string, len(portList.ChildLists)),
	}
	for index, childList := range portList.ChildLists {
		edit.ChildListIDs[index] = childList.ID
	}

	return *edit
}

// PortListEntry represents an entry in a port list.
type PortListEntry struct {
	Begin int  `json:"begin"`
	End   *int `json:"end,omitempty"`
}

// PortLists represents a page of PortList results.
type PortLists struct {
	PortLists []PortList `json:"portList"`

	PagedResult
}

// Request body for creating a port list.
type createPortList struct {
	Name            string          `json:"name"`
	Description     string          `json:"description"`
	NetworkDomainID string          `json:"networkDomainId"`
	Ports           []PortListEntry `json:"port"`
	ChildListIDs    []string        `json:"childPortListId"`
}

// EditPortList represents the request body for editing a port list.
type EditPortList struct {
	ID           string          `json:"id"`
	Description  string          `json:"description"`
	Ports        []PortListEntry `json:"port"`
	ChildListIDs []string        `json:"childPortListId"`
}

// Request body for deleting a port list.
type deletePortList struct {
	ID string `json:"id"`
}

// GetPortList retrieves the port list with the specified Id.
// id is the Id of the port list to retrieve.
// Returns nil if no portList is found with the specified Id.
func (client *Client) GetPortList(id string) (portList *PortList, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return nil, err
	}

	requestURI := fmt.Sprintf("%s/network/portList/%s", organizationID, id)
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

		return nil, apiResponse.ToError("Request to retrieve port list failed with status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	portList = &PortList{}
	err = json.Unmarshal(responseBody, portList)

	return portList, err
}

// ListPortLists retrieves all port lists associated with the specified network domain.
func (client *Client) ListPortLists(networkDomainID string) (portLists *PortLists, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return nil, err
	}

	requestURI := fmt.Sprintf("%s/network/portList?networkDomainId=%s", organizationID, networkDomainID)
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

		return nil, apiResponse.ToError("Request to list port lists failed with status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	portLists = &PortLists{}
	err = json.Unmarshal(responseBody, portLists)

	return portLists, err
}

// CreatePortList creates a new port list.
// Returns the Id of the new port list.
//
// This operation is synchronous.
func (client *Client) CreatePortList(name string, description string, ipVersion string, networkDomainID string, ports []PortListEntry, childListIDs []string) (portListID string, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return "", err
	}

	requestURI := fmt.Sprintf("%s/network/createPortList", organizationID)
	request, err := client.newRequestV22(requestURI, http.MethodPost, &createPortList{
		Name:            name,
		Description:     description,
		Ports:           ports,
		ChildListIDs:    childListIDs,
		NetworkDomainID: networkDomainID,
	})
	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return "", err
	}

	apiResponse, err := readAPIResponseAsJSON(responseBody, statusCode)
	if err != nil {
		return "", err
	}

	if apiResponse.ResponseCode != ResponseCodeOK {
		return "", apiResponse.ToError("Request to create port list '%s' failed with status code %d (%s): %s", name, statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	// Expected: "info" { "name": "portListId", "value": "the-Id-of-the-new-IP-port-list" }
	portListIDMessage := apiResponse.GetFieldMessage("portListId")
	if portListIDMessage == nil {
		return "", apiResponse.ToError("Received an unexpected response (missing 'portListId') with status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	return *portListIDMessage, nil
}

// EditPortList updates the configuration for a port list.
//
// Note that this operation is not additive; it *replaces* the configuration for the port list.
// You can PortList.BuildEditRequest() to create an EditPortList request that copies the current state of the PortList (and then apply customisations).
//
// This operation is synchronous.
func (client *Client) EditPortList(id string, edit EditPortList) error {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return err
	}

	requestURI := fmt.Sprintf("%s/network/editPortList", organizationID)
	request, err := client.newRequestV22(requestURI, http.MethodPost, edit)
	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return err
	}

	apiResponse, err := readAPIResponseAsJSON(responseBody, statusCode)
	if err != nil {
		return err
	}

	if apiResponse.ResponseCode != ResponseCodeOK {
		return apiResponse.ToError("Request to edit port list failed with unexpected status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	return nil
}

// DeletePortList deletes an existing port list.
// Returns an error if the operation was not successful.
//
// This operation is synchronous.
func (client *Client) DeletePortList(id string) (err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return err
	}

	requestURI := fmt.Sprintf("%s/network/deletePortList", organizationID)
	request, err := client.newRequestV22(requestURI, http.MethodPost, &deletePortList{id})
	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return err
	}

	apiResponse, err := readAPIResponseAsJSON(responseBody, statusCode)
	if err != nil {
		return err
	}

	if apiResponse.ResponseCode != ResponseCodeOK {
		return apiResponse.ToError("Request to delete port list failed with unexpected status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	return nil
}
