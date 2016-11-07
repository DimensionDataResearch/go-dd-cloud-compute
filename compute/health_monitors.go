package compute

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// HealthMonitor represents a load-balancer persistence (stickiness) profile.
type HealthMonitor struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	IsNodeCompatible bool   `json:"nodeCompatible"`
	IsPoolCompatible bool   `json:"poolCompatible"`
}

// HealthMonitors represents a page of HealthMonitor results.
type HealthMonitors struct {
	Items []HealthMonitor `json:"defaultHealthMonitor"`

	PagedResult
}

// ListDefaultHealthMonitors retrieves a list of all default load-balancing health monitors in the specified network domain.
func (client *Client) ListDefaultHealthMonitors(networkDomainID string, paging *Paging) (healthMonitors *HealthMonitors, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return nil, err
	}

	requestURI := fmt.Sprintf("%s/networkDomainVip/defaultHealthMonitor?networkDomainId=%s&%s",
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

		return nil, apiResponse.ToError("Request to list default health monitors in network domain '%s' failed with status code %d (%s): %s", networkDomainID, statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	healthMonitors = &HealthMonitors{}
	err = json.Unmarshal(responseBody, healthMonitors)
	if err != nil {
		return nil, err
	}

	return healthMonitors, nil
}
