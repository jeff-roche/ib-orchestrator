package secrets

import (
	"fmt"
	"os"
)

func GetPullSecret(fname string) (string, error) {
	var pullSecret string

	// See if a pull secret file was specified
	if fname != "" {
		// Try to get the pull secret from the file
		data, err := os.ReadFile(fname)
		if err != nil {
			return "", fmt.Errorf("unable to read pull secret from %s: %w", fname, err)
		}

		pullSecret = string(data)
		if pullSecret == "" {
			return "", fmt.Errorf("the contents of %s were empty", fname)
		}

		return pullSecret, nil
	}

	// Try the PULL_SECRET env var
	pullSecret = os.Getenv("PULL_SECRET")
	if pullSecret == "" {
		return "", fmt.Errorf("no pull secret specified")
	}

	return pullSecret, nil
}
