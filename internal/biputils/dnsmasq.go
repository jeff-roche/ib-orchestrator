package biputils

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"snoman/internal/logger"

	"go.uber.org/zap"
)

const (
	BIP_NETWORK_CONF_FILE = "/etc/NetworkManager/conf.d/bip.conf"
)

func CreateDnsmasqConfig(address string) error {
	log := logger.Get()
	if !userIsRoot() {
		return fmt.Errorf("bootstrap in place requires elevated priveleges. Please run again as root")
	}

	if err := setDnsMasqConfig(address, log); err != nil {
		return fmt.Errorf("unable to configure dnsmasq: %w", err)
	}

	return nil
}

// setDnsMasqConfig requires sudo
func setDnsMasqConfig(address string, log *zap.SugaredLogger) error {
	log.Infof("writing dnsmasq config for bootstrap in place to %s", BIP_NETWORK_CONF_FILE)

	filedata := []byte(fmt.Sprintf("[main]\ndns=dnsmasq\naddress=%s", address))

	os.WriteFile(BIP_NETWORK_CONF_FILE, filedata, 0644)

	log.Infof("reloading NetworkManager after new config was added")
	cmd := exec.Command("systemctl", "reload", "NetworkManager.service")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("unable to restart NetworkManager: %w", err)
	}

	return nil
}

func userIsRoot() bool {
	user, _ := user.Current()

	if os.Getegid() == 0 && os.Geteuid() == 0 && user.Username == "root" {
		return true
	}

	return false
}
