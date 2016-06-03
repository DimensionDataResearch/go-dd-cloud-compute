package compute

import "encoding/xml"

// Account represents the details for a compute account.
type Account struct {
	XMLName        xml.Name `xml:"ns3:Account"`
	UserName       string   `xml:"ns3:userName"`
	FullName       string   `xml:"ns3:fullName"`
	FirstName      string   `xml:"ns3:firstName"`
	LastName       string   `xml:"ns3:lastName"`
	EmailAddress   string   `xml:"ns3:emailAddress"`
	Department     string   `xml:"ns3:department"`
	OrganizationID string   `xml:"ns3:orgId"`
}
