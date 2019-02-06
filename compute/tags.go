package compute

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// Tag represents a tag applied to an asset.
type Tag struct {
	Name  string `json:"tagKeyName"`
	Value string `json:"value"`
}

// TagDetail represents detailed information about a tag applied to an asset.
type TagDetail struct {
	AssetType        string `json:"assetType"`
	AssetID          string `json:"assetId"`
	AssetName        string `json:"assetName"`
	DataCenterID     string `json:"datacenterId"`
	TagKeyID         string `json:"tagKeyId"`
	Name             string `json:"tagKeyName"`
	Value            string `json:"value"`
	IsValueRequired  bool   `json:"valueRequired"`
	DisplayOnReports bool   `json:"displayOnReport"`
}

// ToTag converts the TagDetail to a Tag.
func (tagDetail *TagDetail) ToTag() Tag {
	return Tag{
		Name:  tagDetail.Name,
		Value: tagDetail.Value,
	}
}

// TagDetails represents a page of TagDetail results.
type TagDetails struct {
	Items []TagDetail `json:"tag"`

	PagedResult
}

// Request body when applying tags to an asset.
type applyTags struct {
	AssetType string `json:"assetType"`
	AssetID   string `json:"assetId"`
	Tags      []Tag  `json:"tag"`
}

// Request body when removing tags from an asset.
type removeTags struct {
	AssetType string   `json:"assetType"`
	AssetID   string   `json:"assetId"`
	TagNames  []string `json:"tagKeyName"`
}

// TagKey represents a key for asset tags.
type TagKey struct {
	ID string `json:"id"`

	tagKey
}

// TagKeys represents a page of TagKey results.
type TagKeys struct {
	Items []TagKey `json:"tagKey"`

	PagedResult
}

// Common fields for a tag key.
type tagKey struct {
	Name             string `json:"name"`
	Description      string `json:"description"`
	IsValueRequired  bool   `json:"valueRequired"`
	DisplayOnReports bool   `json:"displayOnReport"`
}

// Request body for deleting a tag key.
type deleteTagKey struct {
	ID string `json:"id"`
}

// GetAssetTagsByType gets all tags applied to assets of the specified type.
//
// datacenterID is optional (pass an empty string for tags from all datacenters).
//
// Note that due to a bug in the CloudControl API, when you go past the last page if results, you'll receive an UNEXPECTED_ERROR response code.
func (client *Client) GetAssetTagsByType(assetType string, datacenterID string, paging *Paging) (tags *TagDetails, err error) {
	if paging == nil {
		paging = DefaultPaging()
	}

	organizationID, err := client.getOrganizationID()
	if err != nil {
		return nil, err
	}

	requestURI := fmt.Sprintf("%s/tag/tag?assetType=%s&%s",
		url.QueryEscape(organizationID),
		url.QueryEscape(assetType),
		paging.toQueryParameters(),
	)
	if datacenterID != "" {
		requestURI += fmt.Sprintf("&datacenterId=%s",
			url.QueryEscape(datacenterID),
		)
	}
	request, err := client.newRequestV25(requestURI, http.MethodGet, nil)
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

	tags = &TagDetails{}
	err = json.Unmarshal(responseBody, tags)

	return tags, err
}

// GetAssetTags gets all tags applied to the specified asset.
//
// Note that due to a bug in the CloudControl API, when you go past the last page if results, you'll receive an UNEXPECTED_ERROR response code.
func (client *Client) GetAssetTags(assetID string, assetType string, paging *Paging) (tags *TagDetails, err error) {
	if paging == nil {
		paging = DefaultPaging()
	}

	organizationID, err := client.getOrganizationID()
	if err != nil {
		return nil, err
	}

	requestURI := fmt.Sprintf("%s/tag/tag?assetId=%s&assetType=%s&%s",
		url.QueryEscape(organizationID),
		url.QueryEscape(assetID),
		url.QueryEscape(assetType),
		paging.toQueryParameters(),
	)
	request, err := client.newRequestV25(requestURI, http.MethodGet, nil)
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

	tags = &TagDetails{}
	err = json.Unmarshal(responseBody, tags)

	return tags, err
}

