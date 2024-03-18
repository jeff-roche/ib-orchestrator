package cmd

import (
	"fmt"
	"snoman/internal/vmmanager"

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
}

var createVmCmd = &cobra.Command{
	Use:   "vm",
	Short: "Create an OCP VM locally",
	Run:   createVM,
}

var createVM = func(cmd *cobra.Command, args []string) {
	err := vmmanager.CreateVirtualMachine()
	if err != nil {
		logger.Error(err)
	}
}
