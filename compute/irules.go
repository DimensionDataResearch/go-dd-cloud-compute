package compute

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// IRule represents a load-balancer iRule.
type IRule struct {
	ID                      string `json:"id"`
	Name                    string `json:"name"`
	VirtualListenerType     string `json:"virtualListenerType"`
	VirtualListenerProtocol string `json:"virtualListenerProtocol"`
}

// GetID retrieves the iRule's ID.
func (iRule *IRule) GetID() string {
	return iRule.ID
}

// GetName retrieves the iRule's name.
func (iRule *IRule) GetName() string {
	return iRule.Name
}

// ToEntityReference creates an EntityReference representing the IRule.
func (iRule *IRule) ToEntityReference() EntityReference {
	return EntityReference{
		ID:   iRule.ID,
		Name: iRule.Name,
	}
}

var _ NamedEntity = &IRule{}

// IRules represents a page of IRule results.
type IRules struct {
	Items []IRule `json:"defaultIRule"`

	PagedResult
}

// ListDefaultIRules retrieves a list of all default load-balancing iRules in the specified network domain.
func (client *Client) ListDefaultIRules(networkDomainID string, paging *Paging) (irules *IRules, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return nil, err
	}

	requestURI := fmt.Sprintf("%s/networkDomainVip/defaultIrule?networkDomainId=%s&%s",
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

		return nil, apiResponse.ToError("Request to list default iRules in network domain '%s' failed with status code %d (%s): %s", networkDomainID, statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	irules = &IRules{}
	err = json.Unmarshal(responseBody, irules)
	if err != nil {
		return nil, err
	}

	return irules, nil
}
