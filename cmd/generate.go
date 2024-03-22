package cmd

import (
	"fmt"
	"os"
	bipagentconfig "snoman/internal/biputils/agentconfig"
	bipinstallconfig "snoman/internal/biputils/installconfig"
	"snoman/internal/biputils/secrets"
	"snoman/internal/vms/machines"
	"snoman/internal/vms/network"

	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate the specified resource config",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Fatalf("Error executing the generate command: %v", ErrResourceTypeNotSpecified)
	},
}

func initGenerateCmd() {
	rootCmd.AddCommand(generateCmd)

	// Subcommands
	// Generate VM configuration
	generateCmd.AddCommand(generateVmCmd)
	generateVmCmd.Flags().String("from", "", "Path to the spec file to use for spec generation. This will just sanitize the input.")

	// Generate network configuration
	generateCmd.AddCommand(generateNetCmd)
	generateNetCmd.Flags().String("from-net", "", "Name or UUID of the libvirt network to use as a base")
	generateNetCmd.Flags().String("from-xml", "", "Path to the XML file to use for spec generation")
	generateNetCmd.Flags().String("from", "", "Path to the spec file to use for spec generation. This will just sanitize the input.")
	generateNetCmd.Flags().Bool("xml", false, "Generate the libvirt xml config for the network")

	// Generate BIP ABI Config
	generateCmd.AddCommand(generateBipAgentConfigCmd)
	generateBipAgentConfigCmd.Flags().String("vm-name", "", "Name of the VM this config will be used for")
	generateBipAgentConfigCmd.Flags().String("host-ip", "", "IP address that the resulting host will use")
	generateBipAgentConfigCmd.Flags().String("host-mac", "", "MAC address that the resulting host will use")
	generateBipAgentConfigCmd.Flags().String("host-route", "", "Route that the resulting host will use")

	// Generate BIP install config
	generateCmd.AddCommand(generateBipInstallConfigCmd)
	generateBipInstallConfigCmd.Flags().String("base-domain", "", "The base domain that will be used for the network")
	generateBipInstallConfigCmd.Flags().String("cluster-name", "", "The name of the cluster")
	generateBipInstallConfigCmd.Flags().String("cluster-network-cidr", "", "The cluster cidr")
	generateBipInstallConfigCmd.Flags().String("cluster-svc-network-cidr", "", "The cluster service cidr")
	generateBipInstallConfigCmd.Flags().String("machine-network-cidr", "", "The machine network cidr")
	generateBipInstallConfigCmd.Flags().String("install-disk", "", "The install disk path")
	generateBipInstallConfigCmd.Flags().String("pull-secret-file", "", "Path to the file containing the cluster pull secret. If left empty the PULL_SECRET env variable will be used")
	generateBipInstallConfigCmd.Flags().String("ssh-pub-key", "", "[Required] Path to the file containing the cluster ssh public key")

}

// Generate VM Network
var generateVmCmd = &cobra.Command{
	Use:   "vm",
	Short: "Generate a libvirt machine spec",
	Long: `
	Generate a libvirt machine spec

	If --from is not specified, the default machine configuration will be used
	`,
	Run: func(cmd *cobra.Command, args []string) {
		source, _ := cmd.Flags().GetString("from")

		var spec *machines.VirtualMachineSpec
		if source != "" {
			data, err := os.ReadFile(source)
			if err != nil {
				logger.Fatalf("unable to read the spec file: %w", err)
			}

			spec = &machines.VirtualMachineSpec{}

			if err := spec.UnmarshalYAML(data); err != nil {
				logger.Fatalf("unable to parse the provided machine spec: %w", err)
			}
		} else {
			spec = machines.GetDefaultVirtualMachineSpec()
		}

		output, err := spec.MarshalYAML()
		if err != nil {
			logger.Fatalf("unable to generate the machine spec: %w", err)
		}

		fmt.Println(output)
	},
}

// Generate VM Network
var generateNetCmd = &cobra.Command{
	Use:   "network",
	Short: "Generate a libvirt network spec",
	Long: `
	Generate a libvirt network spec

	If --from-xml or --from are not specified, the default configuration will be used
	`,
	Run: func(cmd *cobra.Command, args []string) {
		outputXML, _ := cmd.Flags().GetBool("xml")
		source, _ := cmd.Flags().GetString("from")
		xmlSource, _ := cmd.Flags().GetString("from-xml")
		netSource, _ := cmd.Flags().GetString("from-net")

		var spec *network.VirtualMachineNetworkSpec
		if source != "" {
			data, err := os.ReadFile(source)
			if err != nil {
				logger.Fatalf("unable to read the spec file: %v", err)
			}

			spec = &network.VirtualMachineNetworkSpec{}
			spec.UnmarshalYAML(data)
		} else if xmlSource != "" {
			data, err := os.ReadFile(xmlSource)
			if err != nil {
				logger.Fatalf("unable to read the xml file: %v", err)
			}

			spec = &network.VirtualMachineNetworkSpec{}
			spec.UnmarshalXML(data)
		} else if netSource != "" {
			var err error
			spec, err = network.Find(netSource)
			if err != nil {
				logger.Fatalf("could not generate spec from network: %v", err)
			}
		} else {
			spec = network.GetDefaultVirtualMachineNetworkSpec()
		}

		var output string
		var err error
		if outputXML {
			output, err = spec.MarshalXML()
		} else {
			output, err = spec.MarshalYAML()
		}

		if err != nil {
			logger.Fatalf("unable to generate spec: %w", err)
		}

		fmt.Println(output)
	},
}

