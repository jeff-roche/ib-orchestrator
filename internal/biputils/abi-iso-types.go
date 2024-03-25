package biputils

import (
	"fmt"
	vmutils "snoman/internal/vms/utils"

	"gopkg.in/yaml.v2"
)

type BootstrapInPlaceIsoSpec struct {
	IsoPath           string `yaml:"iso_path,omitempty" validate:"omitempty,file"`
	AbiPath           string `yaml:"agent_based_installer_path,omitempty" validate:"omitempty,file"`
	CustomManifestDir string `yaml:"custom_manifests_path,omitempty" validate:"omitempty,dir"`
	OpenshiftVersion  string `yaml:"ocp_release_version" validate:"omitempty,semver"`
	OpenshiftArch     string `yaml:"ocp_release_arch" validate:"omitempty"`
	ReleaseImage      string `yaml:"release_image,omitempty" validate:"omitempty"`
	AgentConfigPath   string `yaml:"agent_config_file" validate:"file"`
	InstallConfigPath string `yaml:"install_config_file" validate:"file"`
}

const (
	DEFAULT_ABI_PATH          = "/bin/openshift-install"
	DEFAULT_RELEASE_IMAGE     = "quay.io/openshift-release-dev/ocp-release"
	DEFAULT_OPENSHIFT_VERSION = "4.15.1"
	DEFAULT_OPENSHIFT_ARCH    = "x86_64"
)

// fillAndValidateIsoGenFields will populate any needed empty fields with defaults and then validate the struct
func (spec *BootstrapInPlaceIsoSpec) FillAndValidateIsoGenFields() error {
	// If an ISO path was given, we don't need anything else
	if spec.IsoPath != "" {
		if err := vmutils.SpecValidator.Struct(spec); err != nil {
			return fmt.Errorf("unable to validate VirtualMachineSpec: %w", err)
		}

		return nil
	}

	// We need the installer path to generate the ISO
	if spec.AbiPath == "" {
		spec.AbiPath = DEFAULT_ABI_PATH
	}

	// The release image is used by ABI to generate the ISO
	if spec.ReleaseImage == "" {
		if spec.OpenshiftVersion == "" {
			spec.OpenshiftVersion = DEFAULT_OPENSHIFT_VERSION
		}

		if spec.OpenshiftArch == "" {
			spec.OpenshiftArch = DEFAULT_OPENSHIFT_ARCH
		}

		spec.ReleaseImage = fmt.Sprintf("%s:%s-%s", DEFAULT_RELEASE_IMAGE, spec.OpenshiftVersion, spec.OpenshiftArch)
	}

	// Validate the fields that exist
	if err := vmutils.SpecValidator.Struct(spec); err != nil {
		return fmt.Errorf("unable to validate VirtualMachineSpec: %w", err)
	}

	return nil
}

func (spec *BootstrapInPlaceIsoSpec) MarshalYAML() (string, error) {
	if err := spec.FillAndValidateIsoGenFields(); err != nil {
		return "", err
	}

	data, err := yaml.Marshal(spec)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (spec *BootstrapInPlaceIsoSpec) UnmarshalYAML(yamlData []byte) error {
	err := yaml.Unmarshal(yamlData, spec)
	if err != nil {
		return fmt.Errorf("unable to parse the spec: %w", err)
	}

	if err := spec.FillAndValidateIsoGenFields(); err != nil {
		return err
	}

	return nil
}
