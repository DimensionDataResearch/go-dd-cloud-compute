package compute

import (
	"fmt"
	"log"
	"time"
)

// WaitForDeploy waits for a resource's pending deployment operation to complete.
func (client *Client) WaitForDeploy(resourceType ResourceType, id string, timeout time.Duration) (resource Resource, err error) {
	return client.waitForPendingOperation(resourceType, id, "Deploy", ResourceStatusPendingAdd, false, timeout)
}

// WaitForServerClone waits for a server's pending clone operation to complete.
//
// Pass the customer image Id, not the server Id.
func (client *Client) WaitForServerClone(customerImageID string, timeout time.Duration) (resource Resource, err error) {
	return client.waitForPendingOperation(ResourceTypeCustomerImage, customerImageID, "Clone", ResourceStatusPendingAdd, false, timeout)
}

// WaitForEdit waits for a resource's pending edit operation to complete.
func (client *Client) WaitForEdit(resourceType ResourceType, id string, timeout time.Duration) (resource Resource, err error) {
	return client.WaitForChange(resourceType, id, "Edit", timeout)
}

// WaitForAdd waits for a resource's pending add operation to complete.
func (client *Client) WaitForAdd(resourceType ResourceType, id string, actionDescription string, timeout time.Duration) (resource Resource, err error) {
	return client.waitForPendingOperation(resourceType, id, actionDescription, ResourceStatusPendingAdd, false, timeout)
}

// WaitForChange waits for a resource's pending change operation to complete.
func (client *Client) WaitForChange(resourceType ResourceType, id string, actionDescription string, timeout time.Duration) (resource Resource, err error) {
	return client.waitForPendingOperation(resourceType, id, actionDescription, ResourceStatusPendingChange, false, timeout)
}

// WaitForNestedDeleteChange waits for a resource's pending change operation (actually the delete of a nested resource) to complete.
func (client *Client) WaitForNestedDeleteChange(resourceType ResourceType, id string, actionDescription string, timeout time.Duration) (resource Resource, err error) {
	return client.waitForPendingOperation(resourceType, id, actionDescription, ResourceStatusPendingChange, true, timeout)
}

// WaitForDelete waits for a resource's pending deletion to complete.
func (client *Client) WaitForDelete(resourceType ResourceType, id string, timeout time.Duration) error {
	_, err := client.waitForPendingOperation(resourceType, id, "Delete", ResourceStatusPendingDelete, true, timeout)

	return err
}

// WaitForServerBackupStatus waits for a server to have the specified backup status.
func (client *Client) WaitForServerBackupStatus(serverID string, actionDescription string, targetStatus string, timeout time.Duration) (server *Server, err error) {
	waitTimeout := time.NewTimer(timeout)
	defer waitTimeout.Stop()

	pollTicker := time.NewTicker(5 * time.Second)
	defer pollTicker.Stop()

	for {
		select {
		case <-waitTimeout.C:
			return nil, fmt.Errorf("Timed out after waiting %d seconds for %s of server '%s' to complete",
				timeout/time.Second,
				actionDescription,
				serverID,
			)

		case <-pollTicker.C:
			log.Printf("Polling status for server '%s'...", serverID)
			if client.isCancellationRequested {
				log.Printf("Client indicates that cancellation of pending requests has been requested.")

				return nil, &OperationCancelledError{
					OperationDescription: fmt.Sprintf("Wait for %s of server '%s'",
						actionDescription,
						serverID,
					),
				}
			}

			server, err := client.GetServer(serverID)
			if err != nil {
				return nil, err
			}

			if server == nil {
				return nil, fmt.Errorf("no server was found with Id '%s'", serverID)
			}
			if server.Backup == nil {
				if targetStatus == ResourceStatusDeleted {
					return nil, nil
				}

				return nil, fmt.Errorf("expected backup status for server '%s' to become '%s', but server backup is now disabled", serverID, targetStatus)
			}

			switch server.Backup.State {
			case ResourceStatusNormal:
				log.Printf("%s of server '%s' has successfully completed.", actionDescription, serverID)

				return server, nil

			case ResourceStatusPendingAdd:
				log.Printf("%s of server '%s' is still in progress...", actionDescription, serverID)

				continue
			case ResourceStatusPendingChange:
				log.Printf("%s of server '%s' is still in progress...", actionDescription, serverID)

				continue
			case ResourceStatusPendingDelete:
				log.Printf("%s of server '%s' is still in progress...", actionDescription, serverID)

				continue
			default:
				log.Printf("Unexpected backup status for server '%s' ('%s').", serverID, server.Backup.State)

				return nil, fmt.Errorf("%s failed for server '%s' ('%s'): encountered unexpected state '%s'", actionDescription, serverID, server.Name, server.Backup.State)
			}
		}
	}
}

