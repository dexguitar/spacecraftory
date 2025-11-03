package path

import (
	"os"
	"path/filepath"
)

// GetProjectRoot searches for the project root directory by checking for the go.work file
func GetProjectRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		panic("failed to get working directory: " + err.Error())
	}

	for {
		_, err = os.Stat(filepath.Join(dir, "go.work"))
		if err == nil {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			panic("failed to find project root (go.work)")
		}

		dir = parent
	}
}
