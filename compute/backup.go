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

// EnableBackup represents the request body when enabling Cloud Backup for a server.
type EnableBackup struct {
	// The XML name for the "EnableBackup" structure
	XMLName xml.Name `xml:"http://oec.api.opsource.net/schemas/backup NewBackup"`

	// The Cloud Backup service plan ("Essentials" or "Advanced") to use.
	ServicePlan string `xml:"servicePlan,attr"`
}

// ChangeBackupServicePlan represents the request body when changing the Cloud Backup service plan for a server.
type ChangeBackupServicePlan struct {
	// The XML name for the "ChangeBackupServicePlan" structure
	XMLName xml.Name `xml:"http://oec.api.opsource.net/schemas/backup ModifyBackup"`

	// The new service plan ("Essentials" or "Advanced") to use.
	ServicePlan string `xml:"servicePlan,attr"`
}

// BackupClientTypes represents the types of backup client enabled for a server.
type BackupClientTypes struct {
	// The XML name for the BackupClientTypes structure
	XMLName xml.Name `xml:"http://oec.api.opsource.net/schemas/backup BackupClientTypes"`

	// Types of backup client enabled for the server.
	Items []BackupClientType `xml:"http://oec.api.opsource.net/schemas/backup backupClientType"`
}

// BackupClientType represents a types of backup client enabled for a server.
type BackupClientType struct {
	// The XML name for the BackupClientType structure
	XMLName xml.Name `xml:"http://oec.api.opsource.net/schemas/backup BackupClientType"`

	Type         string `xml:"type,attr"`
	IsFileSystem bool   `xml:"isFileSystem,attr"`
	Description  string `xml:"description,attr"`
}

// BackupStoragePolicies represents a list of backup storage policies.
type BackupStoragePolicies struct {
	// The XML name for the BackupClientType structure
	XMLName xml.Name `xml:"http://oec.api.opsource.net/schemas/backup BackupStoragePolicies"`

	// The storage policies.
	Items []BackupStoragePolicy `xml:"http://oec.api.opsource.net/schemas/backup storagePolicy"`
}

// BackupStoragePolicy represents a Cloud Backup storage policy.
type BackupStoragePolicy struct {
	// The XML name for the BackupStoragePolicy structure
	XMLName xml.Name `xml:"http://oec.api.opsource.net/schemas/backup StoragePolicy"`

	// The policy name.
	Name string `xml:"name,attr"`

	// The policy's backup retention period (in days).
	RetentionPeriodInDays int `xml:"retentionPeriodInDays,attr"`

	// The secondary location where backups are stored.
	SecondaryLocation string `xml:"secondaryLocation,attr"`
}

// BackupSchedulePolicies represents a list of backup schedule policies.
type BackupSchedulePolicies struct {
	// The XML name for the BackupClientType structure
	XMLName xml.Name `xml:"http://oec.api.opsource.net/schemas/backup BackupSchedulePolicies"`

	// The schedule policies.
	Items []BackupSchedulePolicy `xml:"http://oec.api.opsource.net/schemas/backup schedulePolicy"`
}

// BackupSchedulePolicy represents a Cloud Backup schedule policy.
type BackupSchedulePolicy struct {
	// The XML name for the BackupSchedulePolicy structure
	XMLName xml.Name `xml:"http://oec.api.opsource.net/schemas/backup SchedulePolicy"`

	// The policy name.
	Name string `xml:"name,attr"`

	// The policy description.
	Description string `xml:"description,attr"`
}

// NewBackupClient represents the request body when adding a backup client to a server.
type NewBackupClient struct {
	// The XML name for the NewBackupClient structure
	XMLName xml.Name `xml:"http://oec.api.opsource.net/schemas/backup NewBackupClient"`

	// The client type (e.g. "FA.Linux").
	Type string `xml:"http://oec.api.opsource.net/schemas/backup type"`

	// The name of the storage policy to use.
	StoragePolicyName string `xml:"http://oec.api.opsource.net/schemas/backup storagePolicyName"`

	// The name of the schedule policy to use.
	SchedulePolicyName string `xml:"http://oec.api.opsource.net/schemas/backup schedulePolicyName"`

	// The client alerting configuration (if any).
	Alerting *BackupClientAlerting `xml:"http://oec.api.opsource.net/schemas/backup alerting,omitempty"`
}

// ModifyBackupClient represents the request body when modifying a server's backup client.
type ModifyBackupClient struct {
	// The XML name for the ModifyBackupClient structure
	XMLName xml.Name `xml:"http://oec.api.opsource.net/schemas/backup ModifyBackupClient"`

	// The client type (e.g. "FA.Linux").
	Type string `xml:"http://oec.api.opsource.net/schemas/backup type"`

	// The name of the storage policy to use.
	StoragePolicyName string `xml:"http://oec.api.opsource.net/schemas/backup storagePolicyName"`

	// The name of the schedule policy to use.
	SchedulePolicyName string `xml:"http://oec.api.opsource.net/schemas/backup schedulePolicyName"`

	// The client alerting configuration (if any).
	Alerting *BackupClientAlerting `xml:"http://oec.api.opsource.net/schemas/backup alerting,omitempty"`
}

// BackupClientAlerting represents the alerting configuration for a backup client.
type BackupClientAlerting struct {
	// The XML name for the BackupClientAlerting structure
	XMLName xml.Name `xml:"http://oec.api.opsource.net/schemas/backup BackupClientAlerting"`

	// When should the alert be triggered?
	//
	// Must be one of "ON_FAILURE", "ON_SUCCESS", "ON_SUCCESS_OR_FAILURE".
	Trigger string `xml:"trigger,attr"`

	// Email addresses for alert notifications.
	EmailAddresses []string `xml:"http://oec.api.opsource.net/schemas/backup emailAddress"`
}

// ServerBackupDetails represents detailed backup information for a server.
type ServerBackupDetails struct {
	// The XML name for the BackupClientDetail structure
	XMLName xml.Name `xml:"http://oec.api.opsource.net/schemas/backup BackupDetails"`

	// The server's associated asset Id.
	AssetID string `xml:"assetId,attr"`

	// The server's associated backup service plan.
	ServicePlan string `xml:"servicePlan,attr"`

	// The server state.
	State string `xml:"state,attr"`

	// Detailed information about the server's backup clients.
	Clients []BackupClientDetail `xml:"http://oec.api.opsource.net/schemas/backup backupClient"`
}

// BackupClientDetail represents the detail for a specific backup client on a server.
type BackupClientDetail struct {
	// The XML name for the BackupClientDetail structure
	XMLName xml.Name `xml:"http://oec.api.opsource.net/schemas/backup BackupClientDetail"`

	// The client Id.
	ID string `xml:"id,attr"`

	// The client type (e.g. "FA.Linux").
	Type string `xml:"type,attr"`

	// Does the backup client operate at the file-system level?
	IsFileSystem bool `xml:"isFileSystem,attr"`

	// A description of the backup client.
	Description string `xml:"description,attr"`

	// The name of the storage policy to use.
	StoragePolicyName string `xml:"http://oec.api.opsource.net/schemas/backup storagePolicyName"`

	// The name of the schedule policy to use.
	SchedulePolicyName string `xml:"http://oec.api.opsource.net/schemas/backup schedulePolicyName"`

	// The client alerting configuration (if any).
	Alerting *BackupClientAlerting `xml:"http://oec.api.opsource.net/schemas/backup alerting,omitempty"`

	// The client download URL.
	DownloadURL string `xml:"http://oec.api.opsource.net/schemas/backup downloadUrl"`
}
