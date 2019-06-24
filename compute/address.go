package compute

import (
	"log"
)

func (client *Client) CheckAddressExists(addresslistId string,
	begin string, network string) (exists bool) {

	// Get existing address list
	addressList, _ := client.GetIPAddressList(addresslistId)

	addresses := addressList.Addresses

	for _, addr := range addresses {
		log.Printf("Address begin: %s, network: %s, addr.Begin: %s", begin, network, addr.Begin)
		if begin == addr.Begin || network == addr.Begin {
			return true
		}
	}

	return false
}

func (client *Client) GetAddressOk(addresslistId string,
	begin string, network string) (address *IPAddressListEntry, exists bool) {

	// Get existing address list
	addressList, _ := client.GetIPAddressList(addresslistId)

	addresses := addressList.Addresses

	for _, addr := range addresses {
		if begin == addr.Begin || network == addr.Begin {
			return &addr, true
		}
	}

	return nil, false
}

func (client *Client) AddAddress(addresslistId string,
	begin string, end string, network string, prefixSize int) (address *IPAddressListEntry, err error) {

	// Get existing address list
	addressList, err := client.GetIPAddressList(addresslistId)

	var newAddress IPAddressListEntry

	if begin != "" {
		newAddress.Begin = begin
	}

	if end != "" {
		newAddress.End = &end
	}

	// Comply to existing IPAddressListEntry model to use begin for both IP Begin and Network
	if network != "" {
		newAddress.Begin = network
	}

	if prefixSize > 0 {
		newAddress.PrefixSize = &prefixSize
	}

	// Append address to current list
	editRequest := addressList.BuildEditRequest()
	editRequest.Addresses = append(editRequest.Addresses, newAddress)

	err = client.EditIPAddressList(editRequest)
	if err != nil {
		return nil, err
	}

	log.Printf("Appended address %+v\n to addresslist '%s'.", newAddress, addresslistId)

	return &newAddress, nil
}

// Note: ipAddress can represent begin or network
func (client *Client) DeleteAddress(addresslistId string,
	ipAddress string) (addressList *IPAddressList, err error) {

	// Get existing address list
	addressList, err = client.GetIPAddressList(addresslistId)

	addresses := addressList.Addresses

	for i, addr := range addresses {

		if ipAddress == addr.Begin {
			// Copy last element to index i.
			addresses[i] = addresses[len(addresses)-1]
			// Erase last element (write zero value)
			addresses[len(addresses)-1] = IPAddressListEntry{}
			// Truncate slice
			addresses = addresses[:len(addresses)-1]
		}
	}

	// Append address to current list
	editRequest := addressList.BuildEditRequest()
	editRequest.Addresses = addresses

	err = client.EditIPAddressList(editRequest)
	if err != nil {
		return nil, err
	}

	log.Printf("Deleted address %+v\n from addresslist '%s'.", ipAddress, addresslistId)

	return nil, nil
}
