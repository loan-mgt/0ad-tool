package handler

import (
	"encoding/json"
	"net/http"

	"0ad/tool/loader"
)

var civs []loader.Civ

func InitCivilisations(civList []loader.Civ) {
	civs = civList
}

func GetCivilisationsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(civs)
}
