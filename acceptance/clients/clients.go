// Package clients contains functions for creating OpenStack service clients
// for use in acceptance tests. It also manages the required environment
// variables to run the tests.
package clients

import (
	"fmt"
	"os"
	"strings"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
)

// AcceptanceTestChoices contains image and flavor selections for use by the acceptance tests.
type AcceptanceTestChoices struct {
	// ImageID contains the ID of a valid image.
	ImageID string

	// FlavorID contains the ID of a valid flavor.
	FlavorID string

	// FlavorIDResize contains the ID of a different flavor available on the same OpenStack installation, that is distinct
	// from FlavorID.
	FlavorIDResize string

	// FloatingIPPool contains the name of the pool from where to obtain floating IPs.
	FloatingIPPoolName string

	// NetworkName is the name of a network to launch the instance on.
	NetworkName string
}

// AcceptanceTestChoicesFromEnv populates a ComputeChoices struct from environment variables.
// If any required state is missing, an `error` will be returned that enumerates the missing properties.
func AcceptanceTestChoicesFromEnv() (*AcceptanceTestChoices, error) {
	imageID := os.Getenv("OS_IMAGE_ID")
	flavorID := os.Getenv("OS_FLAVOR_ID")
	flavorIDResize := os.Getenv("OS_FLAVOR_ID_RESIZE")
	networkName := os.Getenv("OS_NETWORK_NAME")
	floatingIPPoolName := os.Getenv("OS_POOL_NAME")

	missing := make([]string, 0, 3)
	if imageID == "" {
		missing = append(missing, "OS_IMAGE_ID")
	}
	if flavorID == "" {
		missing = append(missing, "OS_FLAVOR_ID")
	}
	if flavorIDResize == "" {
		missing = append(missing, "OS_FLAVOR_ID_RESIZE")
	}
	if floatingIPPoolName == "" {
		missing = append(missing, "OS_POOL_NAME")
	}
	if networkName == "" {
		networkName = "private"
	}

	notDistinct := ""
	if flavorID == flavorIDResize {
		notDistinct = "OS_FLAVOR_ID and OS_FLAVOR_ID_RESIZE must be distinct."
	}

	if len(missing) > 0 || notDistinct != "" {
		text := "You're missing some important setup:\n"
		if len(missing) > 0 {
			text += " * These environment variables must be provided: " + strings.Join(missing, ", ") + "\n"
		}
		if notDistinct != "" {
			text += " * " + notDistinct + "\n"
		}

		return nil, fmt.Errorf(text)
	}

	return &AcceptanceTestChoices{ImageID: imageID, FlavorID: flavorID, FlavorIDResize: flavorIDResize, FloatingIPPoolName: floatingIPPoolName, NetworkName: networkName}, nil
}

// NewBlockStorageV1Client returns a *ServiceClient for making calls
// to the OpenStack Block Storage v1 API. An error will be returned
// if authentication or client creation was not possible.
func NewBlockStorageV1Client() (*gophercloud.ServiceClient, error) {
	ao, err := openstack.AuthOptionsFromEnv()
	if err != nil {
		return nil, err
	}

	client, err := openstack.AuthenticatedClient(ao)
	if err != nil {
		return nil, err
	}

	return openstack.NewBlockStorageV1(client, gophercloud.EndpointOpts{
		Region: os.Getenv("OS_REGION_NAME"),
	})
}

// NewComputeV2Client returns a *ServiceClient for making calls
// to the OpenStack Compute v2 API. An error will be returned
// if authentication or client creation was not possible.
func NewComputeV2Client() (*gophercloud.ServiceClient, error) {
	ao, err := openstack.AuthOptionsFromEnv()
	if err != nil {
		return nil, err
	}

	client, err := openstack.AuthenticatedClient(ao)
	if err != nil {
		return nil, err
	}

	return openstack.NewComputeV2(client, gophercloud.EndpointOpts{
		Region: os.Getenv("OS_REGION_NAME"),
	})
}

// NewIdentityV2Client returns a *ServiceClient for making calls
// to the OpenStack Identity v2 API. An error will be returned
// if authentication or client creation was not possible.
func NewIdentityV2Client() (*gophercloud.ServiceClient, error) {
	ao, err := openstack.AuthOptionsFromEnv()
	if err != nil {
		return nil, err
	}

	client, err := openstack.AuthenticatedClient(ao)
	if err != nil {
		return nil, err
	}

	return openstack.NewIdentityV2(client, gophercloud.EndpointOpts{
		Region: os.Getenv("OS_REGION_NAME"),
	})
}
