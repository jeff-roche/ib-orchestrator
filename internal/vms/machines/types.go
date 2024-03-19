package machines

import (
	"snoman/internal/vms/network"
)

type VirtualMachineSpec struct {
	Name    string                             `yaml:"name" validate:"required"`
	Network *network.VirtualMachineNetworkSpec `yaml:"network"`
}

const (
	DEFAULT_VM_NAME string = "default-sno-vm"
)

func GetDefaultVirtualMachineSpec() *VirtualMachineSpec {
	return &VirtualMachineSpec{
		Name:    DEFAULT_VM_NAME,
		Network: network.GetDefaultVirtualMachineNetworkSpec(),
	}
}
