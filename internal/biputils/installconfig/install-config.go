package installconfig

import (
	"fmt"
	"os"
	"path/filepath"
)

func CreateBootstrapInstallConfigFile(spec *BootstrapInPlaceInstallConfigSpec, workdir string) error {
	config, err := GetBipInstallConfig(spec)
	if err != nil {
		return fmt.Errorf("unable to generate bootstrap install config: %w", err)
	}

	if err := os.WriteFile(filepath.Join(workdir, "install-config.yaml"), []byte(config), 0644); err != nil {
		return fmt.Errorf("unable to create install-config.yaml: %w", err)
	}

	return nil
}
