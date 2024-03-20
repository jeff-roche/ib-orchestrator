package network

import (
	"fmt"
	vmutils "snoman/internal/vms/utils"
)

var ErrNetworkNotFound = fmt.Errorf("the specified network could not be found")

// Find will use libvirt to search for the network by name or uuid and return the network spec object
func Find(id string) (*VirtualMachineNetworkSpec, error) {
	lvc, err := vmutils.GetLibvirtConnection()
	if err != nil {
		return nil, fmt.Errorf("unable to initialize libvirt connection: %w", err)
	}
	defer lvc.Close()

	// Make sure we have an active libvirt connection
	if alive, err := lvc.IsAlive(); !alive {
		return nil, fmt.Errorf("can not create virtual machine network, libvirt connection is not alive: %w", err)
	}

	libvirtnet := findNetworkByNameOrUUID(id, lvc)
	if libvirtnet == nil {
		return nil, fmt.Errorf("could not find libvirt network by identifier '%s'", id)
	}

	netxml, err := libvirtnet.GetXMLDesc(0)
	if err != nil {
		return nil, fmt.Errorf("unable to get libvirt network xml description: %w", err)
	}

	spec := &VirtualMachineNetworkSpec{}
	spec.UnmarshalXML([]byte(netxml))

	return spec, nil
}
