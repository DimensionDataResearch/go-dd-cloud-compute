package compute

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// SSLDomainCertificate represents an SSL certificate applicable to a domain name.
type SSLDomainCertificate struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	State           string `json:"state"`
	DatacenterID    string `json:"datacenterId"`
	NetworkDomainID string `json:"networkDomainId"`
}

// GetID returns the domain certificate's Id.
func (domainCertificate *SSLDomainCertificate) GetID() string {
	return domainCertificate.ID
}

// GetResourceType returns the domain certificate's resource type.
func (domainCertificate *SSLDomainCertificate) GetResourceType() ResourceType {
	return ResourceTypeSSLDomainCertificate
}

// GetName returns the domain certificate's name.
func (domainCertificate *SSLDomainCertificate) GetName() string {
	return domainCertificate.Name
}

// GetState returns the domain certificate's current state.
func (domainCertificate *SSLDomainCertificate) GetState() string {
	return domainCertificate.State
}

// IsDeleted determines whether the domain certificate has been deleted (is nil).
func (domainCertificate *SSLDomainCertificate) IsDeleted() bool {
	return domainCertificate == nil
}

var _ Resource = &SSLDomainCertificate{}

// ToEntityReference creates an EntityReference representing the domain certificate.
func (domainCertificate *SSLDomainCertificate) ToEntityReference() EntityReference {
	return EntityReference{
		ID:   domainCertificate.ID,
		Name: domainCertificate.Name,
	}
}

var _ NamedEntity = &SSLDomainCertificate{}

// SSLDomainCertificates represents a page of SSLDomainCertificate results.
type SSLDomainCertificates struct {
	Items []SSLDomainCertificate `json:"sslDomainCertificate"`

	PagedResult
}

// Request body when importing an SSL certificate for a domain name.
type importSSLDomainCertificate struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	Certificate     string `json:"certificate"`
	Key             string `json:"key"`
	NetworkDomainID string `json:"networkDomainId"`
}

// Request body when deleting an SSL certificate for a domain name.
type deleteSSLDomainCertificate struct {
	// The SSL certificate Id.
	ID string `json:"id"`
}

// ListSSLDomainCertificatesInNetworkDomain retrieves a list of all SSL domain certificates in the specified network domain.
func (client *Client) ListSSLDomainCertificatesInNetworkDomain(networkDomainID string, paging *Paging) (pools *SSLDomainCertificates, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return nil, err
	}

	requestURI := fmt.Sprintf("%s/networkDomainVip/sslDomainCertificate?networkDomainId=%s&%s",
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

		return nil, apiResponse.ToError("Request to list SSL domain certificates in network domain '%s' failed with status code %d (%s): %s", networkDomainID, statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	pools = &SSLDomainCertificates{}
	err = json.Unmarshal(responseBody, pools)
	if err != nil {
		return nil, err
	}

	return pools, nil
}

// GetSSLDomainCertificate retrieves the SSL domain certificate with the specified Id.
// Returns nil if no SSL domain certificate is found with the specified Id.
func (client *Client) GetSSLDomainCertificate(id string) (pool *SSLDomainCertificate, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return nil, err
	}

	requestURI := fmt.Sprintf("%s/networkDomainVip/sslDomainCertificate/%s",
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

		return nil, apiResponse.ToError("Request to retrieve SSL domain certificate with Id '%s' failed with status code %d (%s): %s", id, statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	pool = &SSLDomainCertificate{}
	err = json.Unmarshal(responseBody, pool)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

// ImportSSLDomainCertificate imports an SSL domain certificate into a network domain.
func (client *Client) ImportSSLDomainCertificate(networkDomainID string, name string, description string, certificate string, key string) (certificateID string, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return "", err
	}

	requestURI := fmt.Sprintf("%s/networkDomainVip/importSslDomainCertificate",
		url.QueryEscape(organizationID),
	)
	request, err := client.newRequestV26(requestURI, http.MethodPost, &importSSLDomainCertificate{
		Name:            name,
		Description:     description,
		Certificate:     certificate,
		Key:             key,
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
		return "", apiResponse.ToError("Request to import SSL domain certificate '%s' failed with status code %d (%s): %s", name, statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	// Expected: "info" { "name": "sslDomainCertificateId", "value": "the-Id-of-the-imported-certificate" }
	sslDomainCertificateIDMessage := apiResponse.GetFieldMessage("sslDomainCertificateId")
	if sslDomainCertificateIDMessage == nil {
		return "", apiResponse.ToError("Received an unexpected response (missing 'sslDomainCertificateId') with status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	return *sslDomainCertificateIDMessage, nil
}

// DeleteSSLDomainCertificate deletes an existing SSL domain certificate.
//
// Returns an error if the operation was not successful.
func (client *Client) DeleteSSLDomainCertificate(id string) (err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return err
	}

	requestURI := fmt.Sprintf("%s/networkDomainVip/deleteSslDomainCertificate",
		url.QueryEscape(organizationID),
	)
	request, err := client.newRequestV26(requestURI, http.MethodPost, &deleteSSLDomainCertificate{id})
	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return err
	}

	apiResponse, err := readAPIResponseAsJSON(responseBody, statusCode)
	if err != nil {
		return err
	}

	if apiResponse.ResponseCode != ResponseCodeOK {
		return apiResponse.ToError("Request to delete SSL domain certificate '%s' failed with unexpected status code %d (%s): %s", id, statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	return nil
}
