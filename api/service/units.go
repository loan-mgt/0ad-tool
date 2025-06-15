package service

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/beevik/etree"
)

// GetUnit loads a unit from the given filePath, parsing the XML and resolving parents using etree.
func GetUnit(filePath string) map[string]any {
	println("================================ start load unit ================================")
	unit := map[string]any{"code": strings.TrimSuffix(filepath.Base(filePath), ".xml")}
	loadUnitRecursiveEtree(unit, filePath)
	println("Loading unit from:", filePath)
	return unit
}

// loadUnitRecursiveEtree loads a unit recursively, merging data from XML files using etree.
func loadUnitRecursiveEtree(unit map[string]any, filePath string) {

	doc := etree.NewDocument()
	err := doc.ReadFromFile(filePath)
	if err != nil {
		return
	}
	root := doc.Root()
	if root == nil {
		return
	}

	mergeElementToMap(unit, root)

	parentAttri := getRootAttr(root)

	if parentValue, ok := parentAttri["parent"]; ok {
		parentPath := resolveParentPath(filePath, parentValue)
		if parentPath != "" {
			loadUnitRecursiveEtree(unit, parentPath)
		}
	}

}

// mergeElementToMap merges an etree.Element into a map[string]any recursively.
func mergeElementToMap(dst map[string]any, elem *etree.Element) {
	if dst == nil {
		dst = map[string]any{}
	}

	for _, child := range elem.ChildElements() {
		if len(child.ChildElements()) > 0 {
			_, exists := dst[child.Tag]
			if !exists {
				dst[child.Tag] = map[string]any{}
			}
			if _, ok := dst[child.Tag].(map[string]any); ok {
				mergeElementToMap(dst[child.Tag].(map[string]any), child)
			}

		} else {
			//println("Child tag:", child.Tag, "Text:", child.Text(), "current:", utils.SPrintMapAny(dst))
			if _, exists := dst[child.Tag]; !exists {
				dst[child.Tag] = child.Text()
			}
		}
	}

}

func getRootAttr(elem *etree.Element) map[string]string {
	attrMap := map[string]string{}
	for _, attr := range elem.Attr {
		attrMap[attr.Key] = attr.Value
	}
	return attrMap
}

// resolveParentPath resolves the file path of a parent XML file, trying various locations.
func resolveParentPath(currentPath, parent string) string {
	// parent can be relative or absolute, or with | fallback
	parts := strings.Split(parent, "|")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if !strings.HasSuffix(p, ".xml") {
			p += ".xml"
		}
		// Try relative to templates/units, then templates/
		tryPaths := []string{
			filepath.Join("../templates/units", p),
			filepath.Join("../templates", p),
			filepath.Join(filepath.Dir(currentPath), p),
		}

		for _, try := range tryPaths {
			if _, err := os.Stat(try); err == nil {
				return try
			}
		}
	}
	return ""
}
