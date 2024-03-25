package cmd

import (
	"os"
	"path/filepath"
	"snoman/internal/biputils"
	"snoman/internal/biputils/secrets"
	"snoman/internal/vms/machines"
	"snoman/internal/workflows/bip"

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

	// Run Bootstrap In Place Command
	runCmd.AddCommand(runBipCmd)

	wd, _ := os.Getwd()
	wd = filepath.Join(wd, "workdir")
	wd, _ = filepath.Abs(wd)

	runBipCmd.Flags().String("ssh-pub-key", "", "[Required] Path to the file containing the cluster ssh public key")
	runBipCmd.Flags().String("vm-config", "", "Path to the configuration yaml for the virtual machine")
	runBipCmd.Flags().String("pull-secret-file", "", "Path to the file containing the cluster pull secret. If left empty the PULL_SECRET env variable will be used")
	runBipCmd.Flags().String("iso-file", "", "Path to the installer iso file to use for the VM")
	runBipCmd.Flags().String("iso-config", "", "Path to the configuration yaml for the iso file")
	runBipCmd.Flags().StringP("workdir", "w", wd, "The working folder to generate any required files in")

	//runCmd.AddCommand(runIbuCmd)
}

// Run Bootstrap In Place
// This uses libvirt and modified local DNS configs to bootstrap an SNO cluster locally
var runBipCmd = &cobra.Command{
	Use:   "bootstrap-in-place",
	Short: "Run bootstrap in place",
	Run: func(cmd *cobra.Command, args []string) {
		spec := &bip.BootstrapInPlaceSpec{}
		var err error

		// SSH Public Key
		sshPubKeyFile, _ := cmd.Flags().GetString("ssh-pub-key")
		if sshPubKeyFile == "" {
			logger.Fatalf("path to your ssh public key file is required")
		}

		spec.PublicKey, err = secrets.GetSshPublicKey(sshPubKeyFile)
		if err != nil {
			logger.Fatalf("unable to parse ssh public key: %v", err)
		}

		// Pull Secret
		pullSecretFile, _ := cmd.Flags().GetString("pull-secret-file")
		spec.PullSecret, err = secrets.GetPullSecret(pullSecretFile)
		if err != nil {
			logger.Fatalf("unable to parse pull secret: %v", err)
		}

		// Virtual Machine Configuration
		vmConfigFile, _ := cmd.Flags().GetString("vm-config")
		if vmConfigFile == "" {
			logger.Info("using default virtual machine configuration")
			spec.MachineConfig = machines.GetDefaultVirtualMachineSpec()
		} else {
			data, err := os.ReadFile(vmConfigFile)
			if err != nil {
				logger.Fatalf("unable to read the virtual machine config file: %v", err)
			}

			spec.MachineConfig = &machines.VirtualMachineSpec{}
			spec.MachineConfig.UnmarshalYAML(data)
		}

		// ISO Configuration
		spec.IsoSpec = &biputils.BootstrapInPlaceIsoSpec{}
		spec.IsoSpec.IsoPath, _ = cmd.Flags().GetString("iso-file")
		if spec.IsoSpec.IsoPath == "" {
			// Check for an ISO config and generate one if needed
			isoConfigFile, _ := cmd.Flags().GetString("iso-config")
			if isoConfigFile == "" {
				logger.Info("using default ISO configuration")
				spec.IsoSpec.FillAndValidateIsoGenFields()
			} else {
				data, err := os.ReadFile(isoConfigFile)
				if err != nil {
					logger.Fatalf("unable to read the iso config file: %v", err)
				}

				if err := spec.IsoSpec.UnmarshalYAML(data); err != nil {
					logger.Fatalf("unable to unmarshal iso config file: %v", err)
				}
			}
		}

		spec.Workdir, _ = cmd.Flags().GetString("workdir")

		if err := bip.Run(spec); err != nil {
			logger.Errorf("unable to run bootstrap in place: %v", err)
		}
	},
}

// Run Image Based Upgrade
// This uses bootstrap in place to do image based upgrades and installs of Openshift
/*var runIbuCmd = &cobra.Command{
	Use:   "image-based-upgrade",
	Short: "Run image based upgrade against a cluster",
	Run: func(cmd *cobra.Command, args []string) {
		if err := ibu.Run(); err != nil {
			logger.Errorf("unable to run image based upgrade: %v", err)
		}
	},
}*/
