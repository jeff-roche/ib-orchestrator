package machines

import (
	"fmt"
	"snoman/internal/vms/network"
	vmutils "snoman/internal/vms/utils"

	"gopkg.in/yaml.v2"
)

type VirtualMachineSpec struct {
	Name    string                             `yaml:"name" validate:"required"`
	Network *network.VirtualMachineNetworkSpec `yaml:"network,omitempty" validate:"omitempty"`
}

const (
	DEFAULT_VM_NAME string = "default-sno-vm"
)

func GetDefaultVirtualMachineSpec() *VirtualMachineSpec {
	spec := &VirtualMachineSpec{
		Name:    DEFAULT_VM_NAME,
		Network: network.GetDefaultVirtualMachineNetworkSpec(),
	}

	return spec
}

func (spec VirtualMachineSpec) Validate() error {
	err := vmutils.SpecValidator.Struct(spec)
	if err != nil {
		return fmt.Errorf("unable to validate VirtualMachineSpec: %w", err)
	}

	return nil
}

func (spec VirtualMachineSpec) MarshalYAML() (string, error) {
	if err := spec.Validate(); err != nil {
		return "", err
	}

	data, err := yaml.Marshal(spec)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (spec *VirtualMachineSpec) UnmarshalYAML(yamlData []byte) error {
	err := yaml.Unmarshal(yamlData, spec)
	if err != nil {
		return fmt.Errorf("unable to parse the spec: %w", err)
	}

	if err := spec.Validate(); err != nil {
		return err
	}

	return nil
}
