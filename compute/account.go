package compute

import "encoding/xml"

// Account represents the details for a compute account.
type Account struct {
	// The XML name for the "Account" data contract
	XMLName xml.Name `xml:"ns3:Account"`

	// The compute API user name.
	UserName string `xml:"ns3:userName"`

	// The user's full name.
	FullName string `xml:"ns3:fullName"`

	// The user's first name.
	FirstName string `xml:"ns3:firstName"`

	// The user's last name.
	LastName string `xml:"ns3:lastName"`

	// The user's email address.
	EmailAddress string `xml:"ns3:emailAddress"`

	// The user's department.
	Department string `xml:"ns3:department"`

	// The Id of the user's organisation.
	OrganizationID string `xml:"ns3:orgId"`
}
