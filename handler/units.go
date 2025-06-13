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
		if strings.HasSuffix(code, "_a") || strings.HasSuffix(code, "_b") || strings.HasSuffix(code, "_e") {
			code = code[:len(code)-2]
		}

		if unitInfo, ok := finalUnits[code]; ok {
			unitInfo.HasEvolution = true
			finalUnits[code] = unitInfo
			continue
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
			HasEvolution: false,
			Icon:         icon,
		}
	}

	c.Header("Access-Control-Allow-Origin", "*")
	// c.JSON automatically sets Content-Type to application/json
	c.JSON(http.StatusOK, finalUnits)

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
