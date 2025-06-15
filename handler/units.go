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

		if strings.HasSuffix(code, "_a") || strings.HasSuffix(code, "_e") || strings.HasSuffix(code, "_packed") || strings.HasSuffix(code, "_house") {
			continue
		}

		hasEvolution := false
		if strings.HasSuffix(code, "_unpacked") || strings.HasSuffix(code, "_b") {
			code = code[:len(code)-2]
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
		case strings.Contains(unitCode, "infantry"):
			res.Infantry = append(res.Infantry, unitInfo)
		case strings.Contains(unitCode, "cavalry"):
			res.Cavalry = append(res.Cavalry, unitInfo)
		case strings.Contains(unitCode, "champion"):
			res.Champion = append(res.Champion, unitInfo)
		case strings.Contains(unitCode, "hero"):
			res.Hero = append(res.Hero, unitInfo)
		case strings.HasPrefix(unitCode, "ship"):
			res.Ship = append(res.Ship, unitInfo)
		case strings.HasPrefix(unitCode, "siege"):
			res.Siege = append(res.Siege, unitInfo)
		case strings.HasPrefix(unitCode, "support"):
			res.Support = append(res.Support, unitInfo)
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

	if _, err := os.Stat(unitFilePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Unit not found"})
		return
	}

	unit := service.GetUnit(unitFilePath)
	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, unit)
}
