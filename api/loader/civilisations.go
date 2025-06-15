package loader

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"strings"
)

type Civ struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type identity struct {
	GenericName string `xml:"GenericName"`
}
type entity struct {
	Identity identity `xml:"Identity"`
}

func LoadCivilisations(unitsPath string) ([]Civ, error) {
	dirs, err := os.ReadDir(unitsPath)
	if err != nil {
		return nil, err
	}
	var civs []Civ
	for _, entry := range dirs {
		if entry.IsDir() {
			civCode := entry.Name()
			cataPath := filepath.Join(unitsPath, civCode, "catafalque.xml")
			file, err := os.Open(cataPath)
			if err != nil {
				continue // skip if file not found
			}
			var ent entity
			xml.NewDecoder(file).Decode(&ent)
			file.Close()
			name := strings.TrimSuffix(ent.Identity.GenericName, " Catafalque")
			civs = append(civs, Civ{Code: civCode, Name: name})
		}
	}
	return civs, nil
}
