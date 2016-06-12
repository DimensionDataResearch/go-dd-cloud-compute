package compute

import (
	"fmt"
	"log"
	"time"
)

// GetResourceByID represents a function for retrieving a Resource using its ID.
// Returns nil for resource if the resource does not exist.
type GetResourceByID func(id string) (resource Resource, err error)

// UseResource represents a function that works with a Resource.
type UseResource func(resource Resource) error

// WaitForDeploy waits for a resource's pending deletion to complete.
func WaitForDeploy(id string, resourceDescription string, getResource GetResourceByID, onResourceDeployed UseResource, timeoutSeconds time.Duration) error {
	deployedResource, err := waitForPendingOperation(id, resourceDescription, "Deploy", ResourceStatusPendingAdd, getResource, timeoutSeconds)
	if err != nil {
		return err
	}

	if onResourceDeployed != nil {
		return onResourceDeployed(deployedResource)
	}

	return nil
}

// WaitForEdit waits for a resource's pending edit operation to complete.
func WaitForEdit(id string, resourceDescription string, getResource GetResourceByID, onResourceEdited UseResource, timeoutSeconds time.Duration) error {
	return WaitForChange(id, resourceDescription, "Edit", getResource, onResourceEdited, timeoutSeconds)
}

// WaitForChange waits for a resource's pending change operation to complete.
func WaitForChange(id string, resourceDescription string, actionDescription string, getResource GetResourceByID, onResourceChanged UseResource, timeoutSeconds time.Duration) error {
	changedResource, err := waitForPendingOperation(id, resourceDescription, actionDescription, ResourceStatusPendingChange, getResource, timeoutSeconds)
	if err != nil {
		return err
	}

	if onResourceChanged != nil {
		return onResourceChanged(changedResource)
	}

	return nil
}

// WaitForDelete waits for a resource's pending deletion to complete.
func WaitForDelete(id string, resourceDescription string, getResource GetResourceByID, timeoutSeconds time.Duration) error {
	_, err := waitForPendingOperation(id, resourceDescription, "Delete", ResourceStatusPendingDelete, getResource, timeoutSeconds)

	return err
}

// waitForPendingOperation waits for a resource's pending operation to complete (i.e. for its status to become ResourceStatusNormal or the resource to disappear).
func waitForPendingOperation(id string, resourceDescription string, actionDescription string, expectedStatus string, getResource GetResourceByID, timeoutSeconds time.Duration) (resource Resource, err error) {
	return waitForResourceStatus(id, resourceDescription, actionDescription, expectedStatus, ResourceStatusNormal, getResource, timeoutSeconds)
}

// waitForResourceStatus polls a resource for its status (which is expected to initially be expectedStatus) until it becomes expectedStatus.
// getResource is a function that, given the resource Id, will retrieve the resource.
// timeout is the length of time before the wait times out.
func waitForResourceStatus(id string, resourceDescription string, actionDescription string, expectedStatus string, targetStatus string, getResource GetResourceByID, timeout time.Duration) (resource Resource, err error) {
	waitTimeout := time.NewTimer(timeout)
	defer waitTimeout.Stop()

	pollTicker := time.NewTicker(5 * time.Second)
	defer pollTicker.Stop()

	for {
		select {
		case <-waitTimeout.C:
			return nil, fmt.Errorf("Timed out after waiting %d seconds for %s of %s '%s' to complete", timeout/time.Second, actionDescription, resourceDescription, id)

		case <-pollTicker.C:
			log.Printf("Polling status for %s '%s'...", resourceDescription, id)
			resource, err := getResource(id)
			if err != nil {
				return nil, err
			}
			if err != nil {
				return nil, err
			}

			if resource.IsDeleted() {
				if expectedStatus == ResourceStatusPendingDelete {
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

				return nil, fmt.Errorf("%s failed for server '%s' ('%s'): encountered unexpected state '%s'", actionDescription, id, resource.GetName(), resource.GetState())
			}
		}
	}
}
