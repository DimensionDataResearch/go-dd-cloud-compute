package compute

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// Virtual listener types
const (
	// VirtualListenerTypeStandard represents a standard virtual listener.
	VirtualListenerTypeStandard = "STANDARD"

	// VirtualListenerTypePerformanceLayer4 represents a high-performance (layer 4) virtual listener.
	VirtualListenerTypePerformanceLayer4 = "PERFORMANCE_LAYER_4"
)

// Protocols supported by a standard virtual listener
const (
	// Any protocol (as supported by a standard virtual listener).
	VirtualListenerStandardProtocolAny = "ANY"

	// Transmission Control Protocol (as supported by a standard virtual listener).
	VirtualListenerStandardProtocolTCP = "TCP"

	// Uniform Datagram Protocol (as supported by a standard virtual listener).
	VirtualListenerStandardProtocolUDP = "UDP"

	// Hypertext Transfer Protocol (as supported by a standard virtual listener).
	VirtualListenerStandardProtocolHTTP = "HTTP"

	// File Transfer Protocol (as supported by a standard virtual listener).
	VirtualListenerStandardProtocolFTP = "FTP"

	// Simple Mail Transfer Protocol (as supported by a standard virtual listener).
	VirtualListenerStandardProtocolSMTP = "SMTP"
)

// Protocols supported by a high-performance (layer 4) virtual listener
const (
	// Any protocol (as supported by a high-performance layer 4 virtual listener).
	VirtualListenerPerformanceLayer4ProtocolAny = "ANY"

	// Transmission Control Protocol (as supported by a high-performance layer 4 virtual listener).
	VirtualListenerPerformanceLayer4ProtocolTCP = "TCP"

	// Uniform Datagram Protocol (as supported by a high-performance layer 4 virtual listener).
	VirtualListenerPerformanceLayer4ProtocolUDP = "UDP"

	// Hypertext Transfer Protocol (as supported by a high-performance layer 4 virtual listener).
	VirtualListenerPerformanceLayer4ProtocolHTTP = "HTTP"
)

// Options for source-port preservation on a virtual listener.
const (
	// Source port preservation is enabled.
	SourcePortPreservationEnabled = "PRESERVE"

	// Source port preservation is enabled (strict mode).
	SourcePortPreservationEnabledStrict = "PRESERVE_STRICT"

	// Source port preservation is disabled.
	SourcePortPreservationDisabled = "CHANGE"
)

// VirtualListener represents a virtual listener.
type VirtualListener struct {
	ID                         string                    `json:"id"`
	Name                       string                    `json:"name"`
	Description                string                    `json:"description"`
	Type                       string                    `json:"type"`
	Protocol                   string                    `json:"protocol"`
	ListenerIPAddress          string                    `json:"listenerIpAddress"`
	Port                       int                       `json:"port"`
	Enabled                    bool                      `json:"enabled"`
	ConnectionLimit            int                       `json:"connectionLimit"`
	ConnectionRateLimit        int                       `json:"connectionRateLimit"`
	SourcePortPreservation     string                    `json:"sourcePortPreservation"`
	Pool                       VirtualListenerVIPPoolRef `json:"pool"`
	ClientClonePool            VirtualListenerVIPPoolRef `json:"ClientClonePool"`
	PersistenceProfile         EntityReference           `json:"persistenceProfile"`
	FallbackPersistenceProfile EntityReference           `json:"fallbackPersistenceProfile"`
	OptimizationProfiles       []string                  `json:"optimizationProfile"`
	IRules                     []EntityReference         `json:"irule"`
	State                      string                    `json:"state"`
	CreateTime                 string                    `json:"createTime"`
	NetworkDomainID            string                    `json:"networkDomainId"`
	DataCenterID               string                    `json:"datacenterId"`
}

// GetID returns the virtual listener's Id.
func (virtualListener *VirtualListener) GetID() string {
	return virtualListener.ID
}

