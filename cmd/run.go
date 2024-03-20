package cmd

import (
	"snoman/internal/workflows/bip"
	"snoman/internal/workflows/ibu"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run the specified workflow",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Fatalf("Error executing run command: %v", ErrResourceTypeNotSpecified)
	},
}

func initRunCmd() {
	rootCmd.AddCommand(runCmd)

	// Subcommands
	runCmd.AddCommand(runBipCmd)
	runCmd.AddCommand(runIbuCmd)
}

// Run Bootstrap In Place
// This uses libvirt and modified local DNS configs to bootstrap an SNO cluster locally
var runBipCmd = &cobra.Command{
	Use:   "bootstrap-in-place",
	Short: "Run bootstrap in place",
	Run: func(cmd *cobra.Command, args []string) {
		if err := bip.Run(); err != nil {
			logger.Errorf("unable to run bootstrap in place: %v", err)
		}
	},
}

// Run Image Based Upgrade
// This uses bootstrap in place to do image based upgrades and installs of Openshift
var runIbuCmd = &cobra.Command{
	Use:   "image-based-upgrade",
	Short: "Run image based upgrade against a cluster",
	Run: func(cmd *cobra.Command, args []string) {
		if err := ibu.Run(); err != nil {
			logger.Errorf("unable to run image based upgrade: %v", err)
		}
	},
}
