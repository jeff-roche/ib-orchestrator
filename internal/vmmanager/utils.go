package vmmanager

import (
	"errors"

	"libvirt.org/go/libvirt"
)

var (
	ErrNotImplemented = errors.New("functionality not yet implemented")
)

func getLibvirtConnection() (*libvirt.Connect, error) {
	return libvirt.NewConnect("qemu:///system") // TODO: make sure this is the right connection string (make configurable?)
}
