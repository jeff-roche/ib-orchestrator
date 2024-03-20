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

type VMNet_DHCP_Host struct {
	Name       string `yaml:"name" validate:"required"`
	MacAddress string `yaml:"mac_address,omitempty" validate:"omitempty,mac"`
	IpAddress  string `yaml:"ip_address,omitempty" validate:"ip,omitempty"`
}

type VirtualMachineNetworkSpec struct {
	Name       string            `yaml:"name" validate:"required"`
	UUID       string            `yaml:"UUID,omitempty" validate:"omitempty,uuid"`
	BridgeName string            `yaml:"bridge" validate:"required"`
	MacAddress string            `yaml:"mac_address,omitempty" validate:"omitempty,mac"`
	CIDR       string            `yaml:"cidr" validate:"required,cidr"`
	Domain     string            `yaml:"domain" validate:"required,fqdn"`
	Hosts      []VMNet_DHCP_Host `yaml:"hosts,omitempty" validate:"omitempty"`
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

	spec.genAdditionalFields()

	spec.addDefaultHost()

	return spec
}

func (spec VirtualMachineNetworkSpec) Validate() error {
	err := vmutils.SpecValidator.Struct(spec)
	if err != nil {
		return fmt.Errorf("unable to validate VirtualMachineNetworkSpec: %w", err)
	}

	return nil
}

func (spec *VirtualMachineNetworkSpec) genAdditionalFields() {
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

	spec.genAdditionalFields()

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

	if len(spec.Hosts) > 0 {
		netcfg.IPs[0].DHCP.Hosts = make([]libvirtxml.NetworkDHCPHost, 0, len(spec.Hosts))

		for _, host := range spec.Hosts {
			netcfg.IPs[0].DHCP.Hosts = append(netcfg.IPs[0].DHCP.Hosts, libvirtxml.NetworkDHCPHost{
				Name: host.Name,
				IP:   host.IpAddress,
				MAC:  host.MacAddress,
			})
		}
	}

	return netcfg.Marshal()
}

func (spec *VirtualMachineNetworkSpec) UnmarshalXML(xmlData []byte) error {
	netcfg := &libvirtxml.Network{}

	if err := netcfg.Unmarshal(string(xmlData)); err != nil {
		return fmt.Errorf("unable to parse libvirt xml: %w", err)
	}

	if err := spec.fromLibvirtxml(netcfg); err != nil {
		return fmt.Errorf("unable to generate spec from xml: %w", err)
	}

	return nil
}

func (spec *VirtualMachineNetworkSpec) fromLibvirtxml(net *libvirtxml.Network) error {
	// Set the fields
	spec.Name = net.Name
	spec.UUID = net.UUID
	spec.BridgeName = net.Bridge.Name
	spec.MacAddress = net.MAC.Address
	spec.Domain = net.Domain.Name

	// CIDR
	if len(net.IPs) > 0 {
		ipParts := strings.Split(net.IPs[0].Address, ".")[0:3]
		spec.CIDR = fmt.Sprintf("%s.%s.%s.0/%d", ipParts[0], ipParts[1], ipParts[2], net.IPs[0].Prefix)
	} else {
		return fmt.Errorf("unable to determine CIDR. No IP range specified in the XML")
	}

	// If IPs are set and it has at least 1 host
	if len(net.IPs) > 0 && len(net.IPs[0].DHCP.Hosts) > 0 {
		spec.Hosts = make([]VMNet_DHCP_Host, 0, len(net.IPs[0].DHCP.Hosts))

		for _, host := range net.IPs[0].DHCP.Hosts {
			spec.Hosts = append(spec.Hosts, VMNet_DHCP_Host{
				Name:       host.Name,
				IpAddress:  host.IP,
				MacAddress: host.MAC,
			})
		}
	}

	spec.genAdditionalFields() // Calculate any hidden fields from the data retrieved

	return nil
}

func (spec *VirtualMachineNetworkSpec) addDefaultHost() {
	spec.Hosts = append(spec.Hosts, VMNet_DHCP_Host{
		Name:       "example_host",
		IpAddress:  fmt.Sprintf("%s.%d", spec.prefix, 10),
		MacAddress: fmt.Sprintf(spec.MacAddress),
	})
}
