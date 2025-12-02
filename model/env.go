package model

import (
	"fmt"
	"os"
)

// MustEnv returns the value corresponding to the key passed as parameter.
// It returns an error if the key is missing from the .env file.
func MustEnv(key string) (string, error) {

	// Is the key in the file?
	if v := os.Getenv(key); v != "" {

		// If so, return the value read from .env file
		return v, nil
	}

	// If there's no key, return an empty string and an error
	return "", fmt.Errorf("missing env %s", key)
}
