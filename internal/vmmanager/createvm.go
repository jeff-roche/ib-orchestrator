package vmmanager

import (
	"fmt"

	"libvirt.org/go/libvirt"
)

func CreateVirtualMachine() error {
	lvc, err := getLibvirtConnection()
	if err != nil {
		return fmt.Errorf("unable to initialize libvirt connection: %w", err)
	}
	defer lvc.Close()

	// TODO this is example placeholder code, change this to create a VM
	doms, err := lvc.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE)
	if err != nil {
		return fmt.Errorf("unable to get libvirt domains: %w", err)
	}

	fmt.Printf("%d running domains:\n", len(doms))
	for _, dom := range doms {
		name, err := dom.GetName()
		if err == nil {
			fmt.Printf("  %s\n", name)
		}
		dom.Free()
	}
	return fmt.Errorf("CreateVirtualMachine %w", ErrNotImplemented)
}