// waitForPendingOperation waits for a resource's pending operation to complete (i.e. for its status to become ResourceStatusNormal or the resource to disappear if expectedStatus is ResourceStatusPendingDelete).
func (client *Client) waitForPendingOperation(resourceType ResourceType, id string, actionDescription string, expectedStatus string, isDelete bool, timeout time.Duration) (resource Resource, err error) {
	return client.waitForResourceStatus(resourceType, id, actionDescription, expectedStatus, ResourceStatusNormal, isDelete, timeout)
}

// waitForResourceStatus polls a resource for its status (which is expected to initially be expectedStatus) until it becomes expectedStatus.
// getResource is a function that, given the resource Id, will retrieve the resource.
// timeout is the length of time before the wait times out.
func (client *Client) waitForResourceStatus(resourceType ResourceType, id string, actionDescription string, expectedStatus string, targetStatus string, isDelete bool, timeout time.Duration) (resource Resource, err error) {
	waitTimeout := time.NewTimer(timeout)
	defer waitTimeout.Stop()

	pollTicker := time.NewTicker(5 * time.Second)
	defer pollTicker.Stop()

	resourceDescription, err := GetResourceDescription(resourceType)
	if err != nil {
		return nil, err
	}

	for {
		select {
		case <-waitTimeout.C:
			return nil, fmt.Errorf("Timed out after waiting %d seconds for %s of %s '%s' to complete",
				timeout/time.Second,
				actionDescription,
				resourceDescription,
				id,
			)

		case <-pollTicker.C:
			log.Printf("Polling status for %s '%s'...", resourceDescription, id)
			if client.isCancellationRequested {
				log.Printf("Client indicates that cancellation of pending requests has been requested.")

				return nil, &OperationCancelledError{
					OperationDescription: fmt.Sprintf("Wait for %s of %s '%s'",
						actionDescription,
						resourceDescription,
						id,
					),
				}
			}

			resource, err := client.GetResource(id, resourceType)
			if err != nil {
				return nil, err
			}

			if resource == nil || resource.IsDeleted() {
				if isDelete {
					log.Printf("%s '%s' has been successfully deleted.", resourceDescription, id)

					return nil, nil
				}

				return nil, fmt.Errorf("No %s was found with Id '%s'", resourceDescription, id)
			}

			switch resource.GetState() {
			case ResourceStatusNormal:
				log.Printf("%s of %s '%s' has successfully completed.", actionDescription, resourceDescription, id)

				return resource, nil

			case ResourceStatusPendingAdd:
				log.Printf("%s of %s '%s' is still in progress...", actionDescription, resourceDescription, id)

				continue
			case ResourceStatusPendingChange:
				log.Printf("%s of %s '%s' is still in progress...", actionDescription, resourceDescription, id)

				continue
			case ResourceStatusPendingDelete:
				log.Printf("%s of %s '%s' is still in progress...", actionDescription, resourceDescription, id)

				continue
			default:
				log.Printf("Unexpected status for %s '%s' ('%s').", resourceDescription, id, resource.GetState())

				return nil, fmt.Errorf("%s failed for %s '%s' ('%s'): encountered unexpected state '%s'", actionDescription, resourceDescription, id, resource.GetName(), resource.GetState())
			}
		}
	}
}
