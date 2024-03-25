package bip

import (
	"snoman/internal/biputils"
	"snoman/internal/vms/machines"
)

type BootstrapInPlaceSpec struct {
	MachineConfig *machines.VirtualMachineSpec
	PullSecret    string
	PublicKey     string
	Workdir       string
	IsoSpec       *biputils.BootstrapInPlaceIsoSpec
}
