package configuration

import (
	"fmt"
	"os"
	"os/exec"
	"snoman/internal/logger"

	"go.uber.org/zap"
)

const (
	BIP_NETWORK_CONF_FILE = "/etc/NetworkManager/conf.d/bip.conf"
)

func ConfigureBootstrapInPlace() error {
	log := logger.Get()
	if !userIsRoot() {
		return fmt.Errorf("bootstrap in place requires elevated priveleges. Please run again as root")
	}

	if err := setDnsMasqConfig(log); err != nil {
		return fmt.Errorf("unable to configure dnsmasq: %w", err)
	}

	return nil
}

// setDnsMasqConfig requires sudo
func setDnsMasqConfig(log *zap.SugaredLogger) error {
	log.Infof("writing dnsmasq config for bootstrap in place to %s", BIP_NETWORK_CONF_FILE)
	os.WriteFile(BIP_NETWORK_CONF_FILE, []byte("[main]\ndns=dnsmasq"), 0644)

	log.Infof("reloading NetworkManager after new config was added")
	cmd := exec.Command("systemctl", "reload", "NetworkManager.service")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("unable to restart NetworkManager: %w", err)
	}

	return nil
}
