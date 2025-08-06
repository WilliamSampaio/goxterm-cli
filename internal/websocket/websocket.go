package websocket

import (
	"fmt"
	"goxterm-cli/internal/config"
	"goxterm-cli/internal/sshclient"
	"goxterm-cli/internal/store"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SshWebSocketHandler(w http.ResponseWriter, r *http.Request) {

	strId := r.URL.Query().Get("id")

	id, err := strconv.Atoi(strId)
	if err != nil {
		log.Println("Error converting string to int:", err)
		return
	}

	log.Println("WebSocket connection request for:", id)

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer ws.Close()

	cfg, err := config.Load()
	if err != nil {
		log.Println("Error loading configuration:", err)
		http.Error(w, "Error loading configuration", http.StatusInternalServerError)
		return
	}

	if !store.Exists(cfg.StorePath) {
		log.Println("Store does not exist or is not located:", cfg.StorePath)
		http.Error(w, "Store does not exist or is not located", http.StatusNotFound)
		return
	}

	db, err := store.Load(cfg.StorePath)
	if err != nil {
		log.Println("Error loading store:", err)
		http.Error(w, "Error loading store", http.StatusInternalServerError)
		return
	}

	credential, err := db.GetSshSession(id)
	if err != nil {
		log.Printf("Connection '%d' not found in the store.\n", id)
		http.Error(w, fmt.Sprintf("Connection '%d' not found in the store", id), http.StatusNotFound)
		return
	}

	client, err := sshclient.ConnectSSH(*credential)
	if err != nil {
		log.Println("SSH dial error:", err)
		http.Error(w, "SSH connection error", http.StatusInternalServerError)
		return
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		log.Println("SSH session error:", err)
		http.Error(w, "Failed to create SSH session", http.StatusInternalServerError)
		return
	}
	defer session.Close()

	if err := sshclient.RequestTTY(session); err != nil {
		log.Println("Error request PTY:", err)
		return
	}

	stdinPipe, _ := session.StdinPipe()
	stdoutPipe, _ := session.StdoutPipe()
	stderrPipe, _ := session.StderrPipe()

	if err := session.Shell(); err != nil {
		log.Println("Failed to start shell:", err)
		return
	}

	go func() {
		io.Copy(&wsWriter{ws}, stdoutPipe)
	}()
	go func() {
		io.Copy(&wsWriter{ws}, stderrPipe)
	}()

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			break
		}
		stdinPipe.Write(msg)
	}
}

type wsWriter struct {
	ws *websocket.Conn
}

func (w *wsWriter) Write(p []byte) (int, error) {
	err := w.ws.WriteMessage(websocket.TextMessage, p)
	return len(p), err
}
