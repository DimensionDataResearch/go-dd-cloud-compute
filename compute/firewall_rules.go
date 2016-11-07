package compute

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	// FirewallRuleActionAccept indicates a firewall rule that, if it matches, will accept the packet and stop processing further rules.
	FirewallRuleActionAccept = "ACCEPT_DECISIVELY"

	// FirewallRuleActionDrop indicates a firewale rule that, if it matches, will drop the packet.
	FirewallRuleActionDrop = "DROP"

	// FirewallRuleIPVersion4 indicates a firewall rule that targets IPv4
	FirewallRuleIPVersion4 = "IPv4"

	// FirewallRuleIPVersion6 indicates a firewale rule that targets IPv6
	FirewallRuleIPVersion6 = "IPv6"

	// FirewallRuleProtocolIP indicates a firewall rule that targets the Internet Protocol (IP)
	FirewallRuleProtocolIP = "IP"

	// FirewallRuleProtocolTCP indicates a firewall rule that targets the Transmission Control Protocol (TCP)
	FirewallRuleProtocolTCP = "TCP"

	// FirewallRuleProtocolICMP indicates a firewall rule that targets the Internet Control Message Protocol (ICMP)
	FirewallRuleProtocolICMP = "ICMP"

	// FirewallRuleMatchAny indicates a firewall rule value that matches any other value in the same scope.
	FirewallRuleMatchAny = "ANY"
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

// GetID returns the firewall rule's Id.
func (rule *FirewallRule) GetID() string {
	return rule.ID
}

// GetResourceType returns the firewall rule's resource type.
func (rule *FirewallRule) GetResourceType() ResourceType {
	return ResourceTypeFirewallRule
}

// GetName returns the firewall rule's name.
func (rule *FirewallRule) GetName() string {
	return rule.Name
}

// GetState returns the firewall rule's current state.
func (rule *FirewallRule) GetState() string {
	return rule.State
}

// IsDeleted determines whether the firewall rule has been deleted (is nil).
func (rule *FirewallRule) IsDeleted() bool {
	return rule == nil
}

var _ Resource = &FirewallRule{}

// FirewallRuleScope represents a scope (IP and / or port) for firewall configuration (source or destination).
type FirewallRuleScope struct {
	IPAddress     *FirewallRuleIPAddress `json:"ip,omitempty"`
	AddressList   *EntityReference       `json:"ipAddressList,omitempty"`
	AddressListID *string                `json:"ipAddressListId,omitempty"`
	Port          *FirewallRulePort      `json:"port,omitempty"`
	PortListID    *string                `json:"portListId,omitempty"`
}

// IsScopeHost determines whether the firewall rule scope matches a host.
func (scope *FirewallRuleScope) IsScopeHost() bool {
	return scope.IPAddress != nil && scope.IPAddress.PrefixSize == nil
}

// IsScopeNetwork determines whether the firewall rule scope matches a network.
func (scope *FirewallRuleScope) IsScopeNetwork() bool {
	return scope.IPAddress != nil && scope.IPAddress.PrefixSize != nil
}

// IsScopePort determines whether the firewall rule scope matches a single port.
func (scope *FirewallRuleScope) IsScopePort() bool {
	return scope.Port != nil && scope.Port.End == nil
}

// IsScopePortRange determines whether the firewall rule scope matches a port range.
func (scope *FirewallRuleScope) IsScopePortRange() bool {
	return scope.Port != nil && scope.Port.End != nil
}

// IsScopeAddressList determines whether the firewall rule scope matches an IP address list.
func (scope *FirewallRuleScope) IsScopeAddressList() bool {
	return scope.AddressList != nil || scope.AddressListID != nil
}

// IsScopeAny determines whether the firewall rule scope matches anything (i.e. the rule is unscoped).
func (scope *FirewallRuleScope) IsScopeAny() bool {
	return scope.IPAddress == nil && scope.AddressList == nil && scope.Port == nil
}

