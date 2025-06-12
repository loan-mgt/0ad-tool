package main

import (
	"0ad/tool/handler"
	"0ad/tool/loader"
	"0ad/tool/middleware"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
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

	r.GET("/civilisations/:civ_folder/units", middleware.VerifyCivFolderMiddleware, handler.GetUnitsHandler)
	r.GET("/civilisations/:civ_folder/units/:unit_code", middleware.VerifyCivFolderMiddleware, handler.GetUnitHandler)

	// Serve static files from the public directory at /static
	r.Static("/static", "./public")

	// Serve index.html at the root route using Go's template engine
	r.GET("/", func(c *gin.Context) {
		tmpl, err := template.New("index.html").Funcs(template.FuncMap{"lower": strings.ToLower}).ParseFiles("./public/index.html")
		if err != nil {
			c.String(http.StatusInternalServerError, "Template parsing error: %v", err)
			return
		}

		// Wrap Civs in a struct for extensibility
		data := struct {
			Civs []loader.Civ
		}{
			Civs: handler.Civs,
		}

		if err := tmpl.Execute(c.Writer, data); err != nil {
			c.String(http.StatusInternalServerError, "Template execution error: %v", err)
		}
	})

	log.Println("API server running on :8081")
	r.Run(":8081")
}
