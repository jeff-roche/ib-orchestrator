package configuration

import (
	"os"
	"os/user"
)

func userIsRoot() bool {
	user, _ := user.Current()

	if os.Getegid() == 0 && os.Geteuid() == 0 && user.Username == "root" {
		return true
	}

	return false
}