// Diff captures the differences (if any) between a FirewallRuleScope and another FirewallRuleScope.
func (scope FirewallRuleScope) Diff(other FirewallRuleScope) (differences []string) {
	if scope.IsScopeHost() {
		if other.IsScopeHost() {
			if scope.IPAddress.Address != other.IPAddress.Address {
				differences = append(differences, fmt.Sprintf(
					"target hosts do not match ('%s' vs '%s')",
					scope.IPAddress.Address,
					other.IPAddress.Address,
				))
			}
		} else if other.IsScopeNetwork() {
			differences = append(differences, "host scope vs network scope")
		} else if other.IsScopeAddressList() {
			differences = append(differences, "host scope vs address list scope")
		} else {
			differences = append(differences, "host scope vs unknown scope")
		}
	} else if scope.IsScopeNetwork() {
		if other.IsScopeNetwork() {
			scopeNetwork := fmt.Sprintf("%s/%d",
				scope.IPAddress.Address,
				*scope.IPAddress.PrefixSize,
			)
			otherNetwork := fmt.Sprintf("%s/%d",
				other.IPAddress.Address,
				*other.IPAddress.PrefixSize,
			)

			if scopeNetwork != otherNetwork {
				differences = append(differences, fmt.Sprintf(
					"target networks do not match ('%s' vs '%s')",
					scopeNetwork,
					otherNetwork,
				))
			}
		} else if other.IsScopeHost() {
			differences = append(differences, "network scope vs host scope")
		} else if other.IsScopeAddressList() {
			differences = append(differences, "network scope vs address list scope")
		} else {
			differences = append(differences, "network scope vs unknown scope")
		}
	} else if scope.IsScopeAddressList() {
		if other.IsScopeAddressList() {
			addressListID := scope.AddressListID
			if addressListID == nil {
				addressListID = &scope.AddressList.ID
			}

			otherAddressListID := other.AddressListID
			if otherAddressListID == nil {
				otherAddressListID = &other.AddressList.ID
			}

			if addressListID != otherAddressListID {
				differences = append(differences, fmt.Sprintf(
					"address lists do not match ('%s' vs '%s')",
					scope.AddressList.ID,
					other.AddressList.ID,
				))
			}
		} else if other.IsScopeHost() {
			differences = append(differences, "address list scope vs host scope")
		} else if other.IsScopeNetwork() {
			differences = append(differences, "address list scope vs network scope")
		} else {
			differences = append(differences, "address list scope vs unknown scope")
		}
	}

	if scope.IsScopePort() {
		if other.IsScopePort() {
			if scope.Port.Begin != other.Port.Begin {
				differences = append(differences, fmt.Sprintf(
					"ports do not match (%d vs %d)",
					scope.Port.Begin,
					scope.Port.End,
				))
			}
		} else if other.IsScopePortRange() {
			differences = append(differences, "port scope vs port-range scope")
		} else {
			differences = append(differences, "port scope vs no scope")
		}
	} else if scope.IsScopePortRange() {
		if other.IsScopePortRange() {
			scopeRange := fmt.Sprintf("%d-%d",
				scope.Port.Begin,
				*scope.Port.End,
			)
			otherRange := fmt.Sprintf("%d-%d",
				other.Port.Begin,
				*other.Port.End,
			)

			differences = append(differences, fmt.Sprintf(
				"port ranges do not match ('%s' vs '%s')",
				scopeRange,
				otherRange,
			))
		} else if other.IsScopePort() {
			differences = append(differences, "port-range scope vs port scope")
		} else {
			differences = append(differences, "port-range scope vs no scope")
		}
	}

	return
}

// FirewallRuleIPAddress represents represents an IP address for firewall configuration.
type FirewallRuleIPAddress struct {
	Address    string `json:"address"`
	PrefixSize *int   `json:"prefixSize,omitempty"`
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
}

// Enable enables the firewall rule.
func (configuration *FirewallRuleConfiguration) Enable() *FirewallRuleConfiguration {
	configuration.Enabled = true

	return configuration
}

// Disable disables the firewall rule.
func (configuration *FirewallRuleConfiguration) Disable() *FirewallRuleConfiguration {
	configuration.Enabled = false

	return configuration
}

