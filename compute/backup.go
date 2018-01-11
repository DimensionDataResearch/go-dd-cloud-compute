package compute

import (
	"encoding/xml"
)

const (
	// BackupServicePlanEssentials represents the basic service plan for Cloud Backup
	BackupServicePlanEssentials = "Essentials"

	// BackupServicePlanAdvanced represents the advanced service plan for Cloud Backup
	BackupServicePlanAdvanced = "Advanced"
)

// EnableBackup represents the request body when enabling Cloud Backup for a Server.
type EnableBackup struct {
	// The XML name for the "EnableBackup" data contract
	XMLName xml.Name `xml:"http://oec.api.opsource.net/schemas/backup NewBackup"`

	// The Cloud Backup service plan ("Essentials" or "Advanced") to use.
	ServicePlan string `xml:"servicePlan,attr"`
}

// ChangeBackupServicePlan represents the request body when changing the Cloud Backup service plan for a Server.
type ChangeBackupServicePlan struct {
	// The XML name for the "ChangeBackupServicePlan" data contract
	XMLName xml.Name `xml:"http://oec.api.opsource.net/schemas/backup ModifyBackup"`

	// The new service plan ("Essentials" or "Advanced") to use.
	ServicePlan string `xml:"servicePlan,attr"`
}

// BackupClientTypes represents the types of backup client enabled for a Server.
type BackupClientTypes struct {
	// The XML name for the "BackupClientTypes" data contract
	XMLName xml.Name `xml:"http://oec.api.opsource.net/schemas/backup BackupClientTypes"`

	// types of backup client enabled for the Server.
	Types []BackupClientType `xml:"http://oec.api.opsource.net/schemas/backup backupClientType"`
}

// BackupClientType represents a types of backup client enabled for a Server.
type BackupClientType struct {
	// The XML name for the "BackupClientType" data contract
	XMLName xml.Name `xml:"http://oec.api.opsource.net/schemas/backup BackupClientType"`

	Type         string `xml:"type,attr"`
	IsFileSystem bool   `xml:"isFileSystem,attr"`
	Description  string `xml:"description,attr"`
}
