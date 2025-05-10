package handler

import (
	"encoding/json"
	"net/http"
)

var civs []string

func InitCivilisations(civList []string) {
	civs = civList
}

func GetCivilisationsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(civs)
}