// Accept sets the firewall rule action to FirewallRuleActionAccept.
func (configuration *FirewallRuleConfiguration) Accept() *FirewallRuleConfiguration {
	configuration.Action = FirewallRuleActionAccept

	return configuration
}

// Drop sets the firewall rule action to FirewallRuleActionDrop.
func (configuration *FirewallRuleConfiguration) Drop() *FirewallRuleConfiguration {
	configuration.Action = FirewallRuleActionDrop

	return configuration
}

// IPv4 sets the firewall rule's target IP version to IPv4.
func (configuration *FirewallRuleConfiguration) IPv4() *FirewallRuleConfiguration {
	configuration.IPVersion = FirewallRuleIPVersion4

	return configuration
}

// IPv6 sets the firewall rule's target IP version to IPv6.
func (configuration *FirewallRuleConfiguration) IPv6() *FirewallRuleConfiguration {
	configuration.IPVersion = FirewallRuleIPVersion4

	return configuration
}

// IP sets the firewall rule's target protocol to IP.
func (configuration *FirewallRuleConfiguration) IP() *FirewallRuleConfiguration {
	configuration.Protocol = FirewallRuleProtocolIP

	return configuration
}

// TCP sets the firewall rule's target protocol to TCP.
func (configuration *FirewallRuleConfiguration) TCP() *FirewallRuleConfiguration {
	configuration.Protocol = FirewallRuleProtocolTCP

	return configuration
}

// ICMP sets the firewall rule's target protocol to ICMP.
func (configuration *FirewallRuleConfiguration) ICMP() *FirewallRuleConfiguration {
	configuration.Protocol = FirewallRuleProtocolICMP

	return configuration
}

// PlaceFirst modifies the configuration so that the firewall rule will be placed in the first available position.
func (configuration *FirewallRuleConfiguration) PlaceFirst() *FirewallRuleConfiguration {
	configuration.Placement = FirewallRulePlacement{
		Position: "FIRST",
	}

	return configuration
}

// PlaceBefore modifies the configuration so that the firewall rule will be placed before the specified rule.
func (configuration *FirewallRuleConfiguration) PlaceBefore(beforeRuleName string) *FirewallRuleConfiguration {
	configuration.Placement = FirewallRulePlacement{
		Position:           "BEFORE",
		RelativeToRuleName: &beforeRuleName,
	}

	return configuration
}

// PlaceAfter modifies the configuration so that the firewall rule will be placed after the specified rule.
func (configuration *FirewallRuleConfiguration) PlaceAfter(afterRuleName string) *FirewallRuleConfiguration {
	configuration.Placement = FirewallRulePlacement{
		Position:           "AFTER",
		RelativeToRuleName: &afterRuleName,
	}

	return configuration
}

// MatchAnySourceAddress modifies the configuration so that the firewall rule will match source IP address.
func (configuration *FirewallRuleConfiguration) MatchAnySourceAddress() *FirewallRuleConfiguration {
	return configuration.MatchSourceAddress(FirewallRuleMatchAny)
}

// MatchSourceAddress modifies the configuration so that the firewall rule will match a specific source IP address.
func (configuration *FirewallRuleConfiguration) MatchSourceAddress(address string) *FirewallRuleConfiguration {
	sourceScope := &configuration.Source
	sourceScope.IPAddress = &FirewallRuleIPAddress{
		Address: strings.ToUpper(address),
	}
	sourceScope.AddressList = nil

	return configuration
}

// MatchSourceNetwork modifies the configuration so that the firewall rule will match any source IP address on the specified network.
func (configuration *FirewallRuleConfiguration) MatchSourceNetwork(baseAddress string, prefixSize int) *FirewallRuleConfiguration {
	sourceScope := &configuration.Source
	sourceScope.IPAddress = &FirewallRuleIPAddress{
		Address:    baseAddress,
		PrefixSize: &prefixSize,
	}
	sourceScope.AddressList = nil

	return configuration
}

