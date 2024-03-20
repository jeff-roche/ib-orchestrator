package network

import (
	"fmt"
	"snoman/internal/logger"
	vmutils "snoman/internal/vms/utils"

	"libvirt.org/go/libvirt"
	"libvirt.org/go/libvirtxml"
)

// AddHostToNetwork will go find the active network by name or uuid (netid) and add the specified host config
func AddHostToNetwork(netid string, hostspec *VMNet_DHCP_Host) error {
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

	net := findNetworkByNameOrUUID(netid, lvc)
	if net == nil {
		return fmt.Errorf("could not find libvirt network by identifier '%s'", netid)
	}

	lvhost := &libvirtxml.NetworkDHCPHost{
		Name: hostspec.Name,
		IP:   hostspec.IpAddress,
		MAC:  hostspec.MacAddress,
	}

	xml, err := lvhost.Marshal()
	if err != nil {
		return fmt.Errorf("could not generate host xml: %w", err)
	}

	err = net.Update(
		libvirt.NETWORK_UPDATE_COMMAND_ADD_LAST,
		libvirt.NETWORK_SECTION_IP_DHCP_HOST,
		0,
		xml,
		libvirt.NETWORK_UPDATE_AFFECT_LIVE)

	if err != nil {
		return fmt.Errorf("could not update network with id '%s': %w", netid, err)
	}

	log.Infow("successfully added host to network", "network", netid, "host", hostspec)

	return nil
}
