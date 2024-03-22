package agentconfig

import (
	"fmt"
	"strings"
	"text/template"
)

func GetBipAgentConfig(data *BootstrapInPlaceAgentConfigSpec) (string, error) {
	ictemplate, err := template.New("bip-agent-config").Parse(installConfigTemplate)
	if err != nil {
		return "", fmt.Errorf("unable to parse agent config template: %w", err)
	}

	contents := new(strings.Builder)

	if err := ictemplate.Execute(contents, data); err != nil {
		return "", fmt.Errorf("unable to generate agent config contents: %w", err)
	}

	return contents.String(), nil
}

var installConfigTemplate = `apiVersion: v1alpha1
kind: AgentConfig
metadata:
  name: {{ .VmName }}-sno-cluster
rendezvousIP: {{ .HostIP }}
hosts:
  - hostname: {{ .VmName }}
    interfaces:
      - name: eno1
        macAddress: {{ .HostMAC }}
    networkConfig:
      interfaces:
        - name: eno1
          type: ethernet
          state: up
          mac-address: {{ .HostMAC }}
          ipv4:
            enabled: true
            address:
              - ip: {{ .HostIP }}
                prefix-length: 24
      routes:
        config:
        - next-hop-address: {{ .HostRoute }}
          next-hop-interface: eno1
          destination: 0.0.0.0/0
      dns-resolver:
        config:
          server:
            - {{ .HostRoute }}`
