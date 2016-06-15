package compute

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// FirewallRule represents a firewall rule.
type FirewallRule struct {
	ID              string            `json:"id"`
	Name            string            `json:"name"`
	Action          string            `json:"action"`
	IPVersion       string            `json:"ipVersion"`
	Protocol        string            `json:"protocol"`
	Source          FirewallRuleScope `json:"source"`
	Destination     FirewallRuleScope `json:"destination"`
	Enabled         bool              `json:"enabled"`
	State           string            `json:"state"`
	NetworkDomainID string            `json:"networkDomainId"`
	DataCenterID    string            `json:"datacenterId"`
	RuleType        string            `json:"ruleType"`
}

// FirewallRuleScope represents a scope (IP and / or port) for firewall configuration (source or destination).
type FirewallRuleScope struct {
	IPAddress   *FirewallRuleIPAddress `json:"ip,omitempty"`
	AddressList *EntitySummary         `json:"ipAddressList,omitempty"`
	Port        *FirewallRulePort      `json:"port,omitempty"`
}

// FirewallRuleIPAddress represents represents an IP address for firewall configuration.
type FirewallRuleIPAddress struct {
	Address    string `json:"address"`
	PrefixSize *int   `json:"PrefixSize,omitempty"`
}

// FirewallRulePort represents a firewall port configuration.
type FirewallRulePort struct {
	Begin int  `json:"begin"`
	End   *int `json:"end"`
}

// FirewallRules represents a page of FirewallRule results.
type FirewallRules struct {
	Rules []FirewallRule `json:"firewallRule"`

	PagedResult
}

// FirewallRuleConfiguration represents the configuration for a new firewall rule.
type FirewallRuleConfiguration struct {
	Name            string                `json:"name"`
	Action          string                `json:"action"`
	Enabled         bool                  `json:"enabled"`
	Placement       FirewallRulePlacement `json:"placement"`
	IPVersion       string                `json:"ipVersion"`
	Protocol        string                `json:"protocol"`
	Source          FirewallRuleScope     `json:"source"`
	Destination     FirewallRuleScope     `json:"destination"`
	NetworkDomainID string                `json:"networkDomainId"`
	DataCenterID    string                `json:"datacenterId"`
}

// PlaceFirst modifies the configuration so that the firewall rule will be placed in the first available position.
func (configuration *FirewallRuleConfiguration) PlaceFirst() {
	configuration.Placement = FirewallRulePlacement{
		Position: "FIRST",
	}
}

// PlaceBefore modifies the configuration so that the firewall rule will be placed before the specified rule.
func (configuration *FirewallRuleConfiguration) PlaceBefore(beforeRuleName string) {
	configuration.Placement = FirewallRulePlacement{
		Position:           "BEFORE",
		RelativeToRuleName: &beforeRuleName,
	}
}

// PlaceAfter modifies the configuration so that the firewall rule will be placed after the specified rule.
func (configuration *FirewallRuleConfiguration) PlaceAfter(afterRuleName string) {
	configuration.Placement = FirewallRulePlacement{
		Position:           "AFTER",
		RelativeToRuleName: &afterRuleName,
	}
}



// MatchAnySource modifies the configuration so that the firewall rule will match any combination of source IP and port.
func (configuration *FirewallRuleConfiguration) MatchAnySource() {
	configuration.Source = FirewallRuleScope{
		IPAddress: &FirewallRuleIPAddress{
			Address: "ANY",
		},
		Port: nil,
	}
}

// MatchSourceAddressAndPort modifies the configuration so that the firewall rule will match a specific source IP address (and, optionall, port).
func (configuration *FirewallRuleConfiguration) MatchSourceAddressAndPort(address string, port *int) {
	sourceScope := &FirewallRuleScope{
		IPAddress: &FirewallRuleIPAddress{
			Address: address,
		},
	}
	if port != nil {
		sourceScope.Port = &FirewallRulePort{
			Begin: *port,
		}
	}
	configuration.Source = *sourceScope
}

// MatchDestinationAddressAndPort modifies the configuration so that the firewall rule will match a specific destination IP address (and, optionall, port).
func (configuration *FirewallRuleConfiguration) MatchDestinationAddressAndPort(address string, port *int) {
	destinationScope := &FirewallRuleScope{
		IPAddress: &FirewallRuleIPAddress{
			Address: address,
		},
	}
	if port != nil {
		destinationScope.Port = &FirewallRulePort{
			Begin: *port,
		}
	}
	configuration.Destination = *destinationScope
}

