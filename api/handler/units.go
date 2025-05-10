package handler

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type Unit struct {
	Code string `json:"code"`
	Name string `json:"name"`
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

	var units []Unit
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".xml") {
			continue
		}
		code := strings.TrimSuffix(file.Name(), ".xml")
		name := formatUnitName(code) // Assuming formatUnitName is defined elsewhere in this package
		units = append(units, Unit{Code: code, Name: name})
	}

	c.Header("Access-Control-Allow-Origin", "*")
	// c.JSON automatically sets Content-Type to application/json
	c.JSON(http.StatusOK, units)

}

// formatUnitName formats the unit code to a display name as specified
// This is a placeholder, assuming it's implemented correctly.
func formatUnitName(code string) string {
	// Example implementation: replace underscores with spaces and capitalize words
	name := strings.ReplaceAll(code, "_", " ")
	name = strings.ReplaceAll(name, "-", " ")
	return strings.Title(name) //nolint:staticcheck // SA1019: strings.Title is deprecated
}
