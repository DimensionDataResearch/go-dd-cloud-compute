package compute

import "encoding/xml"

// Account represents the details for a compute account.
type Account struct {
	// The XML name for the "Account" data contract
	XMLName xml.Name `xml:"Account"`

	// The compute API user name.
	UserName string `xml:"userName"`

	// The user's full name.
	FullName string `xml:"fullName"`

	// The user's first name.
	FirstName string `xml:"firstName"`

	// The user's last name.
	LastName string `xml:"lastName"`

	// The user's email address.
	EmailAddress string `xml:"emailAddress"`

	// The user's department.
	Department string `xml:"department"`

	// The Id of the user's organisation.
	OrganizationID string `xml:"orgId"`
}
