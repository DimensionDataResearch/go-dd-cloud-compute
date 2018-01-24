package compute

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

const (
	// BackupServicePlanEssentials represents the basic service plan for Cloud Backup
	BackupServicePlanEssentials = "Essentials"

	// BackupServicePlanAdvanced represents the advanced service plan for Cloud Backup
	BackupServicePlanAdvanced = "Advanced"

	// BackupServicePlanEnterprise represents the enterprise service plan for Cloud Backup
	BackupServicePlanEnterprise = "Enterprise"

	// BackupClientStatusOffline indicates that a backup client is not currently contactable by Cloud Backup.
	BackupClientStatusOffline = "Offline"

	// BackupClientStatusUnregistered indicates that a backup client has never registered with Cloud Backup.
	BackupClientStatusUnregistered = "Unregistered"

	// BackupClientStatusUnconfigured indicates that a backup client has not been configured by Cloud Backup.
	BackupClientStatusUnconfigured = "Unconfigured"

	// BackupClientStatusActive indicates that a backup client is currently contactable by Cloud Backup and is ready to service requests.
	BackupClientStatusActive = "Active"
)

// ServerBackup represents the backup configuration for a server.
type ServerBackup struct {
	ServicePlan string `json:"servicePlan"`
	State       string `json:"state"`
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
	XMLName xml.Name // Always a child element, so we'll accept element name from containing element's declaration

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
	XMLName xml.Name // Always a child element, so we'll accept element name from containing element's declaration

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
	XMLName xml.Name // Always a child element, so we'll accept element name from containing element's declaration

	// The policy name.
	Name string `xml:"name,attr"`

	// The policy description.
	Description string `xml:"description,attr"`
}

