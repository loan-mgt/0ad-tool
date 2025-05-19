package service

import (
	"0ad/tool/utils"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/beevik/etree"
)

var (
	unitCache  = make(map[string]map[string]any)
	cacheMutex = &sync.Mutex{}
)

// GetUnit loads a unit from the given filePath, parsing the XML and resolving parents using etree.
func GetUnit(filePath string) map[string]any {
	unit := map[string]any{"code": strings.TrimSuffix(filepath.Base(filePath), ".xml")}
	loadUnitRecursiveEtree(unit, filePath)
	println("Loading unit from:", filePath)
	return unit
}

// loadUnitRecursiveEtree loads a unit recursively, merging data from XML files using etree.
func loadUnitRecursiveEtree(unit map[string]any, filePath string) {

	cacheMutex.Lock()
	if cachedUnit, exists := unitCache[filePath]; exists {
		println("Using cached unit for:", filePath)
		deepCopy(unit, cachedUnit)

		cacheMutex.Unlock()
		return
	}
	cacheMutex.Unlock()

	doc := etree.NewDocument()
	err := doc.ReadFromFile(filePath)
	if err != nil {
		return
	}
	root := doc.Root()
	if root == nil {
		return
	}
	println("Loaded XML root:", root.Tag)
	mergeElementToMap(unit, root)

	parentAttri := getRootAttr(root)

	if parentValue, ok := parentAttri["parent"]; ok {
		parentPath := resolveParentPath(filePath, parentValue)
		println("Parent path:", parentPath)
		if parentPath != "" {
			println(filePath, "Loading parent unit from:", parentPath)
			loadUnitRecursiveEtree(unit, parentPath)
		}
	}

	cacheMutex.Lock()
	if _, exists := unitCache[filePath]; !exists {
		unitCache[filePath] = unit
	}
	cacheMutex.Unlock()
	println("Unit loaded and cached:", filePath)
	println("Current unit state:", utils.SPrintMapAny(unit))

}

func deepCopy(dst map[string]any, src map[string]any) {
	for key, value := range src {
		switch v := value.(type) {
		case map[string]any:
			var childDst map[string]any
			if existing, exists := dst[key]; exists {
				if m, ok := existing.(map[string]any); ok {
					childDst = m
				} else {
					childDst = make(map[string]any)
				}
			} else {
				childDst = make(map[string]any)
			}
			deepCopy(childDst, v)
			dst[key] = childDst
		default:
			dst[key] = v
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
			if _, exists := dst[child.Tag]; !exists || (exists && (dst[child.Tag] == nil || (func(v any) bool { _, ok := v.(map[string]any); return !ok }(dst[child.Tag])))) {
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

		println("Trying paths:")
		for _, try := range tryPaths {
			println(try)
		}

		for _, try := range tryPaths {
			if _, err := os.Stat(try); err == nil {
				return try
			}
		}
	}
	return ""
}