// MatchSourceAddressList modifies the configuration so that the firewall rule will match a specific source IP address list.
func (configuration *FirewallRuleConfiguration) MatchSourceAddressList(addressListID string) *FirewallRuleConfiguration {
	sourceScope := &configuration.Source
	sourceScope.IPAddress = nil
	sourceScope.AddressList = nil
	sourceScope.AddressListID = &addressListID

	return configuration
}

// MatchAnySourcePort modifies the configuration so that the firewall rule will match any source port.
func (configuration *FirewallRuleConfiguration) MatchAnySourcePort() *FirewallRuleConfiguration {
	sourceScope := &configuration.Source
	sourceScope.Port = nil
	sourceScope.PortListID = nil

	return configuration
}

// MatchSourcePort modifies the configuration so that the firewall rule will match a specific source port.
func (configuration *FirewallRuleConfiguration) MatchSourcePort(port int) *FirewallRuleConfiguration {
	sourceScope := &configuration.Source
	sourceScope.Port = &FirewallRulePort{
		Begin: port,
	}
	sourceScope.PortListID = nil

	return configuration
}

// MatchSourcePortRange modifies the configuration so that the firewall rule will match any source port in the specified range.
func (configuration *FirewallRuleConfiguration) MatchSourcePortRange(beginPort int, endPort int) *FirewallRuleConfiguration {
	sourceScope := &configuration.Source
	sourceScope.Port = &FirewallRulePort{
		Begin: beginPort,
		End:   &endPort,
	}
	sourceScope.PortListID = nil

	return configuration
}

// MatchSourcePortList modifies the configuration so that the firewall rule will match any source port appearing on the specified port list (or its children).
func (configuration *FirewallRuleConfiguration) MatchSourcePortList(portListID string) *FirewallRuleConfiguration {
	sourceScope := &configuration.Source
	sourceScope.Port = nil
	sourceScope.PortListID = &portListID

	return configuration
}

// MatchAnyDestinationAddress modifies the configuration so that the firewall rule will match any destination IP address.
func (configuration *FirewallRuleConfiguration) MatchAnyDestinationAddress() *FirewallRuleConfiguration {
	return configuration.MatchDestinationAddress(FirewallRuleMatchAny)
}

// MatchDestinationAddress modifies the configuration so that the firewall rule will match a specific destination IP address.
func (configuration *FirewallRuleConfiguration) MatchDestinationAddress(address string) *FirewallRuleConfiguration {
	destinationScope := &configuration.Destination
	destinationScope.IPAddress = &FirewallRuleIPAddress{
		Address: strings.ToUpper(address),
	}
	destinationScope.AddressList = nil

	return configuration
}

// MatchDestinationNetwork modifies the configuration so that the firewall rule will match any destination IP address on the specified network.
func (configuration *FirewallRuleConfiguration) MatchDestinationNetwork(baseAddress string, prefixSize int) *FirewallRuleConfiguration {
	destinationScope := &configuration.Destination
	destinationScope.IPAddress = &FirewallRuleIPAddress{
		Address:    baseAddress,
		PrefixSize: &prefixSize,
	}
	destinationScope.AddressList = nil

	return configuration
}

// MatchDestinationAddressList modifies the configuration so that the firewall rule will match a specific destination IP address list (and, optionally, port).
func (configuration *FirewallRuleConfiguration) MatchDestinationAddressList(addressListID string) *FirewallRuleConfiguration {
	destinationScope := &configuration.Destination
	destinationScope.AddressList = nil
	destinationScope.AddressListID = &addressListID

	return configuration
}

// MatchAnyDestinationPort modifies the configuration so that the firewall rule will match any destination port.
func (configuration *FirewallRuleConfiguration) MatchAnyDestinationPort() *FirewallRuleConfiguration {
	destinationScope := &configuration.Destination
	destinationScope.Port = nil
	destinationScope.PortListID = nil

	return configuration
}

// MatchDestinationPort modifies the configuration so that the firewall rule will match a specific destination port.
func (configuration *FirewallRuleConfiguration) MatchDestinationPort(port int) *FirewallRuleConfiguration {
	destinationScope := &configuration.Destination
	destinationScope.Port = &FirewallRulePort{
		Begin: port,
	}
	destinationScope.PortListID = nil

	return configuration
}