// BackupClientAlerting represents the alerting configuration for a backup client.
type BackupClientAlerting struct {
	// The XML name for the BackupClientAlerting structure
	XMLName xml.Name // Always a child element, so we'll accept element name from containing element's declaration

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

// GetClientByID retrieves the BackupClientDetail (if any) with the specified Id.
func (backupDetails *ServerBackupDetails) GetClientByID(id string) *BackupClientDetail {
	for index := range backupDetails.Clients {
		backupClient := &backupDetails.Clients[index]
		if backupClient.ID == id {
			return backupClient
		}
	}

	return nil
}

// BackupClientDetail represents the detail for a specific backup client on a server.
type BackupClientDetail struct {
	// The XML name for the BackupClientDetail structure
	XMLName xml.Name // Always a child element, so we'll accept element name from containing element's declaration

	// The client Id.
	ID string `xml:"id,attr"`

	// The client type (e.g. "FA.Linux").
	Type string `xml:"type,attr"`

	// Does the backup client operate at the file-system level?
	IsFileSystem bool `xml:"isFileSystem,attr"`

	// The backup client's status (e.g. "Unregistered", etc).
	Status string `xml:"status,attr"`

	// A description of the backup client.
	Description string `xml:"http://oec.api.opsource.net/schemas/backup description"`

	// The name of the storage policy to use.
	StoragePolicyName string `xml:"http://oec.api.opsource.net/schemas/backup storagePolicyName"`

	// The name of the schedule policy to use.
	SchedulePolicyName string `xml:"http://oec.api.opsource.net/schemas/backup schedulePolicyName"`

	// The client alerting configuration (if any).
	Alerting *BackupClientAlerting `xml:"http://oec.api.opsource.net/schemas/backup alerting,omitempty"`

	// The server's total backup size (in GB).
	TotalBackupSizeGb int `xml:"http://oec.api.opsource.net/schemas/backup totalBackupSizeGb"`

	// The client download URL.
	DownloadURL string `xml:"http://oec.api.opsource.net/schemas/backup downloadUrl"`
}

// newBackup represents the request body when enabling Cloud Backup for a server.
type newBackup struct {
	// The XML name for the "newBackup" structure
	XMLName xml.Name `xml:"http://oec.api.opsource.net/schemas/backup NewBackup"`

	// The Cloud Backup service plan ("Essentials" or "Advanced") to use.
	ServicePlan string `xml:"servicePlan,attr"`
}

// modifyBackup represents the request body when changing the Cloud Backup service plan for a server.
type modifyBackup struct {
	// The XML name for the "modifyBackup" structure
	XMLName xml.Name `xml:"http://oec.api.opsource.net/schemas/backup ModifyBackup"`

	// The new service plan ("Essentials" or "Advanced") to use.
	ServicePlan string `xml:"servicePlan,attr"`
}

// newBackupClient represents the request body when adding a backup client to a server.
type newBackupClient struct {
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

// modifyBackupClient represents the request body when modifying a server's backup client.
type modifyBackupClient struct {
	// The XML name for the ModifyBackupClient structure
	XMLName xml.Name `xml:"http://oec.api.opsource.net/schemas/backup ModifyBackupClient"`

	// The name of the storage policy to use.
	StoragePolicyName string `xml:"http://oec.api.opsource.net/schemas/backup storagePolicyName"`

	// The name of the schedule policy to use.
	SchedulePolicyName string `xml:"http://oec.api.opsource.net/schemas/backup schedulePolicyName"`

	// The client alerting configuration (if any).
	Alerting *BackupClientAlerting `xml:"http://oec.api.opsource.net/schemas/backup alerting,omitempty"`
}

// GetServerBackupDetails retrieves detailed information about a server's Cloud Backup status
func (client *Client) GetServerBackupDetails(serverID string) (*ServerBackupDetails, error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return nil, err
	}

	requestURI := fmt.Sprintf("%s/server/%s/backup",
		url.QueryEscape(organizationID),
		url.QueryEscape(serverID),
	)
	request, err := client.newRequestV1(requestURI, http.MethodGet, nil)

	if err != nil {
		return nil, errors.Wrapf(err, "failed to create request for retrieving backup details of server '%s'", serverID)
	}

	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to execute request for retrieving backup details of server '%s'", serverID)
	}

	if statusCode != http.StatusOK {
		response := &APIResponseV1{}
		err = xml.Unmarshal(responseBody, response)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse error response for retrieving backup details of server '%s'", serverID)
		}

		if response.ResultCode == ResultCodeBackupNotEnabledForServer || response.ResultCode == ResultCodeBackupEnablementInProgressForServer {
			return nil, nil
		}

		return nil, response.ToError("failed to retrieve backup details of server '%s' (HTTP %d / %s): %s",
			serverID,
			statusCode,
			response.ResultCode,
			response.Message,
		)
	}

	serverBackupDetails := &ServerBackupDetails{}
	err = xml.Unmarshal(responseBody, serverBackupDetails)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse response for retrieval backup details of server '%s'", serverID)
	}

	return serverBackupDetails, nil
}

// EnableServerBackup enables Cloud Backup for a server
func (client *Client) EnableServerBackup(serverID string, servicePlan string) error {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return err
	}

	requestURI := fmt.Sprintf("%s/server/%s/backup",
		url.QueryEscape(organizationID),
		url.QueryEscape(serverID),
	)
	request, err := client.newRequestV1(requestURI, http.MethodPost, &newBackup{
		ServicePlan: servicePlan,
	})
	if err != nil {
		return errors.Wrapf(err, "failed to create request for enabling backup on server '%s'", serverID)
	}

	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return errors.Wrapf(err, "failed to execute request for enabling backup on server '%s'", serverID)
	}

	response := &APIResponseV1{}
	err = xml.Unmarshal(responseBody, response)
	if err != nil {
		return errors.Wrapf(err, "failed to parse response for enabling backup on server '%s'", serverID)
	}

	if response.Result != ResultSuccess {
		return response.ToError("failed to enable backup for server '%s' (HTTP %d / %s): %s",
			serverID,
			statusCode,
			response.ResultCode,
			response.Message,
		)
	}

	return nil
}

