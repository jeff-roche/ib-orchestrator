package cmd

import (
	"fmt"
	"os"
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
