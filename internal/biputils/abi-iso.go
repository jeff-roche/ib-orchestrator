package biputils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"snoman/internal/logger"
)

const installer_image_override_env = "OPENSHIFT_INSTALL_RELEASE_IMAGE_OVERRIDE"

func GenerateIso(spec *BootstrapInPlaceIsoSpec, workdir string) error {
	if err := spec.fillAndValidateIsoGenFields(); err != nil {
		return fmt.Errorf("unable to validate the required fields for iso generation: %w", err)
	}

	// Move the agent config to the working directory
	agentConfigData, err := os.ReadFile(spec.AgentConfigPath)
	if err != nil {
		return fmt.Errorf("could not read agent config '%s': %w", spec.AgentConfigPath, err)
	}

	agentConfigPath := filepath.Join(workdir, "agent-config.yaml")
	if err := os.WriteFile(agentConfigPath, agentConfigData, 0644); err != nil {
		return fmt.Errorf("could not write %s: %w", agentConfigPath, err)
	}

	// Move the install config to the working directory
	installConfigData, err := os.ReadFile(spec.InstallConfigPath)
	if err != nil {
		return fmt.Errorf("could not read install config '%s': %w", spec.InstallConfigPath, err)
	}

	installConfigPath := filepath.Join(workdir, "install-config.yaml")
	if err := os.WriteFile(installConfigPath, installConfigData, 0644); err != nil {
		return fmt.Errorf("could not write %s: %w", installConfigPath, err)
	}

	// Run openshift-install to generate the iso image
	// ${INSTALLER_BIN} agent create image --log-level debug --dir="${INSTALLER_WORKDIR}"

	args := []string{
		"agent", "create", "image",
		"--log-level", "debug",
		fmt.Sprintf("--dir=%s", workdir),
	}

	// Create the command
	isoGenCmd := exec.Command(spec.AbiPath, args...)
	isoGenCmd.Env = append(os.Environ(), fmt.Sprintf("%s=%s", installer_image_override_env, spec.ReleaseImage))

	out, err := isoGenCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error running the agent based installer: %w", err)
	}

	if err := logger.WriteLogFile(out, workdir, "openshift-install-create-image.log"); err != nil {
		return fmt.Errorf("unable to write openshift-install logs: %w", err)
	}

	return nil
}
