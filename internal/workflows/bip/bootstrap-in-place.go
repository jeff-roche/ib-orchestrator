package bip

import (
	"fmt"
	"os"
	"path/filepath"
	"snoman/internal/biputils"
	"snoman/internal/biputils/agentconfig"
	"snoman/internal/biputils/installconfig"
	"snoman/internal/logger"
	"snoman/internal/vms/machines"

	"go.uber.org/zap"
)

func Run(spec *BootstrapInPlaceSpec) error {
	log := logger.Get()

	// Make sure the workdir exists and create the directory for openshift-install to operate in
	log.Info("validating the working folder exists")
	if _, err := os.Stat(spec.Workdir); err != nil {
		if os.IsNotExist(err) {
			log.Infof("creating working folder at '%s'", spec.Workdir)
			if err := os.MkdirAll(spec.Workdir, 0755); err != nil {
				return fmt.Errorf("unable to create the working folder: %w", err)
			}
		} else {
			return fmt.Errorf("unable to check for the working folder: %w", err)
		}
	}

	// Generate the ISO if needed
	if spec.IsoSpec.IsoPath == "" {
		log.Info("generating installer ISO image")
		if err := generateISO(spec, log); err != nil {
			return fmt.Errorf("unable to generate installer ISO: %w", err)
		}
	}

	// Validate the ISO file exists
	if _, err := os.Stat(spec.IsoSpec.IsoPath); err != nil {
		return fmt.Errorf("could not find the installer iso: %w", err)
	}

	// Create the dnsmasq config
	dnsmasqAddr := fmt.Sprintf("/api.%s.%s/%s", spec.MachineConfig.Name, spec.MachineConfig.Network.Domain, spec.MachineConfig.Network.Hosts[0].IpAddress)
	if err := biputils.CreateDnsmasqConfig(dnsmasqAddr); err != nil {
		return fmt.Errorf("could not create dnsmasq config: %w", err)
	}

	// Create the virtual machine and network
	if err := machines.CreateVirtualMachine(spec.MachineConfig); err != nil {
		return fmt.Errorf("could not create the virtual machine: %w", err)
	}

	return nil
}

func generateISO(spec *BootstrapInPlaceSpec, log *zap.SugaredLogger) error {
	installerWorkdir := filepath.Join(spec.Workdir, "clusterconfig")

	// Remove any previous installer config folder and create a new one
	log.Info("clearing any previous openshift-installer configurations")
	if err := os.RemoveAll(installerWorkdir); err != nil {
		return fmt.Errorf("unable to clear previous bootstrap configs: %w", err)
	}

	log.Infof("creating openshift-install working folder at '%s'", installerWorkdir)
	if err := os.MkdirAll(installerWorkdir, 0755); err != nil {
		return fmt.Errorf("unable to create installer working folder: %w", err)
	}

	// Generate the install-config.yaml
	log.Info("creating install-config.yaml")
	icspec := &installconfig.BootstrapInPlaceInstallConfigSpec{
		BaseDomain:        spec.MachineConfig.Network.Domain,
		ClusterName:       spec.MachineConfig.Name,
		ClusterNetwork:    spec.MachineConfig.Network.ClusterNetworkCIDR,
		ClusterSvcNetwork: spec.MachineConfig.Network.ClusterSvcNetworkCIDR,
		MachineNetwork:    spec.MachineConfig.Network.CIDR,
		InstallDisk:       spec.MachineConfig.Disk.InstallDisk,
		PullSecret:        spec.PullSecret,
		SshPubKey:         spec.PublicKey,
	}

	if err := installconfig.CreateBootstrapInstallConfigFile(icspec, spec.Workdir); err != nil {
		return fmt.Errorf("unable to generate bootstrap install config file: %w", err)
	}

	spec.IsoSpec.InstallConfigPath = filepath.Join(spec.Workdir, "install-config.yaml")

	// Generate the agent-config.yaml
	log.Info("creating agent-config.yaml")

	// We only support 1 host in this workflow
	if len(spec.MachineConfig.Network.Hosts) == 0 {
		return fmt.Errorf("no VM network host configuration provided")
	} else if len(spec.MachineConfig.Network.Hosts) > 1 {
		log.Warn("more than one VM network host config provided. Only the first configuration will be used")
	}

	acspec := &agentconfig.BootstrapInPlaceAgentConfigSpec{
		VmName:    spec.MachineConfig.Name,
		HostIP:    spec.MachineConfig.Network.Hosts[0].IpAddress,
		HostMAC:   spec.MachineConfig.Network.Hosts[0].MacAddress,
		HostRoute: spec.MachineConfig.Network.Hosts[0].IpAddress,
	}

	if err := agentconfig.CreateBootstrapAgentConfigFile(acspec, spec.Workdir); err != nil {
		return fmt.Errorf("unable to generate bootstrap agent config file: %w", err)
	}

	spec.IsoSpec.AgentConfigPath = filepath.Join(spec.Workdir, "agent-config.yaml")

	// Generate the ISO image
	log.Info("generating install iso")
	if err := biputils.GenerateIso(spec.IsoSpec, installerWorkdir); err != nil {
		installerLog := filepath.Join(installerWorkdir, ".openshift_install.log")
		if _, err := os.Stat(installerLog); err == nil {
			log.Debugf("more failure information available in %s", installerLog)
		}

		return fmt.Errorf("unable to generate iso image: %w", err)
	}

	isoFileName := fmt.Sprintf("agent.%s.iso", spec.IsoSpec.OpenshiftArch)
	spec.IsoSpec.IsoPath = filepath.Join(installerWorkdir, isoFileName)

	return nil

}
