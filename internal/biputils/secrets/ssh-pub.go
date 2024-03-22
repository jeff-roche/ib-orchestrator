package secrets

import (
	"fmt"
	"os"
)

func GetSshPublicKey(fname string) (string, error) {
	var pKey string

	if fname == "" {
		return "", fmt.Errorf("no path to a public key was provided")
	}

	// Try to get the pull secret from the file
	data, err := os.ReadFile(fname)
	if err != nil {
		return "", fmt.Errorf("unable to read public key from %s: %w", fname, err)
	}

	pKey = string(data)
	if pKey == "" {
		return "", fmt.Errorf("the contents of %s were empty", fname)
	}

	return pKey, nil
}
