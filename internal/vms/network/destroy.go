package network

import (
	"fmt"
	"snoman/internal/logger"
	vmutils "snoman/internal/vms/utils"
)

// Destroy will use the id to delete the network. The id can be a network name or UUID
func Destroy(id string) error {
	log := logger.Get()

	lvc, err := vmutils.GetLibvirtConnection()
	if err != nil {
		return fmt.Errorf("unable to initialize libvirt connection: %w", err)
	}
	defer lvc.Close()

	net := findNetworkByNameOrUUID(id, lvc)
	if net == nil {
		return fmt.Errorf("could not find network with identifier '%s'", id)
	}

	if active, _ := net.IsActive(); active {
		if err := net.Destroy(); err != nil {
			return fmt.Errorf("could not destroy the network: %w", err)
		}
	}

	if err := net.Undefine(); err != nil {
		return fmt.Errorf("could not undefine the network: %w", err)
	}

	log.Infow("successfully deleted network by identifier", "id", id)

	return nil
}
