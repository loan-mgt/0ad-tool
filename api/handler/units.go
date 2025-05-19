package handler

import (
	"0ad/tool/service"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

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

	var units []map[string]interface{}
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".xml") {
			continue
		}
		units = append(units, service.GetUnit(filepath.Join(civDir, file.Name())))
	}

	c.Header("Access-Control-Allow-Origin", "*")
	// c.JSON automatically sets Content-Type to application/json
	c.JSON(http.StatusOK, units)

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