// DisableServerBackup disables Cloud Backup for a server
func (client *Client) DisableServerBackup(serverID string) error {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return err
	}

	requestURI := fmt.Sprintf("%s/server/%s/backup?disable",
		url.QueryEscape(organizationID),
		url.QueryEscape(serverID),
	)
	request, err := client.newRequestV1(requestURI, http.MethodGet, nil)
	if err != nil {
		return errors.Wrapf(err, "failed to create request for disabling backup on server '%s'", serverID)
	}

	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return errors.Wrapf(err, "failed to execute request for disabling backup on server '%s'", serverID)
	}

	response := &APIResponseV1{}
	err = xml.Unmarshal(responseBody, response)
	if err != nil {
		return errors.Wrapf(err, "failed to parse response for disabling backup on server '%s'", serverID)
	}

	if response.Result != ResultSuccess {
		return response.ToError("failed to disable backup for server '%s' (HTTP %d / %s): %s",
			serverID,
			statusCode,
			response.ResultCode,
			response.Message,
		)
	}

	return nil
}

// ChangeServerBackupServicePlan changes a server's Cloud Backup service plan
func (client *Client) ChangeServerBackupServicePlan(serverID string, servicePlan string) error {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return err
	}

	requestURI := fmt.Sprintf("%s/server/%s/backup/modify",
		url.QueryEscape(organizationID),
		url.QueryEscape(serverID),
	)
	request, err := client.newRequestV1(requestURI, http.MethodPost, &modifyBackup{
		ServicePlan: servicePlan,
	})
	if err != nil {
		return errors.Wrapf(err, "failed to create request for changing backup service plan on server '%s'", serverID)
	}

	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return errors.Wrapf(err, "failed to execute request for changing backup service plan on server '%s'", serverID)
	}

	response := &APIResponseV1{}
	err = xml.Unmarshal(responseBody, response)
	if err != nil {
		return errors.Wrapf(err, "failed to parse response for changing backup service plan on server '%s'", serverID)
	}

	if response.Result != ResultSuccess {
		return response.ToError("failed to change backup service plan for server '%s' (HTTP %d / %s): %s",
			serverID,
			statusCode,
			response.ResultCode,
			response.Message,
		)
	}

	return nil
}

// GetServerBackupClientTypes retrieves a list of a server's configured Cloud Backup clients.
func (client *Client) GetServerBackupClientTypes(serverID string) (*BackupClientTypes, error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return nil, err
	}

	requestURI := fmt.Sprintf("%s/server/%s/backup/client/type",
		url.QueryEscape(organizationID),
		url.QueryEscape(serverID),
	)
	request, err := client.newRequestV1(requestURI, http.MethodGet, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create request for retrieving backup client types on server '%s'", serverID)
	}

	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to execute request for retrieving backup client types on server '%s'", serverID)
	}

	if statusCode != http.StatusOK {
		response := &APIResponseV1{}
		err = xml.Unmarshal(responseBody, response)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse response for retrieving backup client types on server '%s'", serverID)
		}

		return nil, response.ToError("failed to retrieve backup client types for server '%s' (HTTP %d / %s): %s",
			serverID,
			statusCode,
			response.ResultCode,
			response.Message,
		)
	}

	clientTypes := &BackupClientTypes{}
	err = xml.Unmarshal(responseBody, clientTypes)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse response for retrieving backup client types on server '%s'", serverID)
	}

	return clientTypes, nil
}

// GetServerBackupStoragePolicies retrieves a list of a server's configured Cloud Backup storage policies.
func (client *Client) GetServerBackupStoragePolicies(serverID string) (*BackupStoragePolicies, error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return nil, err
	}

	requestURI := fmt.Sprintf("%s/server/%s/backup/client/storagePolicy",
		url.QueryEscape(organizationID),
		url.QueryEscape(serverID),
	)
	request, err := client.newRequestV1(requestURI, http.MethodGet, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create request for retrieving backup storage policies on server '%s'", serverID)
	}

	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to execute request for retrieving backup storage policies on server '%s'", serverID)
	}

	if statusCode != http.StatusOK {
		response := &APIResponseV1{}
		err = xml.Unmarshal(responseBody, response)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse response for retrieving backup storage policies on server '%s'", serverID)
		}

		return nil, response.ToError("failed to retrieve backup storage policies for server '%s' (HTTP %d / %s): %s",
			serverID,
			statusCode,
			response.ResultCode,
			response.Message,
		)
	}

	clientTypes := &BackupStoragePolicies{}
	err = xml.Unmarshal(responseBody, clientTypes)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse response for retrieving backup storage policies on server '%s'", serverID)
	}

	return clientTypes, nil
}

