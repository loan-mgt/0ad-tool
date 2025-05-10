package main

import (
	"0ad/tool/handler"
	"0ad/tool/loader"
	"0ad/tool/middleware"
	"log"

	"github.com/gin-contrib/cors" // Added import for CORS middleware
	"github.com/gin-gonic/gin"
)

func main() {
	civs, err := loader.LoadCivilisations("../templates/units")
	if err != nil {
		log.Fatalf("Failed to read civilisations: %v", err)
	}
	handler.InitCivilisations(civs)

	r := gin.Default()

	// Add CORS middleware. Default() allows all origins
	r.Use(cors.Default())

	r.GET("/civilisations", handler.GetCivilisationsHandler)

	// Updated route to use the middleware from the new package and correct handler reference
	r.GET("/civilisations/:civ_folder/units", middleware.VerifyCivFolderMiddleware, handler.GetUnitsHandler)

	log.Println("API server running on :8081")
	r.Run(":8081")
}