// MatchDestinationPortRange modifies the configuration so that the firewall rule will match any destination port in the specified range.
func (configuration *FirewallRuleConfiguration) MatchDestinationPortRange(beginPort int, endPort int) *FirewallRuleConfiguration {
	destinationScope := &configuration.Destination
	destinationScope.Port = &FirewallRulePort{
		Begin: beginPort,
		End:   &endPort,
	}
	destinationScope.PortListID = nil

	return configuration
}

// MatchDestinationPortList modifies the configuration so that the firewall rule will match any destination port appearing on the specified port list (or its children).
func (configuration *FirewallRuleConfiguration) MatchDestinationPortList(portListID string) *FirewallRuleConfiguration {
	destinationScope := &configuration.Destination
	destinationScope.Port = nil
	destinationScope.PortListID = &portListID

	return configuration
}

// ToFirewallRule converts the FirewallRuleConfiguration to a FirewallRule (for use in test scenarios).
func (configuration *FirewallRuleConfiguration) ToFirewallRule() FirewallRule {
	return FirewallRule{
		Name:            configuration.Name,
		IPVersion:       configuration.IPVersion,
		Protocol:        configuration.Protocol,
		Source:          configuration.Source,
		Destination:     configuration.Destination,
		Action:          configuration.Action,
		Enabled:         configuration.Enabled,
		NetworkDomainID: configuration.NetworkDomainID,
	}
}

// FirewallRulePlacement describes the placement for a firewall rule.
type FirewallRulePlacement struct {
	Position           string  `json:"position"`
	RelativeToRuleName *string `json:"relativeToRule,omitempty"`
}

type editFirewallRule struct {
	ID      string `json:"id"`
	Enabled bool   `json:"enabled"`
}

type deleteFirewallRule struct {
	ID string `json:"id"`
}

// GetFirewallRule retrieves the Firewall rule with the specified Id.
// Returns nil if no Firewall rule is found with the specified Id.
func (client *Client) GetFirewallRule(id string) (rule *FirewallRule, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return nil, err
	}

	requestURI := fmt.Sprintf("%s/network/firewallRule/%s",
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
func (client *Client) ListFirewallRules(networkDomainID string, paging *Paging) (rules *FirewallRules, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return nil, err
	}

	requestURI := fmt.Sprintf("%s/network/firewallRule?networkDomainId=%s&%s",
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

	requestURI := fmt.Sprintf("%s/network/createFirewallRule",
		url.QueryEscape(organizationID),
	)
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
	firewallRuleIDMessage := apiResponse.GetFieldMessage("firewallRuleId")
	if firewallRuleIDMessage == nil {
		return "", apiResponse.ToError("Received an unexpected response (missing 'firewallRuleId') with status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	return *firewallRuleIDMessage, nil
}

// EditFirewallRule updates the configuration for a firewall rule (enable / disable).
// This operation is synchronous.
func (client *Client) EditFirewallRule(id string, enabled bool) error {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return err
	}

	requestURI := fmt.Sprintf("%s/network/editFirewallRule",
		url.QueryEscape(organizationID),
	)
	request, err := client.newRequestV22(requestURI, http.MethodPost, &editFirewallRule{
		ID:      id,
		Enabled: enabled,
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
		return apiResponse.ToError("Request to edit firewall rule failed with unexpected status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	return nil
}

// DeleteFirewallRule deletes the specified FirewallRule rule.
func (client *Client) DeleteFirewallRule(id string) error {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return err
	}

	requestURI := fmt.Sprintf("%s/network/deleteFirewallRule",
		url.QueryEscape(organizationID),
	)
	request, err := client.newRequestV22(requestURI, http.MethodPost,
		&deleteFirewallRule{id},
	)
	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return err
	}

	apiResponse, err := readAPIResponseAsJSON(responseBody, statusCode)
	if err != nil {
		return err
	}

	if apiResponse.ResponseCode != ResponseCodeOK {
		return apiResponse.ToError("Request to delete firewall rule '%s' failed with unexpected status code %d (%s): %s", id, statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	return nil
}
