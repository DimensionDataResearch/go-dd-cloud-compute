// Package compute contains the Go client for Dimension Data's cloud compute API.
package compute

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Client is the client for Dimension Data's cloud compute API.
type Client struct {
	baseAddress string
	username    string
	password    string
	client      *http.Client
}

// NewClient creates a new cloud compute API client.
// region is the cloud compute region identifier.
func NewClient(region string, username string, password string) *Client {
	baseAddress := fmt.Sprintf("https://api-%s.dimensiondata.com", region)

	return &Client{
		baseAddress,
		username,
		password,
		&http.Client{},
	}
}

// GetMyAccount retrieves the current user's account information
func (client *Client) GetMyAccount() (*Account, error) {
	request, err := baseRequestV1(client.baseAddress, "GET", client.username, client.password)
	if err != nil {
		return nil, err
	}
	defer request.Body.Close()

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return nil, err
	}

	account := &Account{}
	err = xml.Unmarshal(body, account)
	if err != nil {
		return nil, err
	}

	return account, nil
}
