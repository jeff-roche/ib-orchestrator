package network

import (
	"fmt"
	"snoman/internal/logger"
	vmutils "snoman/internal/vms/utils"
)

func Create(spec *VirtualMachineNetworkSpec) error {
	log := logger.Get()

	lvc, err := vmutils.GetLibvirtConnection()
	if err != nil {
		return fmt.Errorf("unable to initialize libvirt connection: %w", err)
	}
	defer lvc.Close()

	// Make sure we have an active libvirt connection
	if alive, err := lvc.IsAlive(); !alive {
		return fmt.Errorf("can not create virtual machine network, libvirt connection is not alive: %w", err)
	}

	// Make sure this network does not already exist
	dupNameNet := findNetworkByNameOrUUID(spec.Name, lvc)
	dupUuidNet := findNetworkByNameOrUUID(spec.UUID, lvc)
	if dupNameNet != nil || dupUuidNet != nil {
		var netxml string
		var err error

		if dupNameNet != nil {
			netxml, _ = dupNameNet.GetXMLDesc(0)
			err = fmt.Errorf("a network with name '%s' already exists", spec.Name)
			dupNameNet.Free()
		} else {
			netxml, _ = dupUuidNet.GetXMLDesc(0)
			err = fmt.Errorf("a network with UUID '%s' already exists", spec.UUID)
			dupUuidNet.Free()
		}

		log.Debugln("Duplicate network found")
		log.Debugf("Duplicate network XML:\n%s", netxml)
		return err
	}

	netxml, err := spec.MarshalXML()
	if err != nil {
		return fmt.Errorf("unable to generate network configuration: %w", err)
	}

	net, err := lvc.NetworkDefineXML(netxml)
	if err != nil {
		return fmt.Errorf("unable to define the vm network: %w", err)
	}

	err = net.SetAutostart(true)
	if err != nil {
		return fmt.Errorf("unable to set the network to autostart: %w", err)
	}

	err = net.Create()
	if err != nil {
		return fmt.Errorf("unable to create the network: %w", err)
	}

	log.Infof("successfully created network '%s' with UUID '%s'", spec.Name, spec.UUID)

	return nil
}
