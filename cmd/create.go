package cmd

import (
	"fmt"
	"snoman/internal/vms/machines"
	"snoman/internal/vms/network"

	"github.com/spf13/cobra"
)

var ErrNotImplemented = fmt.Errorf("command not yet implemented")
var ErrResourcetTypeNotSpecified = fmt.Errorf("you must specify the resource type")

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create the specified resource",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Fatalf("Error executing create command: %v", ErrResourcetTypeNotSpecified)
	},
}

func initCreateCmd() {
	rootCmd.AddCommand(createCmd)

	// Subcommands
	createCmd.AddCommand(createVmCmd)
	createCmd.AddCommand(createNetCmd)
}

// Create VM
var createVmCmd = &cobra.Command{
	Use:   "vm",
	Short: "Create an OCP VM locally",
	Run:   createVM,
}

var createVM = func(cmd *cobra.Command, args []string) {
	err := machines.CreateVirtualMachine(machines.GetDefaultVirtualMachineSpec())
	if err != nil {
		logger.Error(err)
	}
}

// Create VM Network
var createNetCmd = &cobra.Command{
	Use:   "network",
	Short: "Create a libvirt network to be used for an OCP VM",
	Long: `
	Create a libvirt network to be used for SNO cluster running inside a virtual machine
	
	**Note**: this can be skipped if the "network" stanza of the virtual machine spec is specified.
	The "network" stanza is defined in the default VM spec.
	`,
	Run: createNet,
}

var createNet = func(cmd *cobra.Command, args []string) {
	err := network.Create(network.GetDefaultVirtualMachineNetworkSpec())
	if err != nil {
		logger.Error(err)
	}
}
