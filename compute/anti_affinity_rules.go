package compute

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
)

// ServerAntiAffinityRule represents an anti-affinity rule between 2 servers.
type ServerAntiAffinityRule struct {
	// The anti-affinity rule Id.
	ID string `json:"id"`

	// The 2 servers that the rule relates to.
	//
	// Only ever contains exactly 2 servers.
	//
	// This is only declared as an array because that's what the CloudControl API returns.
	Servers []ServerSummary `json:"serverSummary"`

	// The network domain's current state.
	State string `json:"state"`

	// The network domain's creation timestamp.
	CreateTime string `json:"created"`

	// The Id of the data centre in which the network domain is located.
	DatacenterID string `json:"datacenterId"`
}

// GetID returns the server anti-affinity rule's Id.
func (rule *ServerAntiAffinityRule) GetID() string {
	return rule.ID
}

// GetResourceType returns the server anti-affinity rule's resource type.
func (rule *ServerAntiAffinityRule) GetResourceType() ResourceType {
	return ResourceTypeServerAntiAffinityRule
}

// GetName returns the server anti-affinity rule's name.
func (rule *ServerAntiAffinityRule) GetName() string {
	return rule.ID
}

// GetState returns the server anti-affinity rule's current state.
func (rule *ServerAntiAffinityRule) GetState() string {
	return rule.State
}

// IsDeleted determines whether the server anti-affinity rule has been deleted (is nil).
func (rule *ServerAntiAffinityRule) IsDeleted() bool {
	return rule == nil
}

// ToEntityReference creates an EntityReference representing the CustomerImage.
func (rule *ServerAntiAffinityRule) ToEntityReference() EntityReference {
	name := ""
	if len(rule.Servers) == 2 {
		name = fmt.Sprintf("%s/%s",
			rule.Servers[0].Name,
			rule.Servers[1].Name,
		)
	}

	return EntityReference{
		ID:   rule.ID,
		Name: name,
	}
}

var _ Resource = &ServerAntiAffinityRule{}

// ServerAntiAffinityRules represents a page of ServerAntiAffinityRule results.
type ServerAntiAffinityRules struct {
	Items []ServerAntiAffinityRule `json:"antiAffinityRule"`

	PagedResult
}

// Request body when creating a new anti-affinity rule.
type newServerAntiAffinityRule struct {
	XMLName xml.Name `xml:"http://oec.api.opsource.net/schemas/server NewAntiAffinityRule"`

	// The Ids of the servers to which the rule relates.
	// Each rule can only apply to exactly 2 servers; we only use an array here because CloudControl (bizarrely) uses the same element name for both server Ids.
	ServerIds []string `xml:"serverId"`
}

// GetServerAntiAffinityRule retrieves the specified server anti-affinity rule (in the specified network domain).
func (client *Client) GetServerAntiAffinityRule(ruleID string, networkDomainID string) (rule *ServerAntiAffinityRule, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return nil, err
	}

	requestURI := fmt.Sprintf("%s/server/antiAffinityRule?id=%s&networkDomainId=%s",
		url.QueryEscape(organizationID),
		url.QueryEscape(ruleID),
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

	rules := &ServerAntiAffinityRules{}
	err = json.Unmarshal(responseBody, rules)
	if err != nil {
		return nil, err
	}

	if rules.IsEmpty() {
		return nil, nil // Rule not found
	}

	return &rules.Items[0], nil
}

// ListServerAntiAffinityRules lists the server anti-affinity rules in a network domain.
func (client *Client) ListServerAntiAffinityRules(networkDomainID string, paging *Paging) (rules *ServerAntiAffinityRules, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return nil, err
	}

	requestURI := fmt.Sprintf("%s/server/antiAffinityRule?networkDomainId=%s&%s",
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

		return nil, apiResponse.ToError("Request failed with status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	rules = &ServerAntiAffinityRules{}
	err = json.Unmarshal(responseBody, rules)
	if err != nil {
		return nil, err
	}

	return rules, nil
}

// CreateServerAntiAffinityRule creates an anti-affinity rule for the 2 specified servers.
// server1Id is the Id of the first server.
// server2Id is the Id of the second server.
//
// Returns the Id of the new anti-affinity rule.
func (client *Client) CreateServerAntiAffinityRule(server1Id string, server2Id string) (ruleID string, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return "", err
	}

	requestURI := fmt.Sprintf("%s/antiAffinityRule",
		url.QueryEscape(organizationID),
	)
	request, err := client.newRequestV1(requestURI, http.MethodPost, &newServerAntiAffinityRule{
		ServerIds: []string{
			server1Id,
			server2Id,
		},
	})
	if err != nil {
		return "", err
	}

	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return "", err
	}

	if statusCode != http.StatusOK {
		var apiResponse *APIResponseV1

		apiResponse, err = readAPIResponseV1(responseBody, statusCode)
		if err != nil {
			return "", err
		}

		return "", apiResponse.ToError("Request failed with status code %d (%s): %s", statusCode, apiResponse.ResultCode, apiResponse.Message)
	}

	apiResponse := &APIResponseV1{}
	err = xml.Unmarshal(responseBody, apiResponse)
	if err != nil {
		return "", err
	}

	newRuleID := apiResponse.GetAdditionalInformation("antiaffinityrule.id")
	if newRuleID == nil {
		return "", apiResponse.ToError("Invalid response (missing 'antiaffinityrule.id')")
	}

	return *newRuleID, nil
}

// DeleteServerAntiAffinityRule deletes the specified server anti-affinity rule.
func (client *Client) DeleteServerAntiAffinityRule(ruleID string, networkDomainID string) error {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return err
	}

	requestURI := fmt.Sprintf("%s/antiAffinityRule/%s?delete",
		url.QueryEscape(organizationID),
		url.QueryEscape(ruleID),
	)
	request, err := client.newRequestV1(requestURI, http.MethodGet, nil)
	if err != nil {
		return err
	}

	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return err
	}

	apiResponse, err := readAPIResponseV1(responseBody, statusCode)
	if err != nil {
		return err
	}

	if apiResponse.Result != ResultSuccess {
		return apiResponse.ToError("Request failed with status code %d (%s): %s", statusCode, apiResponse.ResultCode, apiResponse.Message)
	}

	return nil
}
