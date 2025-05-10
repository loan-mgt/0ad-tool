package loader

import (
	"os"
)

func LoadCivilisations(unitsPath string) ([]string, error) {
	dirs, err := os.ReadDir(unitsPath)
	if err != nil {
		return nil, err
	}
	var civs []string
	for _, entry := range dirs {
		if entry.IsDir() {
			civs = append(civs, entry.Name())
		}
	}
	return civs, nil
}
