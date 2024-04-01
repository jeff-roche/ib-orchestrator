package machines

import (
	"fmt"
	"os/exec"

	"snoman/internal/vms/network"
	vmutils "snoman/internal/vms/utils"
)

func CreateVirtualMachine(spec *VirtualMachineSpec) error {
	err := vmutils.SpecValidator.Struct(spec)
	if err != nil {
		return fmt.Errorf("unable to validate VirtualMachineSpec: %w", err)
	}

	if spec.Network != nil {
		err = network.Create(spec.Network)
		if err != nil {
			return fmt.Errorf("unable to create vm network: %w", err)
		}
	}

	if err := startVirtualMachine(spec); err != nil {
		return fmt.Errorf("unable to start virtual machine: %w", err)
	}

	return nil
}

func startVirtualMachine(spec *VirtualMachineSpec) error {
	args := []string{
		"--connect", "qemu:///system",
		"-n", spec.Name,
		"-r", fmt.Sprint(spec.RAM),
		"--vcpus", fmt.Sprint(spec.CPU),
		"--os-variant", "rhel8.1",
		"--import",
		fmt.Sprintf("--network=network:\"%s\",mac=\"%s\"", spec.Network.Name, spec.Network.MacAddress),
		"--graphics=none",
		"--events", "on_reboot=restart",
		"--cdrom", spec.BipSpec.IsoPath,
		"--disk", fmt.Sprintf("pool=\"%s\",size=\"%d\"", spec.Disk.Pool, spec.Disk.Size),
		"--check", "disk_size=off",
		"--boot", "hd,cdrom",
		"--noautoconsole",
		"--wait=-1",
	}

	cmd := exec.Command("virt-install", args...)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error executing virt-install: %w", err)
	}

	return nil
}
