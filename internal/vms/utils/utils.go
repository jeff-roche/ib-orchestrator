package utils

import (
	"snoman/internal/logger"

	"gopkg.in/yaml.v2"
	"libvirt.org/go/libvirt"
)

func GetLibvirtConnection() (*libvirt.Connect, error) {
	return libvirt.NewConnect("qemu:///system") // TODO: make sure this is the right connection string (make configurable?)
}

func LogYaml(data interface{}) {
	log := logger.Get()

	yamlData, err := yaml.Marshal(data)
	if err != nil {
		log.Fatalf("unable to log yaml data: %w", err)
	}

	log.Infof("\n%s", string(yamlData))
}
