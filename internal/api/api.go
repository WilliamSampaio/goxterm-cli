package api

import (
	"encoding/json"
	"goxterm-cli/internal/config"
	"goxterm-cli/internal/store"
	"net/http"
)

func GetListCredentials(w http.ResponseWriter, r *http.Request) {
	//Allow CORS here By * or specific origin
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cfg, err := config.Load()
	if err != nil {
		http.Error(w, "Error loading configuration", http.StatusInternalServerError)
		return
	}

	if !store.Exists(cfg.StorePath) {
		http.Error(w, "Store does not exist or is not located", http.StatusNotFound)
		return
	}

	db, err := store.Load(cfg.StorePath)
	if err != nil {
		http.Error(w, "Error loading store", http.StatusInternalServerError)
		return
	}

	var list []string

	for k := range db.Credentials {
		list = append(list, k)
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(list); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}
