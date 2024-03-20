package machines

import (
	"fmt"

	"snoman/internal/vms/network"
	vmutils "snoman/internal/vms/utils"
)

func CreateVirtualMachine(spec *VirtualMachineSpec) error {
	err := vmutils.SpecValidator.Struct(spec)
	if err != nil {
		return fmt.Errorf("unable to validate VirtualMachineSpec: %w", err)
	}

	lvc, err := vmutils.GetLibvirtConnection()
	if err != nil {
		return fmt.Errorf("unable to initialize libvirt connection: %w", err)
	}
	defer lvc.Close()

	// Create the network for the VM
	err = network.Create(spec.Network)
	if err != nil {
		return fmt.Errorf("unable to create vm network: %w", err)
	}

	return fmt.Errorf("CreateVirtualMachine %w", vmutils.ErrNotImplemented)
}
