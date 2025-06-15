package middleware

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func VerifyCivFolderMiddleware(c *gin.Context) {
	civFolder := c.Param("civ_folder")
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
