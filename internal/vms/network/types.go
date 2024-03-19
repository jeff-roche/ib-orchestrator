package network

import (
	"fmt"
	vmutils "snoman/internal/vms/utils"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"gopkg.in/yaml.v2"
	"libvirt.org/go/libvirtxml"
)

type VirtualMachineNetworkSpec struct {
	Name       string `yaml:"name" validate:"required"`
	UUID       string `yaml:"UUID,omitempty" validate:"omitempty,uuid"`
	BridgeName string `yaml:"bridge" validate:"required"`
	MacAddress string `yaml:"mac_address,omitempty" validate:"omitempty,mac"`
	CIDR       string `yaml:"cidr" validate:"required,cidr"`
	Domain     string `yaml:"domain" validate:"required,fqdn"`
	prefix     string
}

const (
	DEFAULT_NETWORK_NAME        string = "sno-network"
	DEFAULT_NETWORK_BRIDGE_NAME string = "sno-bridge"
	DEFAULT_NETWORK_CIDR        string = "192.168.126.0/24"
	DEFAULT_DOMAIN              string = "sno.rhlocal.com"
)

func GetDefaultVirtualMachineNetworkSpec() *VirtualMachineNetworkSpec {
	prefix, _ := getNetworkPrefixFromCIDR(DEFAULT_NETWORK_CIDR)

	spec := &VirtualMachineNetworkSpec{
		Name:       DEFAULT_NETWORK_NAME,
		UUID:       uuid.NewString(),
		BridgeName: DEFAULT_NETWORK_BRIDGE_NAME,
		MacAddress: getRandomMacAddress(),
		CIDR:       DEFAULT_NETWORK_CIDR,
		Domain:     DEFAULT_DOMAIN,
		prefix:     prefix,
	}

	spec.genHiddenFields()

	return spec
}

func (spec VirtualMachineNetworkSpec) Validate() error {
	err := vmutils.SpecValidator.Struct(spec)
	if err != nil {
		return fmt.Errorf("unable to validate VirtualMachineNetworkSpec: %w", err)
	}

	return nil
}

func (spec *VirtualMachineNetworkSpec) genHiddenFields() {
	spec.prefix, _ = getNetworkPrefixFromCIDR(spec.CIDR)
}

func (spec VirtualMachineNetworkSpec) MarshalYAML() (string, error) {
	if err := spec.Validate(); err != nil {
		return "", err
	}

	data, err := yaml.Marshal(spec)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (spec *VirtualMachineNetworkSpec) UnmarshalYAML(yamlData []byte) error {
	err := yaml.Unmarshal(yamlData, spec)
	if err != nil {
		return fmt.Errorf("unable to parse the spec: %w", err)
	}

	if err := spec.Validate(); err != nil {
		return err
	}

	spec.genHiddenFields()

	return nil
}

func (spec VirtualMachineNetworkSpec) MarshalXML() (string, error) {
	if err := spec.Validate(); err != nil {
		return "", err
	}

	ipprefix, err := strconv.ParseUint(strings.Split(spec.CIDR, "/")[1], 10, 64)
	if err != nil {
		return "", fmt.Errorf("unable to parse IP prefix from CIDR '%s': %w", spec.CIDR, err)
	}

	// Create the net config xml
	netcfg := &libvirtxml.Network{
		Name: spec.Name,
		UUID: spec.UUID,
		Forward: &libvirtxml.NetworkForward{
			Mode: "nat",
		},
		Bridge: &libvirtxml.NetworkBridge{
			Name:  spec.BridgeName,
			STP:   "on",
			Delay: "0",
		},
		MTU: &libvirtxml.NetworkMTU{
			Size: 1500,
		},
		MAC: &libvirtxml.NetworkMAC{
			Address: spec.MacAddress,
		},
		Domain: &libvirtxml.NetworkDomain{
			Name:      spec.Domain,
			LocalOnly: "yes",
		},
		DNS: &libvirtxml.NetworkDNS{
			Enable: "yes",
		},
		IPs: []libvirtxml.NetworkIP{
			{
				Family:  "ipv4",
				Address: strings.Join([]string{spec.prefix, "1"}, "."),
				Prefix:  uint(ipprefix),
				DHCP: &libvirtxml.NetworkDHCP{
					Ranges: []libvirtxml.NetworkDHCPRange{
						{
							Start: strings.Join([]string{spec.prefix, "2"}, "."),
							End:   strings.Join([]string{spec.prefix, "254"}, "."),
						},
					},
				},
			},
		},
	}

	return netcfg.Marshal()
}

func (spec *VirtualMachineNetworkSpec) UnmarshalXML(xmlData []byte) error {
	netcfg := &libvirtxml.Network{}

	if err := netcfg.Unmarshal(string(xmlData)); err != nil {
		return fmt.Errorf("unable to generate a spec from the xml: %w", err)
	}

	// Set the fields
	spec.Name = netcfg.Name
	spec.UUID = netcfg.UUID
	spec.BridgeName = netcfg.Bridge.Name
	spec.MacAddress = netcfg.MAC.Address
	spec.Domain = netcfg.Domain.Name

	// CIDR
	if len(netcfg.IPs) > 0 {
		ipParts := strings.Split(netcfg.IPs[0].Address, ".")[0:3]
		spec.CIDR = fmt.Sprintf("%s.%s.%s.0/%d", ipParts[0], ipParts[1], ipParts[2], netcfg.IPs[0].Prefix)
	} else {
		return fmt.Errorf("unable to determine CIDR. No IP range specified in the XML")
	}

	spec.genHiddenFields() // Calculate any hidden fields from the data retrieved

	return nil
}
