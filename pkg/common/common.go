package common

import (
	"os"
	"os/exec"
)

func CommandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func DotDir(dotType string) (string, error) {
	switch dotType {
	case "config":
		dir, err := os.UserConfigDir()
		if err != nil {
			return "", err
		}
		return dir, nil
	case "home":
		dir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return dir, nil
	}
	return "", nil
}