// GetServerBackupSchedulePolicies retrieves a list of a server's configured Cloud Backup schedule policies.
func (client *Client) GetServerBackupSchedulePolicies(serverID string) (*BackupSchedulePolicies, error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return nil, err
	}

	requestURI := fmt.Sprintf("%s/server/%s/backup/client/schedulePolicy",
		url.QueryEscape(organizationID),
		url.QueryEscape(serverID),
	)
	request, err := client.newRequestV1(requestURI, http.MethodGet, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create request for retrieving backup schedule policies on server '%s'", serverID)
	}

	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to execute request for retrieving backup schedule policies on server '%s'", serverID)
	}

	if statusCode != http.StatusOK {
		response := &APIResponseV1{}
		err = xml.Unmarshal(responseBody, response)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse response for retrieving backup schedule policies on server '%s'", serverID)
		}

		return nil, response.ToError("failed to retrieve backup schedule policies for server '%s' (HTTP %d / %s): %s",
			serverID,
			statusCode,
			response.ResultCode,
			response.Message,
		)
	}

	clientTypes := &BackupSchedulePolicies{}
	err = xml.Unmarshal(responseBody, clientTypes)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse response for retrieving backup schedule policies on server '%s'", serverID)
	}

	return clientTypes, nil
}

// AddServerBackupClient adds a backup client to a server.
func (client *Client) AddServerBackupClient(serverID string, clientType string, schedulePolicyName string, storagePolicyName string, alerting *BackupClientAlerting) (clientID string, clientDownloadURL string, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return "", "", err
	}

	requestURI := fmt.Sprintf("%s/server/%s/backup/client",
		url.QueryEscape(organizationID),
		url.QueryEscape(serverID),
	)
	request, err := client.newRequestV1(requestURI, http.MethodPost, &newBackupClient{
		Type:               clientType,
		SchedulePolicyName: schedulePolicyName,
		StoragePolicyName:  storagePolicyName,
		Alerting:           alerting,
	})
	if err != nil {
		return "", "", errors.Wrapf(err, "failed to create request for adding '%s' backup client to server '%s'", clientType, serverID)
	}

	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return "", "", errors.Wrapf(err, "failed to execute request for adding '%s' backup client to server '%s'", clientType, serverID)
	}

	response := &APIResponseV1{}
	err = xml.Unmarshal(responseBody, response)
	if err != nil {
		return "", "", errors.Wrapf(err, "failed to parse response for adding '%s' backup client to server '%s'", clientType, serverID)
	}

	if response.Result != ResultSuccess {
		return "", "", response.ToError("failed to add '%s' backup client to server '%s' (HTTP %d / %s): %s",
			clientType,
			serverID,
			statusCode,
			response.ResultCode,
			response.Message,
		)
	}

	backupClientID := response.GetAdditionalInformation("backupClient.id")
	if backupClientID == nil {
		return "", "", response.ToError("request to add '%s' backup client to server '%s' succeeded, but the CloudControl API did not return a valid Id for the new backup client",
			clientType,
			serverID,
		)
	}

	backupClientDownloadURL := response.GetAdditionalInformation("backupClient.downloadUrl")
	if backupClientDownloadURL == nil {
		return *backupClientID, "", nil // No specific download URL for this client.
	}

	return *backupClientID, *backupClientDownloadURL, nil
}

