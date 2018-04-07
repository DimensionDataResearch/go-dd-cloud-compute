package compute

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// PersistenceProfile represents a load-balancer persistence (stickiness) profile.
type PersistenceProfile struct {
	ID                      string `json:"id"`
	Name                    string `json:"name"`
	IsFallbackCompatible    bool   `json:"fallbackCompatible"`
	VirtualListenerType     string `json:"virtualListenerType"`
	VirtualListenerProtocol string `json:"virtualListenerProtocol"`
}

// GetID retrieves the persistence profile's ID.
func (profile *PersistenceProfile) GetID() string {
	return profile.ID
}

// GetName retrieves the persistence profile's name.
func (profile *PersistenceProfile) GetName() string {
	return profile.Name
}

// ToEntityReference creates an EntityReference representing the PersistenceProfile.
func (profile *PersistenceProfile) ToEntityReference() EntityReference {
	return EntityReference{
		ID:   profile.ID,
		Name: profile.Name,
	}
}

var _ NamedEntity = &PersistenceProfile{}

// PersistenceProfiles represents a page of PersistenceProfile results.
type PersistenceProfiles struct {
	Items []PersistenceProfile `json:"defaultPersistenceProfile"`

	PagedResult
}

// ListDefaultPersistenceProfiles retrieves a list of all default load-balancing persistence profiles in the specified network domain.
func (client *Client) ListDefaultPersistenceProfiles(networkDomainID string, paging *Paging) (persistenceProfiles *PersistenceProfiles, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return nil, err
	}

	requestURI := fmt.Sprintf("%s/networkDomainVip/defaultPersistenceProfile?networkDomainId=%s&%s",
		url.QueryEscape(organizationID),
		url.QueryEscape(networkDomainID),
		paging.EnsurePaging().toQueryParameters(),
	)
	request, err := client.newRequestV26(requestURI, http.MethodGet, nil)
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

		return nil, apiResponse.ToError("Request to list default persistence profiles in network domain '%s' failed with status code %d (%s): %s", networkDomainID, statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	persistenceProfiles = &PersistenceProfiles{}
	err = json.Unmarshal(responseBody, persistenceProfiles)
	if err != nil {
		return nil, err
	}

	return persistenceProfiles, nil
}
