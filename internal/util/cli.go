package util

import (
	"os"
	"path/filepath"
)

func ResolveWorkingDir(wd string) (string, error) {
	if wd == "" {
		v, err := os.Getwd()
		if err != nil {
			return "", err
		}
		return v, nil
	}

	absPath, err := filepath.Abs(wd)
	if err != nil {
		return "", err
	}
	return absPath, nil
}
