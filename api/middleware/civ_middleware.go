package middleware

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// VerifyCivFolderMiddleware checks if the 'civ_folder' parameter corresponds to an existing directory.
// It assumes the API binary runs from a directory where '../templates/units/' is a valid relative path
// to the templates directory (e.g., running from 'api/' means templates are at '/workspaces/0ad-tool/templates/units/').
func VerifyCivFolderMiddleware(c *gin.Context) {
	civFolder := c.Param("civ_folder")
	// This path construction assumes the executable is run from a directory
	// such that "../templates/units/" correctly resolves.
	// For example, if the binary is in /workspaces/0ad-tool/api/,
	// this will point to /workspaces/0ad-tool/templates/units/.
	basePath := "../templates/units/"
	civDirPath := basePath + civFolder

	fileInfo, err := os.Stat(civDirPath)
	if os.IsNotExist(err) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Civilisation folder '" + civFolder + "' not found."})
		return
	}
	if err != nil {
		log.Printf("Error accessing civilisation folder path '%s': %v", civDirPath, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Server error while checking civilisation folder."})
		return
	}
	if !fileInfo.IsDir() {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Path '" + civFolder + "' is not a valid civilisation folder (not a directory)."})
		return
	}
	c.Next()
}