// ApplyAssetTags applies the specified tags to an asset.
func (client *Client) ApplyAssetTags(assetID string, assetType string, tags ...Tag) (response *APIResponseV2, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return nil, err
	}

	requestURI := fmt.Sprintf("%s/tag/applyTags",
		url.QueryEscape(organizationID),
	)
	request, err := client.newRequestV25(requestURI, http.MethodPost, &applyTags{
		AssetID:   assetID,
		AssetType: assetType,
		Tags:      tags,
	})
	if err != nil {
		return nil, err
	}
	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return nil, err
	}

	return readAPIResponseAsJSON(responseBody, statusCode)
}

// RemoveAssetTags removes the specified tags from an asset.
func (client *Client) RemoveAssetTags(assetID string, assetType string, tagNames ...string) (response *APIResponseV2, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return nil, err
	}

	requestURI := fmt.Sprintf("%s/tag/removeTags",
		url.QueryEscape(organizationID),
	)
	request, err := client.newRequestV25(requestURI, http.MethodPost, &removeTags{
		AssetID:   assetID,
		AssetType: assetType,
		TagNames:  tagNames,
	})
	if err != nil {
		return nil, err
	}
	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return nil, err
	}

	return readAPIResponseAsJSON(responseBody, statusCode)
}

// GetTagKey retrieves the tag key with the specified Id.
// Returns nil if no tag key is found with the specified Id.
func (client *Client) GetTagKey(id string) (tagKey *TagKey, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return nil, err
	}

	requestURI := fmt.Sprintf("%s/tags/tagKey/%s",
		url.QueryEscape(organizationID),
		url.QueryEscape(id),
	)
	request, err := client.newRequestV25(requestURI, http.MethodGet, nil)
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

		return nil, apiResponse.ToError("Request to retrieve tag key '%s' failed with status code %d (%s): %s", id, statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	tagKey = &TagKey{}
	err = json.Unmarshal(responseBody, tagKey)
	if err != nil {
		return nil, err
	}

	return tagKey, nil
}

// ListTagKeys lists all tag keys that apply to the specified network domain.
func (client *Client) ListTagKeys(paging *Paging) (tagKeys *TagKeys, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return nil, err
	}

	requestURI := fmt.Sprintf("%s/tag/tagKey?orderBy=name&%s",
		url.QueryEscape(organizationID),
		paging.toQueryParameters(),
	)
	request, err := client.newRequestV25(requestURI, http.MethodGet, nil)
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

		return nil, apiResponse.ToError("Request to list tag keys failed with status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	tagKeys = &TagKeys{}
	err = json.Unmarshal(responseBody, tagKeys)

	return tagKeys, err
}

// CreateTagKey creates a new tag key.
func (client *Client) CreateTagKey(name string, description string, isValueRequired bool, displayOnReports bool) (tagKeyID string, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return "", err
	}

	requestURI := fmt.Sprintf("%s/tag/createTagKey",
		url.QueryEscape(organizationID),
	)
	request, err := client.newRequestV25(requestURI, http.MethodPost, &tagKey{
		Name:             name,
		Description:      description,
		IsValueRequired:  isValueRequired,
		DisplayOnReports: displayOnReports,
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
		return "", apiResponse.ToError("Request to create tag key '%s' failed with unexpected status code %d (%s): %s", name, statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	// Expected: "info" { "name": "tagKeyId", "value": "the-Id-of-the-new-tag-key" }
	if len(apiResponse.FieldMessages) != 1 || apiResponse.FieldMessages[0].FieldName != "tagKeyId" {
		return "", apiResponse.ToError("Received an unexpected response (missing 'tagKeyId') with status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	return apiResponse.FieldMessages[0].Message, nil
}

// DeleteTagKey deletes the specified TagKey rule.
func (client *Client) DeleteTagKey(id string) error {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return err
	}

	requestURI := fmt.Sprintf("%s/tag/deleteTagKey",
		url.QueryEscape(organizationID),
	)
	request, err := client.newRequestV25(requestURI, http.MethodPost,
		&deleteTagKey{id},
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
		return apiResponse.ToError("Request to delete tag key '%s' failed with unexpected status code %d (%s): %s", id, statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	return nil
}
