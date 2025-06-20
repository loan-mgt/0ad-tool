package handler

import (
	"net/http"

	"0ad/tool/loader"

	"github.com/gin-gonic/gin" // Added for gin.Context
)

var Civs []loader.Civ

func InitCivilisations(civList []loader.Civ) {
	Civs = civList
}

// GetCivilisationsHandler now uses gin.Context
func GetCivilisationsHandler(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	// Content-Type is automatically set to application/json by c.JSON
	c.JSON(http.StatusOK, Civs)
}
