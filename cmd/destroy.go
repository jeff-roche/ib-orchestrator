package cmd

import (
	"snoman/internal/vms/network"

	"github.com/spf13/cobra"
)

var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy the specified resource",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Fatalf("Error executing destroy command: %v", ErrResourceTypeNotSpecified)
	},
}

func initDestroyCmd() {
	rootCmd.AddCommand(destroyCmd)

	// Subcommands
	destroyCmd.AddCommand(destroyNetCmd)
}

// Destroy VM Network
var destroyNetCmd = &cobra.Command{
	Use:   "network [name or uuid]",
	Short: "Destroy a libvirt network by name or UUID",
	Long: `
	Destroy a libvirt network by name or UUID

	if no name UUID is provided, the default network name will be used
	`,
	Run: func(cmd *cobra.Command, args []string) {
		netname := network.DEFAULT_NETWORK_NAME

		if len(args) > 0 {
			netname = args[0]
		}

		err := network.Destroy(netname)
		if err != nil {
			logger.Error(err)
		}
	},
}