// Generate Bootstrap In Place ABI Config
var generateBipAgentConfigCmd = &cobra.Command{
	Use:   "bootstrap-agent-config",
	Short: "Generate an agent based agent-config for BIP",
	Long:  "Generate an agent based agent-config that is specific to bootstrap in place",
	Run: func(cmd *cobra.Command, args []string) {
		srcspec := machines.GetDefaultVirtualMachineSpec()
		acspec := &bipagentconfig.BootstrapInPlaceAgentConfigSpec{}

		// VM Name
		acspec.VmName = srcspec.Name
		if name, _ := cmd.Flags().GetString("vm-name"); name != "" {
			acspec.VmName = name
		}

		// Host IP
		acspec.HostIP = srcspec.Network.Hosts[0].IpAddress
		if hostip, _ := cmd.Flags().GetString("host-ip"); hostip != "" {
			acspec.HostIP = hostip
		}

		// Host MAC
		acspec.HostMAC = srcspec.Network.Hosts[0].MacAddress
		if hostmac, _ := cmd.Flags().GetString("host-mac"); hostmac != "" {
			srcspec.Network.Hosts[0].MacAddress = hostmac
		}

		// Host Route
		acspec.HostRoute = srcspec.Network.Hosts[0].IpAddress
		if hostroute, _ := cmd.Flags().GetString("host-route"); hostroute != "" {
			acspec.HostRoute = hostroute
		}

		config, err := bipagentconfig.GetBipAgentConfig(acspec)
		if err != nil {
			logger.Fatalf("unable to generate bootstrap agent config: %v", err)
		}

		fmt.Println(config)
	},
}

// Generate Bootstrap In Place install config
var generateBipInstallConfigCmd = &cobra.Command{
	Use:   "bootstrap-install-config",
	Short: "Generate an agent based install-config for BIP",
	Long:  "Generate an agent based install-config that is specific to bootstrap in place",
	Run: func(cmd *cobra.Command, args []string) {
		srcspec := machines.GetDefaultVirtualMachineSpec()
		icspec := &bipinstallconfig.BootstrapInPlaceInstallConfigSpec{
			BaseDomain:        srcspec.Network.Domain,
			ClusterName:       srcspec.Name,
			ClusterNetwork:    srcspec.Network.ClusterNetworkCIDR,
			ClusterSvcNetwork: srcspec.Network.ClusterSvcNetworkCIDR,
			MachineNetwork:    srcspec.Network.CIDR,
			InstallDisk:       srcspec.Disk.InstallDisk,
		}

		if basedomain, _ := cmd.Flags().GetString("base-domain"); basedomain != "" {
			icspec.BaseDomain = basedomain
		}

		if clustername, _ := cmd.Flags().GetString("cluster-name"); clustername != "" {
			icspec.ClusterName = clustername
		}

		if clusterNetworkCIDR, _ := cmd.Flags().GetString("cluster-network-cidr"); clusterNetworkCIDR != "" {
			icspec.ClusterNetwork = clusterNetworkCIDR
		}

		if clusteSvcNetworkCIDR, _ := cmd.Flags().GetString("cluster-svc-network-cidr"); clusteSvcNetworkCIDR != "" {
			icspec.ClusterSvcNetwork = clusteSvcNetworkCIDR
		}

		if machineNetworkCIDR, _ := cmd.Flags().GetString("machine-network-cidr"); machineNetworkCIDR != "" {
			icspec.ClusterNetwork = machineNetworkCIDR
		}

		if installdisk, _ := cmd.Flags().GetString("install-disk"); installdisk != "" {
			icspec.InstallDisk = installdisk
		}

		// Pull Secret
		pullSecretFile, _ := cmd.Flags().GetString("pull-secret-file")

		var err error
		icspec.PullSecret, err = secrets.GetPullSecret(pullSecretFile)
		if err != nil {
			logger.Fatalf("unable to parse pull secret: %v", err)
		}

		// SSH Public Key
		sshPubKeyFile, _ := cmd.Flags().GetString("ssh-pub-key")
		icspec.SshPubKey, err = secrets.GetSshPublicKey(sshPubKeyFile)
		if err != nil {
			logger.Fatalf("unable to parse ssh public key: %v", err)
		}

		config, err := bipinstallconfig.GetBipInstallConfig(icspec)
		if err != nil {
			logger.Fatalf("unable to generate bootstrap agent config: %v", err)
		}

		fmt.Println(config)
	},
}
