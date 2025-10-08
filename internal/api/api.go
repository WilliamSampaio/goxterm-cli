package api

import (
	"encoding/json"
	"goxterm-cli/internal/constants"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

type PingResult struct {
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

	var result = PingResult{}

	client := http.Client{
		Timeout: 2 * time.Second,
	}

	start := time.Now()

	resp, err := client.Get("https://www.google.com")
	if err != nil {
		result.Error = err.Error()
		log.Println("Ping failed:", err)
	} else {
		resp.Body.Close()
	}

	duration := time.Since(start)

	result.Alive = err == nil && resp.StatusCode == http.StatusOK
	result.Duration = float64(duration.Microseconds())

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func GetInfo(w http.ResponseWriter, r *http.Request) {
	headers(w)

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	type Shell struct {
		Bin     string `json:"bin"`
		Path    string `json:"path"`
		Default bool   `json:"default"`
	}

	type Info struct {
		AppName string  `json:"app_name"`
		Version string  `json:"version"`
		Shells  []Shell `json:"shells"`
	}

	info := Info{
		AppName: constants.AppName,
		Version: constants.AppVersion,
	}

	DefaultPath := os.Getenv("SHELL")

	parts := strings.Split(DefaultPath, "/")

	DefaultBin := parts[len(parts)-1]

	info.Shells = append(info.Shells, Shell{Bin: DefaultBin, Path: DefaultPath, Default: true})

	bins := []string{"bash", "zsh", "sh"}

	for _, bin := range bins {
		if bin == DefaultBin {
			continue
		}
		if path, err := exec.LookPath(bin); err == nil {
			info.Shells = append(info.Shells, Shell{Bin: bin, Path: path, Default: false})
		}
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(info); err != nil {
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
