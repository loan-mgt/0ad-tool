package service

import (
	"0ad/tool/utils"
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
	println("Merging unit from:", filePath, "\n\n\n\n", "==========================")
	mergeElementToMap(unit, root)

	println("Merged unit:", utils.SPrintMapAny(unit))

	parentAttri := getRootAttr(root)

	if parentValue, ok := parentAttri["parent"]; ok {
		parentPath := resolveParentPath(filePath, parentValue)
		println("Parent path:", parentPath)
		if parentPath != "" {
			println(filePath, "Loading parent unit from:", parentPath)
			loadUnitRecursiveEtree(unit, parentPath)
		}
	}

}

// mergeElementToMap merges an etree.Element into a map[string]any recursively.
func mergeElementToMap(dst map[string]any, elem *etree.Element) {
	if dst == nil {
		dst = map[string]any{}
	}
	println(">==========================")
	println("Before merge:", utils.SPrintMapAny(dst))
	for _, child := range elem.ChildElements() {
		if len(child.ChildElements()) > 0 {
			var childMap map[string]any
			existing, exists := dst[child.Tag]
			if !exists {
				childMap = map[string]any{}
			} else {
				// Check if existing value is a map[string]any
				if m, ok := existing.(map[string]any); ok {
					childMap = m
				}
			}
			mergeElementToMap(childMap, child)
			dst[child.Tag] = childMap
		} else {
			//println("Child tag:", child.Tag, "Text:", child.Text(), "current:", utils.SPrintMapAny(dst))
			if _, exists := dst[child.Tag]; !exists || (exists && (dst[child.Tag] == nil || (func(v any) bool { _, ok := v.(map[string]any); return !ok }(dst[child.Tag])))) {
				dst[child.Tag] = child.Text()
			}
		}
	}
	println("After ", utils.SPrintMapAny(dst))
	println(">==========================")

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
