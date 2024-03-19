package network

import (
	"crypto/rand"
	"fmt"
	"net"
	"snoman/internal/vms/utils"
	"strings"

	"libvirt.org/go/libvirt"
)

const (
	maclocalbit     = 0b10
	macmulticastbit = 0b1
)

func getRandomMacAddress() string {
	buf := make([]byte, 6)
	var mac net.HardwareAddr

	rand.Read(buf)

	// Set the local and multicast bits
	buf[0] = buf[0]&^macmulticastbit | maclocalbit

	mac = append(mac, buf...)

	return mac.String()
}

func getNetworkPrefixFromCIDR(cidr string) (string, error) {
	err := utils.SpecValidator.Var(cidr, "cidr")
	if err != nil {
		return "", fmt.Errorf("can not get network prefix from CIDR \"%v\": %w", cidr, err)
	}

	return strings.Join(strings.Split(cidr, ".")[0:3], "."), nil
}

// findNetworkByNameOrUUID will try to find the network and return nil if the network could not be found
func findNetworkByNameOrUUID(id string, lvc *libvirt.Connect) (net *libvirt.Network) {
	net, _ = lvc.LookupNetworkByName(id)

	if net == nil {
		net, _ = lvc.LookupNetworkByUUIDString(id)
	}

	return
}