// MatchSourceAddressListAndPort modifies the configuration so that the firewall rule will match a specific source IP address list (and, optionall, port).
func (configuration *FirewallRuleConfiguration) MatchSourceAddressListAndPort(addressListID string, port *int) {
	sourceScope := &FirewallRuleScope{
		AddressList: &EntitySummary{
			ID: addressListID,
		},
	}
	if port != nil {
		sourceScope.Port = &FirewallRulePort{
			Begin: *port,
		}
	}
	configuration.Source = *sourceScope
}

// MatchDestinationAddressListAndPort modifies the configuration so that the firewall rule will match a specific destination IP address list (and, optionall, port).
func (configuration *FirewallRuleConfiguration) MatchDestinationAddressListAndPort(addressListID string, port *int) {
	destinationScope := &FirewallRuleScope{
		AddressList: &EntitySummary{
			ID: addressListID,
		},
	}
	if port != nil {
		destinationScope.Port = &FirewallRulePort{
			Begin: *port,
		}
	}
	configuration.Destination = *destinationScope
}

// FirewallRulePlacement describes the placement for a firewall rule.
type FirewallRulePlacement struct {
	Position           string  `json:"position"`
	RelativeToRuleName *string `json:"relativeToRule,omitempty"`
}

// GetFirewallRule retrieves the Firewall rule with the specified Id.
// Returns nil if no Firewall rule is found with the specified Id.
func (client *Client) GetFirewallRule(id string) (rule *FirewallRule, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return nil, err
	}

	requestURI := fmt.Sprintf("%s/network/firewallRule/%s", organizationID, id)
	request, err := client.newRequestV22(requestURI, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		var apiResponse *APIResponse

		apiResponse, err = readAPIResponseAsJSON(responseBody, statusCode)
		if err != nil {
			return nil, err
		}

		if apiResponse.ResponseCode == ResponseCodeResourceNotFound {
			return nil, nil // Not an error, but was not found.
		}

		return nil, apiResponse.ToError("Request to retrieve firewall rule failed with status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	rule = &FirewallRule{}
	err = json.Unmarshal(responseBody, rule)
	if err != nil {
		return nil, err
	}

	return rule, nil
}

// ListFirewallRules lists all firewall rules that apply to the specified network domain.
func (client *Client) ListFirewallRules(networkDomainID string) (rules *FirewallRules, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return nil, err
	}

	requestURI := fmt.Sprintf("%s/network/firewallRule?networkDomainId=%s", organizationID, networkDomainID)
	request, err := client.newRequestV22(requestURI, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		var apiResponse *APIResponse

		apiResponse, err = readAPIResponseAsJSON(responseBody, statusCode)
		if err != nil {
			return nil, err
		}

		return nil, apiResponse.ToError("Request to list firewall rules for network domain '%s' failed with status code %d (%s): %s", networkDomainID, statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	rules = &FirewallRules{}
	err = json.Unmarshal(responseBody, rules)

	return rules, err
}

// CreateFirewallRule creates a new firewall rule.
func (client *Client) CreateFirewallRule(configuration FirewallRuleConfiguration) (firewallRuleID string, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return "", err
	}

	requestURI := fmt.Sprintf("%s/network/createFirewallRule", organizationID)
	request, err := client.newRequestV22(requestURI, http.MethodPost, &configuration)
	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return "", err
	}

	apiResponse, err := readAPIResponseAsJSON(responseBody, statusCode)
	if err != nil {
		return "", err
	}

	if apiResponse.ResponseCode != ResponseCodeOK {
		return "", apiResponse.ToError("Request to create firewall rule in network domain '%s' failed with unexpected status code %d (%s): %s", configuration.NetworkDomainID, statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	// Expected: "info" { "name": "firewallRuleId", "value": "the-Id-of-the-new-firewall-rule" }
	if len(apiResponse.FieldMessages) != 1 || apiResponse.FieldMessages[0].FieldName != "firewallRuleId" {
		return "", apiResponse.ToError("Received an unexpected response (missing 'firewallRuleId') with status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	return apiResponse.FieldMessages[0].Message, nil
}
