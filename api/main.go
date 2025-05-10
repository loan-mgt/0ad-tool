package main

import (
	"log"
	"net/http"

	"0ad/tool/handler"
	"0ad/tool/loader"
)

func main() {
	civs, err := loader.LoadCivilisations("../templates/units")
	if err != nil {
		log.Fatalf("Failed to read civilisations: %v", err)
	}
	handler.InitCivilisations(civs)

	http.HandleFunc("/civilisations", handler.GetCivilisationsHandler)
	log.Println("API server running on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
