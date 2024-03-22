package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"snoman/internal/biputils"
	"snoman/internal/vms/machines"
	"snoman/internal/vms/network"

	"github.com/spf13/cobra"
)

var ErrNotImplemented = fmt.Errorf("command not yet implemented")
var ErrResourceTypeNotSpecified = fmt.Errorf("you must specify the resource type")

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create the specified resource",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Fatalf("Error executing create command: %v", ErrResourceTypeNotSpecified)
	},
}

func initCreateCmd() {
	rootCmd.AddCommand(createCmd)

	// Subcommands
	createCmd.AddCommand(createVmCmd)

	// Create network
	createCmd.AddCommand(createNetCmd)
	createNetCmd.Flags().String("from-xml", "", "Path to the XML file to use for network creation")
	createNetCmd.Flags().String("from", "", "Path to the spec file to use for network creation")

	// Boostrap ISO
	createCmd.AddCommand(createBootstrapIsoCmd)
	wd, _ := os.Getwd()
	wd, _ = filepath.Abs(wd)
	createBootstrapIsoCmd.Flags().String("agent-config-file", "", "[Required] The path to the agent-config.yaml that will be used to generate the iso")
	createBootstrapIsoCmd.Flags().String("install-config-file", "", "[Required] The path to the install-config.yaml that will be used to generate the iso")
	createBootstrapIsoCmd.Flags().StringP("workdir", "w", wd, "The working folder to generate the iso in")
	createBootstrapIsoCmd.Flags().String("abi-path", biputils.DEFAULT_ABI_PATH, fmt.Sprintf("Path to the agent based installer (default: %s)", biputils.DEFAULT_ABI_PATH))
	createBootstrapIsoCmd.Flags().String("custom-manifests", "", "Path to the folder containing custom manifests to be added to the ISO")
	createBootstrapIsoCmd.Flags().String("ocp", biputils.DEFAULT_OPENSHIFT_VERSION, fmt.Sprintf("The version of OCP to use in the ISO (default: %s)", biputils.DEFAULT_OPENSHIFT_VERSION))
	createBootstrapIsoCmd.Flags().String("arch", biputils.DEFAULT_OPENSHIFT_ARCH, fmt.Sprintf("The OCP arch to use in the ISO (default: %s)", biputils.DEFAULT_OPENSHIFT_ARCH))
	createBootstrapIsoCmd.Flags().String("release-image", "", "The pull spec of the OCP release image to use for creating the ISO")
}

// Create VM
var createVmCmd = &cobra.Command{
	Use:   "vm",
	Short: "Create an OCP VM locally",
	Run: func(cmd *cobra.Command, args []string) {
		spec := machines.GetDefaultVirtualMachineSpec()
		host := spec.Network.Hosts[0]
		spec.Network.Hosts = nil

		err := machines.CreateVirtualMachine(spec)
		if err != nil {
			logger.Error(err)
		}

		if err := network.AddHostToNetwork(spec.Network.UUID, &host); err != nil {
			logger.Error(err)
		}
	},
}

// Create VM Network
var createNetCmd = &cobra.Command{
	Use:   "network",
	Short: "Create a libvirt network to be used for an OCP VM",
	Long: `
	Create a libvirt network to be used for SNO cluster running inside a virtual machine
	
	If --from-xml or --from are not specified, the default configuration will be used

	**Note**: this can be skipped if the "network" stanza of the virtual machine spec is specified.
	The "network" stanza is defined in the default VM spec.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		source, _ := cmd.Flags().GetString("from")
		xmlSource, _ := cmd.Flags().GetString("from-xml")
		var spec *network.VirtualMachineNetworkSpec

		if source != "" {
			data, err := os.ReadFile(source)
			if err != nil {
				logger.Fatalf("unable to read the spec file: %w", err)
			}

			spec = &network.VirtualMachineNetworkSpec{}
			spec.UnmarshalYAML(data)
		} else if xmlSource != "" {
			data, err := os.ReadFile(xmlSource)
			if err != nil {
				logger.Fatalf("unable to read the xml file: %w", err)
			}

			spec = &network.VirtualMachineNetworkSpec{}
			spec.UnmarshalXML(data)
		} else {
			spec = network.GetDefaultVirtualMachineNetworkSpec()
		}

		err := network.Create(spec)
		if err != nil {
			logger.Error(err)
		}
	},
}

// Create BIP ISO
var createBootstrapIsoCmd = &cobra.Command{
	Use:   "boostrap-iso",
	Short: "Create a boostrap ready OCP ISO using the agent based installer",
	Run: func(cmd *cobra.Command, args []string) {
		spec := &biputils.BootstrapInPlaceIsoSpec{}

		// Required fields
		spec.AgentConfigPath, _ = cmd.Flags().GetString("agent-config-file")
		if spec.AgentConfigPath == "" {
			logger.Fatal("argument agent-config-file is required when generating a bootstap ISO")
		}

		spec.InstallConfigPath, _ = cmd.Flags().GetString("install-config-file")
		if spec.InstallConfigPath == "" {
			logger.Fatal("argument install-config-file is required when generating a bootstap ISO")
		}

		// Optional fields
		workdir, _ := cmd.Flags().GetString("workdir")
		spec.AbiPath, _ = cmd.Flags().GetString("abi-path")
		spec.CustomManifestDir, _ = cmd.Flags().GetString("custom-manifests")
		spec.OpenshiftVersion, _ = cmd.Flags().GetString("ocp")
		spec.OpenshiftArch, _ = cmd.Flags().GetString("arch")
		spec.ReleaseImage, _ = cmd.Flags().GetString("release-image")

		if err := biputils.GenerateIso(spec, workdir); err != nil {
			logger.Fatalf("unable to generate bootstrap iso: %v", err)
		}

		logger.Info("iso generation complete")
	},
}