// GetResourceType returns the virtual listener's resource type.
func (virtualListener *VirtualListener) GetResourceType() ResourceType {
	return ResourceTypeVirtualListener
}

// GetName returns the virtual listener's name.
func (virtualListener *VirtualListener) GetName() string {
	return virtualListener.Name
}

// GetState returns the virtual listener's current state.
func (virtualListener *VirtualListener) GetState() string {
	return virtualListener.State
}

// IsDeleted determines whether the virtual listener has been deleted (is nil).
func (virtualListener *VirtualListener) IsDeleted() bool {
	return virtualListener == nil
}

// ToEntityReference creates an EntityReference representing the CustomerImage.
func (virtualListener *VirtualListener) ToEntityReference() EntityReference {
	return EntityReference{
		ID:   virtualListener.ID,
		Name: virtualListener.Name,
	}
}

var _ Resource = &VirtualListener{}

// VirtualListeners represents a page of VirtualListener results.
type VirtualListeners struct {
	Items []VirtualListener `json:"virtualListener"`

	PagedResult
}

// VirtualListenerVIPPoolRef represents a VirtualListener's reference to a VIP pool
type VirtualListenerVIPPoolRef struct {
	LoadBalanceMethod string            `json:"loadBalanceMethod"`
	ServiceDownAction string            `json:"serviceDownAction"`
	HealthMonitors    []EntityReference `json:"healthMonitor"`

	EntityReference
}

// NewVirtualListenerConfiguration represents the configuration for a new virtual listener.
type NewVirtualListenerConfiguration struct {
	Name                         string   `json:"name"`
	Description                  string   `json:"description"`
	Type                         string   `json:"type"`
	Protocol                     string   `json:"protocol"`
	ListenerIPAddress            *string  `json:"listenerIpAddress,omitempty"`
	Port                         int      `json:"port,omitempty"`
	Enabled                      bool     `json:"enabled"`
	ConnectionLimit              int      `json:"connectionLimit"`
	ConnectionRateLimit          int      `json:"connectionRateLimit"`
	SourcePortPreservation       string   `json:"sourcePortPreservation"`
	PoolID                       *string  `json:"poolId,omitempty"`
	ClientClonePoolID            *string  `json:"clientClonePoolId,omitempty"`
	PersistenceProfileID         *string  `json:"persistenceProfileId,omitempty"`
	FallbackPersistenceProfileID *string  `json:"fallbackPersistenceProfileId,omitempty"`
	IRuleIDs                     []string `json:"iruleId"`
	OptimizationProfiles         []string `json:"optimizationProfile"`
	NetworkDomainID              string   `json:"networkDomainId"`
}

// EditVirtualListenerConfiguration represents the configuration for editing a virtual listener.
type EditVirtualListenerConfiguration struct {
	ID                     string    `json:"id"`
	Description            *string   `json:"description,omitempty"`
	Enabled                *bool     `json:"enabled,omitempty"`
	ConnectionLimit        *int      `json:"connectionLimit,omitempty"`
	ConnectionRateLimit    *int      `json:"connectionRateLimit,omitempty"`
	SourcePortPreservation *string   `json:"sourcePortPreservation,omitempty"`
	PoolID                 *string   `json:"poolId,omitempty"`
	PersistenceProfileID   *string   `json:"persistenceProfileId,omitempty"`
	IRuleIDs               *[]string `json:"iruleId,omitempty"`
	OptimizationProfiles   *[]string `json:"optimizationProfile,omitempty"`
}

// Request body for deleting a virtual listener.
type deleteVirtualListener struct {
	// The virtual listener ID.
	ID string `json:"id"`
}

