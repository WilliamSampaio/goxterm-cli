package api

import (
	"encoding/json"
	"goxterm-cli/internal/config"
	"goxterm-cli/internal/constants"
	"goxterm-cli/internal/store"
	"log"
	"net"
	"net/http"
	"os/exec"
	"time"
)

type PingResult struct {
	IP       string  `json:"ip"`
	Alive    bool    `json:"alive"`
	Duration float64 `json:"duration_ms"`
	Error    string  `json:"error,omitempty"`
}

func Ping(w http.ResponseWriter, r *http.Request) {
	headers(w)

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ip := "8.8.8.8"

	start := time.Now()

	conn, err := net.DialTimeout("tcp", net.JoinHostPort(ip, "53"), 2*time.Second)

	duration := time.Since(start)

	result := PingResult{
		IP:       ip,
		Alive:    err == nil,
		Duration: float64(duration.Microseconds()),
	}

	if err != nil {
		log.Println("Ping failed:", err)
	} else {
		conn.Close()
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func GetInfo(w http.ResponseWriter, r *http.Request) {
	headers(w)

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	type Info struct {
		AppName string   `json:"app_name"`
		Version string   `json:"version"`
		Shells  []string `json:"shells"`
	}

	info := Info{
		AppName: constants.AppName,
		Version: constants.AppVersion,
	}

	if _, err := exec.LookPath("bash"); err == nil {
		info.Shells = append(info.Shells, "bash")
	}

	if _, err := exec.LookPath("zsh"); err == nil {
		info.Shells = append(info.Shells, "zsh")
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(info); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func GetListCredentials(w http.ResponseWriter, r *http.Request) {
	headers(w)

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

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(db.SshSessions); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func headers(w http.ResponseWriter) {
	//Allow CORS here By * or specific origin
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	w.Header().Set("Content-Type", "application/json")
}
