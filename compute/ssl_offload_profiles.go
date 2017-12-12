package compute

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// SSLOffloadProfile represents an SSL-offload profile.
type SSLOffloadProfile struct {
	ID                   string          `json:"id"`
	Name                 string          `json:"name"`
	Description          string          `json:"description"`
	SSLDomainCertificate EntityReference `json:"sslDomainCertificate"`
	SSLCertificateChain  EntityReference `json:"sslCertificateChain"`
	Ciphers              string          `json:"ciphers"`
	State                string          `json:"state"`
	DatacenterID         string          `json:"datacenterId"`
	NetworkDomainID      string          `json:"networkDomainId"`
}

// GetID returns the offload profile's Id.
func (offloadProfile *SSLOffloadProfile) GetID() string {
	return offloadProfile.ID
}

// GetResourceType returns the offload profile's resource type.
func (offloadProfile *SSLOffloadProfile) GetResourceType() ResourceType {
	return ResourceTypeSSLOffloadProfile
}

// GetName returns the offload profile's name.
func (offloadProfile *SSLOffloadProfile) GetName() string {
	return offloadProfile.Name
}

// GetState returns the offload profile's current state.
func (offloadProfile *SSLOffloadProfile) GetState() string {
	return offloadProfile.State
}

// IsDeleted determines whether the offload profile has been deleted (is nil).
func (offloadProfile *SSLOffloadProfile) IsDeleted() bool {
	return offloadProfile == nil
}

var _ Resource = &SSLOffloadProfile{}

// ToEntityReference creates an EntityReference representing the offload profile.
func (offloadProfile *SSLOffloadProfile) ToEntityReference() EntityReference {
	return EntityReference{
		ID:   offloadProfile.ID,
		Name: offloadProfile.Name,
	}
}

var _ NamedEntity = &SSLOffloadProfile{}

// SSLOffloadProfiles represents a page of SSLOffloadProfile results.
type SSLOffloadProfiles struct {
	Items []SSLOffloadProfile `json:"sslOffloadProfile"`

	PagedResult
}

// Request body when createing an SSL-offload profile.
type createSSLOffloadProfile struct {
	Name                   string `json:"name"`
	Description            string `json:"description"`
	Ciphers                string `json:"ciphers"`
	SSLDomainCertificateID string `json:"sslDomainCertificateId"`
	SSLCertificateChainID  string `json:"sslCertificateChainId,omitempty"`
	NetworkDomainID        string `json:"networkDomainId"`
}

// Request body when deleting an SSL-offload profile.
type deleteSSLOffloadProfile struct {
	// The SSL-offload profile Id.
	ID string `json:"id"`
}

// ListSSLOffloadProfilesInNetworkDomain retrieves a list of all SSL-offload profiles in the specified network domain.
func (client *Client) ListSSLOffloadProfilesInNetworkDomain(networkDomainID string, paging *Paging) (pools *SSLOffloadProfiles, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return nil, err
	}

	requestURI := fmt.Sprintf("%s/networkDomainVip/sslOffloadProfile?networkDomainId=%s&%s",
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

		return nil, apiResponse.ToError("Request to list SSL-offload profiles in network domain '%s' failed with status code %d (%s): %s", networkDomainID, statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	pools = &SSLOffloadProfiles{}
	err = json.Unmarshal(responseBody, pools)
	if err != nil {
		return nil, err
	}

	return pools, nil
}

// GetSSLOffloadProfile retrieves the SSL-offload profile with the specified Id.
//
// Returns nil if no SSL-offload profile is found with the specified Id.
func (client *Client) GetSSLOffloadProfile(id string) (pool *SSLOffloadProfile, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return nil, err
	}

	requestURI := fmt.Sprintf("%s/networkDomainVip/sslOffloadProfile/%s",
		url.QueryEscape(organizationID),
		url.QueryEscape(id),
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

		if apiResponse.ResponseCode == ResponseCodeResourceNotFound {
			return nil, nil // Not an error, but was not found.
		}

		return nil, apiResponse.ToError("Request to retrieve SSL-offload profile with Id '%s' failed with status code %d (%s): %s", id, statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	pool = &SSLOffloadProfile{}
	err = json.Unmarshal(responseBody, pool)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

// CreateSSLOffloadProfile creates an SSL-offload profile in a network domain.
func (client *Client) CreateSSLOffloadProfile(networkDomainID string, name string, description string, ciphers string, sslDomainCertificateID string, sslCertificateChainID string) (offloadProfileID string, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return "", err
	}

	requestURI := fmt.Sprintf("%s/networkDomainVip/createSslOffloadProfile",
		url.QueryEscape(organizationID),
	)
	request, err := client.newRequestV26(requestURI, http.MethodPost, &createSSLOffloadProfile{
		Name:                   name,
		Description:            description,
		Ciphers:                ciphers,
		SSLDomainCertificateID: sslDomainCertificateID,
		SSLCertificateChainID:  sslCertificateChainID,
		NetworkDomainID:        networkDomainID,
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
		return "", apiResponse.ToError("Request to create SSL-offload profile '%s' failed with status code %d (%s): %s", name, statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	// Expected: "info" { "name": "sslOffloadProfileId", "value": "the-Id-of-the-createed-certificate" }
	sslOffloadProfileIDMessage := apiResponse.GetFieldMessage("sslOffloadProfileId")
	if sslOffloadProfileIDMessage == nil {
		return "", apiResponse.ToError("Received an unexpected response (missing 'sslOffloadProfileId') with status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	return *sslOffloadProfileIDMessage, nil
}

// DeleteSSLOffloadProfile deletes an existing SSL-offload profile.
//
// Returns an error if the operation was not successful.
func (client *Client) DeleteSSLOffloadProfile(id string) (err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return err
	}

	requestURI := fmt.Sprintf("%s/networkDomainVip/deleteSslOffloadProfile",
		url.QueryEscape(organizationID),
	)
	request, err := client.newRequestV26(requestURI, http.MethodPost, &deleteSSLOffloadProfile{id})
	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return err
	}

	apiResponse, err := readAPIResponseAsJSON(responseBody, statusCode)
	if err != nil {
		return err
	}

	if apiResponse.ResponseCode != ResponseCodeOK {
		return apiResponse.ToError("Request to delete SSL-offload profile '%s' failed with unexpected status code %d (%s): %s", id, statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	return nil
}
