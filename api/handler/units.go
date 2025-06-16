package handler

import (
	"0ad/tool/service"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type UnitShortInfo struct {
	Code         string `json:"code"`
	Name         string `json:"name"`
	Class        string `json:"class"`
	HasEvolution bool   `json:"has_evolution"`
	Icon         string `json:"icon"`
}

type UnitCategory struct {
	Infantry []UnitShortInfo `json:"infantry"`
	Cavalry  []UnitShortInfo `json:"cavalry"`
	Champion []UnitShortInfo `json:"champion"`
	Hero     []UnitShortInfo `json:"hero"`
	Ship     []UnitShortInfo `json:"ship"`
	Siege    []UnitShortInfo `json:"siege"`
	Support  []UnitShortInfo `json:"support"`
	Other    []UnitShortInfo `json:"other"`
}

func GetUnitsHandler(c *gin.Context) {
	unitsPath := "../templates/units"
	civFolder := c.Param("civ_folder") // Use c.Param to get the path parameter
	civDir := filepath.Join(unitsPath, civFolder)

	files, err := os.ReadDir(civDir)
	if err != nil {
		// Log the error for server-side diagnostics
		// log.Printf("Error reading directory %s: %v", civDir, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read civilisation units directory"})
		return
	}

	var units []map[string]any
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".xml") {
			continue
		}
		units = append(units, service.GetUnit(filepath.Join(civDir, file.Name())))
	}

	finalUnits := make(map[string]UnitShortInfo)

	for _, unit := range units {
		identity, ok := unit["Identity"].(map[string]any)
		if !ok {
			continue
		}

		specificName, ok := identity["SpecificName"].(string)
		if !ok {
			continue
		}

		code, ok := unit["code"].(string)
		if !ok {
			continue
		}

		// not handling these case
		if strings.HasSuffix(code, "_trireme") {
			code = strings.TrimSuffix(code, "_trireme")
		}

		if strings.HasSuffix(code, "_a") || strings.HasSuffix(code, "_e") || strings.HasSuffix(code, "_packed") || strings.HasSuffix(code, "_house") || strings.HasSuffix(code, "_pike") || strings.HasSuffix(code, "_infantry") || strings.HasSuffix(code, "_fire_fire") {
			continue
		}

		// skip champion variants
		if code == "hero_boudicca_cavalry_javelineer" || code == "hero_boudicca_sword" || code == "hero_wei_qing_horse" || code == "hero_wei_qing_chariot" || code == "hero_liu_bang_horse" || code == "hero_han_xin_horse" || code == "hero_xerxes_i_chariot" {
			continue
		}

		hasEvolution := false
		if strings.HasSuffix(code, "_b") {
			code = strings.TrimSuffix(code, "_b")
			hasEvolution = true
		}

		if strings.HasSuffix(code, "_unpacked") {
			code = strings.TrimSuffix(code, "_unpacked")
			hasEvolution = true
		}

		icon, ok := identity["Icon"].(string)
		if !ok {
			icon = "" // Default to empty string if Icon is not found
		}

		class, ok := identity["VisibleClasses"].(string)
		if !ok {
			continue
		}

		finalUnits[code] = UnitShortInfo{
			Code:         code,
			Name:         specificName,
			Class:        class,
			HasEvolution: hasEvolution,
			Icon:         icon,
		}
	}

	res := UnitCategory{
		Infantry: make([]UnitShortInfo, 0),
		Cavalry:  make([]UnitShortInfo, 0),
		Champion: make([]UnitShortInfo, 0),
		Hero:     make([]UnitShortInfo, 0),
		Ship:     make([]UnitShortInfo, 0),
		Siege:    make([]UnitShortInfo, 0),
		Support:  make([]UnitShortInfo, 0),
		Other:    make([]UnitShortInfo, 0),
	}

	for code, unitInfo := range finalUnits {
		unitCode := strings.ToLower(code)
		switch {
		case strings.Contains(unitCode, "hero"):
			res.Hero = append(res.Hero, unitInfo)
			break
		case strings.Contains(unitCode, "champion"):
			res.Champion = append(res.Champion, unitInfo)
			break
		case strings.Contains(unitCode, "infantry"):
			res.Infantry = append(res.Infantry, unitInfo)
			break
		case strings.Contains(unitCode, "skirmisher"):
			res.Infantry = append(res.Infantry, unitInfo)
			break
		case strings.Contains(unitCode, "hoplite"):
			res.Infantry = append(res.Infantry, unitInfo)
			break
		case strings.Contains(unitCode, "arstibara"):
			res.Infantry = append(res.Infantry, unitInfo)
			break
		case strings.Contains(unitCode, "cavalry"):
			res.Cavalry = append(res.Cavalry, unitInfo)
			break
		case strings.HasPrefix(unitCode, "ship"):
			res.Ship = append(res.Ship, unitInfo)
			break
		case strings.HasPrefix(unitCode, "siege"):
			res.Siege = append(res.Siege, unitInfo)
			break
		case strings.HasPrefix(unitCode, "support"):
			res.Support = append(res.Support, unitInfo)
			break
		case strings.HasPrefix(unitCode, "elephant"):
			res.Infantry = append(res.Infantry, unitInfo)
			break
		case strings.HasPrefix(unitCode, "war"):
			res.Infantry = append(res.Infantry, unitInfo)
			break
		default:
			res.Other = append(res.Other, unitInfo)
		}
	}

	c.Header("Access-Control-Allow-Origin", "*")
	// c.JSON automatically sets Content-Type to application/json
	c.JSON(http.StatusOK, res)

}

func GetUnitHandler(c *gin.Context) {
	civFolder := c.Param("civ_folder")
	unitCode := c.Param("unit_code")
	unitsPath := "../templates/units"
	civDir := filepath.Join(unitsPath, civFolder)
	unitFilePath := filepath.Join(civDir, unitCode+".xml")

	// Try the main file, then _b, then _packed
	tryFiles := []string{
		unitFilePath,
		filepath.Join(civDir, unitCode+"_b.xml"),
		filepath.Join(civDir, unitCode+"_unpacked.xml"),
	}

	var foundFile string
	for _, f := range tryFiles {
		if _, err := os.Stat(f); err == nil {
			foundFile = f
			break
		}
	}

	if foundFile == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Unit not found"})
		return
	}

	units := []string{}
	switch {
	case strings.HasSuffix(foundFile, "_b.xml"):
		units = append(units, foundFile)
		units = append(units, strings.TrimSuffix(foundFile, "_b.xml")+"_a.xml")
		units = append(units, strings.TrimSuffix(foundFile, "_b.xml")+"_e.xml")
	case strings.HasSuffix(foundFile, "_unpacked.xml"):
		units = append(units, foundFile)
		units = append(units, strings.TrimSuffix(foundFile, "_unpacked.xml")+"_packed.xml")
	default:
		units = append(units, foundFile)
	}

	res := make([]map[string]any, 0)
	for _, unit := range units {
		res = append(res, service.GetUnit(unit))
	}

	keysToKeep := []string{
		"Identity", "Health", "Cost", "Attack",
		"Resistance", "ResourceGatherer", "Vision", "UnitMotion", "Pack",
	}

	cleanedUnit := make([]map[string]any, len(res))
	for i, unit := range res {
		cleanedUnit[i] = make(map[string]any) // Initialize each map in the slice
		for _, key := range keysToKeep {
			if value, ok := unit[key].(map[string]any); ok {
				cleanedUnit[i][key] = value
			}
		}
	}

	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, cleanedUnit)
}