// RemoveServerBackupClient removes a backup client from a server.
func (client *Client) RemoveServerBackupClient(serverID string, clientID string) error {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return err
	}

	requestURI := fmt.Sprintf("%s/server/%s/backup/client/%s?remove",
		url.QueryEscape(organizationID),
		url.QueryEscape(serverID),
		url.QueryEscape(clientID),
	)
	request, err := client.newRequestV1(requestURI, http.MethodGet, nil)
	if err != nil {
		return errors.Wrapf(err, "failed to create request for removing backup client '%s' from server '%s'", clientID, serverID)
	}

	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return errors.Wrapf(err, "failed to execute request for removing backup client '%s' from server '%s'", clientID, serverID)
	}

	response := &APIResponseV1{}
	err = xml.Unmarshal(responseBody, response)
	if err != nil {
		return errors.Wrapf(err, "failed to parse response for removing backup client '%s' from server '%s'", clientID, serverID)
	}

	if response.Result != ResultSuccess {
		return response.ToError("failed to remove backup client '%s' from server '%s' (HTTP %d / %s): %s",
			clientID,
			serverID,
			statusCode,
			response.ResultCode,
			response.Message,
		)
	}

	return nil
}

// ModifyServerBackupClient modifies one of a server's existing backup clients.
func (client *Client) ModifyServerBackupClient(serverID string, clientID string, schedulePolicyName string, storagePolicyName string, alerting *BackupClientAlerting) (clientDownloadURL string, err error) {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return "", err
	}

	requestURI := fmt.Sprintf("%s/server/%s/backup/client/%s",
		url.QueryEscape(organizationID),
		url.QueryEscape(serverID),
		url.QueryEscape(clientID),
	)
	request, err := client.newRequestV1(requestURI, http.MethodPost, &modifyBackupClient{
		SchedulePolicyName: schedulePolicyName,
		StoragePolicyName:  storagePolicyName,
		Alerting:           alerting,
	})
	if err != nil {
		return "", errors.Wrapf(err, "failed to create request for modifying backup client '%s' in server '%s'", clientID, serverID)
	}

	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return "", errors.Wrapf(err, "failed to execute request for modifying backup client '%s' in server '%s'", clientID, serverID)
	}

	response := &APIResponseV1{}
	err = xml.Unmarshal(responseBody, response)
	if err != nil {
		return "", errors.Wrapf(err, "failed to parse response for modifying backup client '%s' in server '%s'", clientID, serverID)
	}

	if response.Result != ResultSuccess {
		return "", response.ToError("failed to modify backup client '%s' in server '%s' (HTTP %d / %s): %s",
			clientID,
			serverID,
			statusCode,
			response.ResultCode,
			response.Message,
		)
	}

	backupClientDownloadURL := response.GetAdditionalInformation("backupClient.downloadUrl")
	if backupClientDownloadURL == nil {
		return "", nil // No specific download URL for this client.
	}

	return *backupClientDownloadURL, nil
}

// CancelBackupClientJobs cancels all running jobs (if any) for a backup client.
func (client *Client) CancelBackupClientJobs(serverID string, clientID string) error {
	organizationID, err := client.getOrganizationID()
	if err != nil {
		return err
	}

	requestURI := fmt.Sprintf("%s/server/%s/backup/client/%s?cancelJob",
		url.QueryEscape(organizationID),
		url.QueryEscape(serverID),
		url.QueryEscape(clientID),
	)
	request, err := client.newRequestV1(requestURI, http.MethodGet, nil)
	if err != nil {
		return errors.Wrapf(err, "failed to create request for canceling all jobs for backup client '%s' on server '%s'", clientID, serverID)
	}

	responseBody, statusCode, err := client.executeRequest(request)
	if err != nil {
		return errors.Wrapf(err, "failed to execute request for canceling all jobs for backup client '%s' on server '%s'", clientID, serverID)
	}

	response := &APIResponseV1{}
	err = xml.Unmarshal(responseBody, response)
	if err != nil {
		return errors.Wrapf(err, "failed to parse response for canceling all jobs for backup client '%s' on server '%s'", clientID, serverID)
	}

	if response.Result != ResultSuccess {
		return response.ToError("failed to cancel jobs for backup client '%s' on server '%s' (HTTP %d / %s): %s",
			clientID,
			serverID,
			statusCode,
			response.ResultCode,
			response.Message,
		)
	}

	return nil
}
