package agentconfig

import (
	"fmt"
	"os"
	"path/filepath"
)

func CreateBootstrapAgentConfigFile(spec *BootstrapInPlaceAgentConfigSpec, workdir string) error {
	config, err := GetBipAgentConfig(spec)
	if err != nil {
		return fmt.Errorf("unable to generate bootstrap agent config: %w", err)
	}

	if err := os.WriteFile(filepath.Join(workdir, "agent-config.yaml"), []byte(config), 0644); err != nil {
		return fmt.Errorf("unable to create agent-config.yaml: %w", err)
	}

	return nil
}
