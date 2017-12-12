package compute

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// SSLCertificateChain represents an SSL certificate applicable to a domain name.
type SSLCertificateChain struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	State           string `json:"state"`
	DatacenterID    string `json:"datacenterId"`
	NetworkDomainID string `json:"networkDomainId"`
}

// GetID returns the certificate chain's Id.
func (certificateChain *SSLCertificateChain) GetID() string {
	return certificateChain.ID
}

// GetResourceType returns the certificate chain's resource type.
func (certificateChain *SSLCertificateChain) GetResourceType() ResourceType {
	return ResourceTypeSSLCertificateChain
}

// GetName returns the certificate chain's name.
func (certificateChain *SSLCertificateChain) GetName() string {
	return certificateChain.Name
}

// GetState returns the certificate chain's current state.
func (certificateChain *SSLCertificateChain) GetState() string {
	return certificateChain.State
}

// IsDeleted determines whether the certificate chain has been deleted (is nil).
func (certificateChain *SSLCertificateChain) IsDeleted() bool {
	return certificateChain == nil
}

var _ Resource = &SSLCertificateChain{}

// ToEntityReference creates an EntityReference representing the certificate chain.
func (certificateChain *SSLCertificateChain) ToEntityReference() EntityReference {
	return EntityReference{
		ID:   certificateChain.ID,
		Name: certificateChain.Name,
	}
}

var _ NamedEntity = &SSLCertificateChain{}

// SSLCertificateChains represents a page of SSLCertificateChain results.
type SSLCertificateChains struct {
	Items []SSLCertificateChain `json:"sslCertificateChain"`

	PagedResult
}

// Request body when importing an SSL certificate chain.
type importSSLCertificateChain struct {
	Name             string `json:"name"`
	Description      string `json:"description"`
	CertificateChain string `json:"certificateChain"`
	NetworkDomainID  string `json:"networkDomainId"`
}

// Request body when deleting an SSL certificate chain.
type deleteSSLCertificateChain struct {
	// The SSL certificate chain Id.
	ID string `json:"id"`
}

// ListSSLCertificateChainsInNetworkDomain retrieves a list of all SSL certificate chains in the specified network domain.
func (client *Client) ListSSLCertificateChainsInNetworkDomain(networkDomainID string, paging *Paging) (pools *SSLCertificateChains, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return nil, err
	}

	requestURI := fmt.Sprintf("%s/networkDomainVip/sslCertificateChain?networkDomainId=%s&%s",
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

		return nil, apiResponse.ToError("Request to list SSL certificate chains in network domain '%s' failed with status code %d (%s): %s", networkDomainID, statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	pools = &SSLCertificateChains{}
	err = json.Unmarshal(responseBody, pools)
	if err != nil {
		return nil, err
	}

	return pools, nil
}

// GetSSLCertificateChain retrieves the SSL certificate chain with the specified Id.
//
// Returns nil if no SSL certificate chain is found with the specified Id.
func (client *Client) GetSSLCertificateChain(id string) (pool *SSLCertificateChain, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return nil, err
	}

	requestURI := fmt.Sprintf("%s/networkDomainVip/sslCertificateChain/%s",
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

		return nil, apiResponse.ToError("Request to retrieve SSL certificate chain with Id '%s' failed with status code %d (%s): %s", id, statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	pool = &SSLCertificateChain{}
	err = json.Unmarshal(responseBody, pool)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

// ImportSSLCertificateChain imports an SSL certificate chain into a network domain.
func (client *Client) ImportSSLCertificateChain(networkDomainID string, name string, description string, certificateChain string) (certificateChainID string, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return "", err
	}

	requestURI := fmt.Sprintf("%s/networkDomainVip/importSslCertificateChain",
		url.QueryEscape(organizationID),
	)
	request, err := client.newRequestV26(requestURI, http.MethodPost, &importSSLCertificateChain{
		Name:             name,
		Description:      description,
		CertificateChain: certificateChain,
		NetworkDomainID:  networkDomainID,
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
		return "", apiResponse.ToError("Request to import SSL certificate chain '%s' failed with status code %d (%s): %s", name, statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	// Expected: "info" { "name": "sslCertificateChainId", "value": "the-Id-of-the-imported-certificate" }
	sslCertificateChainIDMessage := apiResponse.GetFieldMessage("sslCertificateChainId")
	if sslCertificateChainIDMessage == nil {
		return "", apiResponse.ToError("Received an unexpected response (missing 'sslCertificateChainId') with status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	return *sslCertificateChainIDMessage, nil
}

// DeleteSSLCertificateChain deletes an existing SSL certificate chain.
//
// Returns an error if the operation was not successful.
func (client *Client) DeleteSSLCertificateChain(id string) (err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return err
	}

	requestURI := fmt.Sprintf("%s/networkDomainVip/deleteSslCertificateChain",
		url.QueryEscape(organizationID),
	)
	request, err := client.newRequestV26(requestURI, http.MethodPost, &deleteSSLCertificateChain{id})
	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return err
	}

	apiResponse, err := readAPIResponseAsJSON(responseBody, statusCode)
	if err != nil {
		return err
	}

	if apiResponse.ResponseCode != ResponseCodeOK {
		return apiResponse.ToError("Request to delete SSL certificate chain '%s' failed with unexpected status code %d (%s): %s", id, statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	return nil
}