// ListVirtualListenersInNetworkDomain retrieves a list of all virtual listeners in the specified network domain.
func (client *Client) ListVirtualListenersInNetworkDomain(networkDomainID string, paging *Paging) (listeners *VirtualListeners, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return nil, err
	}

	requestURI := fmt.Sprintf("%s/networkDomainVip/virtualListener?networkDomainId=%s&%s",
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

		return nil, apiResponse.ToError("Request to list virtual listeners in network domain '%s' failed with status code %d (%s): %s", networkDomainID, statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	listeners = &VirtualListeners{}
	err = json.Unmarshal(responseBody, listeners)
	if err != nil {
		return nil, err
	}

	return listeners, nil
}

// GetVirtualListener retrieves the virtual listener with the specified Id.
// Returns nil if no virtual listener is found with the specified Id.
func (client *Client) GetVirtualListener(id string) (listener *VirtualListener, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return nil, err
	}

	requestURI := fmt.Sprintf("%s/networkDomainVip/virtualListener/%s",
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

		return nil, apiResponse.ToError("Request to retrieve virtual listener with Id '%s' failed with status code %d (%s): %s", id, statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	listener = &VirtualListener{}
	err = json.Unmarshal(responseBody, listener)
	if err != nil {
		return nil, err
	}

	return listener, nil
}

// CreateVirtualListener creates a new VIP node.
// Returns the Id of the new node.
func (client *Client) CreateVirtualListener(listenerConfiguration NewVirtualListenerConfiguration) (nodeID string, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return "", err
	}

	requestURI := fmt.Sprintf("%s/networkDomainVip/createVirtualListener",
		url.QueryEscape(organizationID),
	)
	request, err := client.newRequestV22(requestURI, http.MethodPost, &listenerConfiguration)
	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return "", err
	}

	apiResponse, err := readAPIResponseAsJSON(responseBody, statusCode)
	if err != nil {
		return "", err
	}

	if apiResponse.ResponseCode != ResponseCodeOK {
		return "", apiResponse.ToError("Request to create virtual listener '%s' failed with status code %d (%s): %s", listenerConfiguration.Name, statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	// Expected: "info" { "name": "virtualListenerId", "value": "the-Id-of-the-new-virtual-listener" ... }
	virtualListenerIDMessage := apiResponse.GetFieldMessage("virtualListenerId")
	if virtualListenerIDMessage == nil {
		return "", apiResponse.ToError("Received an unexpected response (missing 'virtualListenerId') with status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	return *virtualListenerIDMessage, nil
}

// EditVirtualListener updates an existing virtual listener.
func (client *Client) EditVirtualListener(id string, listenerConfiguration EditVirtualListenerConfiguration) error {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return err
	}

	editListenerConfiguration := &listenerConfiguration
	editListenerConfiguration.ID = id

	requestURI := fmt.Sprintf("%s/networkDomainVip/editVirtualListener",
		url.QueryEscape(organizationID),
	)
	request, err := client.newRequestV22(requestURI, http.MethodPost, editListenerConfiguration)
	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return err
	}

	apiResponse, err := readAPIResponseAsJSON(responseBody, statusCode)
	if err != nil {
		return err
	}

	if statusCode != http.StatusOK {
		return apiResponse.ToError("Request to edit virtual listener '%s' failed with status code %d (%s): %s", listenerConfiguration.ID, statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	return nil
}

// DeleteVirtualListener deletes an existing virtual listener.
// Returns an error if the operation was not successful.
func (client *Client) DeleteVirtualListener(id string) (err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return err
	}

	requestURI := fmt.Sprintf("%s/networkDomainVip/deleteVirtualListener",
		url.QueryEscape(organizationID),
	)
	request, err := client.newRequestV22(requestURI, http.MethodPost, &deleteVirtualListener{id})
	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return err
	}

	apiResponse, err := readAPIResponseAsJSON(responseBody, statusCode)
	if err != nil {
		return err
	}

	if apiResponse.ResponseCode != ResponseCodeOK {
		return apiResponse.ToError("Request to delete virtual listener '%s' failed with unexpected status code %d (%s): %s", id, statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	return nil
}
