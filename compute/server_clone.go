package compute

import (
	"fmt"
	"net/http"
	"net/url"
)

type cloneServer struct {
	ServerID             string `json:"id"`
	ImageName            string `json:"imageName"`
	ImageDescription     string `json:"description,omitempty"`
	GuestOsCustomization bool   `json:"guestOsCustomization"`
}

// CloneServer clones a server to create a customer image.
func (client *Client) CloneServer(serverID string, imageName string, imageDescription string, preventGuestOSCustomisation bool) (imageID string, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return "", err
	}

	requestURI := fmt.Sprintf("%s/server/cloneServer",
		url.QueryEscape(organizationID),
	)
	request, err := client.newRequestV24(requestURI, http.MethodPost, &cloneServer{
		ServerID:             serverID,
		ImageName:            imageName,
		ImageDescription:     imageDescription,
		GuestOsCustomization: !preventGuestOSCustomisation,
	})
	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return "", err
	}

	apiResponse, err := readAPIResponseAsJSON(responseBody, statusCode)
	if err != nil {
		return "", err
	}

	if apiResponse.ResponseCode != ResponseCodeInProgress {
		return "", apiResponse.ToError("Request to clone '%s' failed with status code %d (%s): %s", serverID, statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	// Expected: "info" { "name": "imageId", "value": "the-Id-of-the-new-image" }
	serverIDMessage := apiResponse.GetFieldMessage("imageId")
	if serverIDMessage == nil {
		return "", apiResponse.ToError("Received an unexpected response (missing 'imageId') with status code %d (%s): %s", statusCode, apiResponse.ResponseCode, apiResponse.Message)
	}

	return *serverIDMessage, nil
}
