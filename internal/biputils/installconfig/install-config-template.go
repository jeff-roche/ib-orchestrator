package installconfig

import (
	"fmt"
	"strings"
	"text/template"
)

func GetBipInstallConfig(data *BootstrapInPlaceInstallConfigSpec) (string, error) {
	ictemplate, err := template.New("bip-install-config").Parse(installConfigTemplate)
	if err != nil {
		return "", fmt.Errorf("unable to parse install config template: %w", err)
	}

	contents := new(strings.Builder)

	if err := ictemplate.Execute(contents, data); err != nil {
		return "", fmt.Errorf("unable to generate install config contents: %w", err)
	}

	return contents.String(), nil
}

var installConfigTemplate = `apiVersion: v1
baseDomain: {{ .BaseDomain }}
compute:
- architecture: amd64
  hyperthreading: Enabled
  name: worker
  platform: {}
  replicas: 0
controlPlane:
  architecture: amd64
  hyperthreading: Enabled
  name: master
  platform: {}
  replicas: 1
metadata:
  creationTimestamp: null
  name: {{ .ClusterName }}
networking:
  clusterNetwork:
  - cidr: {{ .ClusterNetwork }}
    hostPrefix: 23
  machineNetwork:
  - cidr: {{ .MachineNetwork }}
  networkType: OVNKubernetes
  serviceNetwork:
  - {{ .ClusterSvcNetwork }}
platform:
  none: {}
BootstrapInPlace:
  InstallationDisk: {{ .InstallDisk }}
publish: External
pullSecret: |
        {{ .PullSecret }}
sshKey: |
        {{ .SshPubKey }}`
