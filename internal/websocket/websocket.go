package websocket

import (
	"fmt"
	"goxterm-cli/internal/config"
	"goxterm-cli/internal/sshclient"
	"goxterm-cli/internal/store"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SshWebSocketHandler(w http.ResponseWriter, r *http.Request) {

	name := r.URL.Query().Get("name")
	log.Println("WebSocket connection request for:", name)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

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

	credential, exists := db.Credentials[name]
	if !exists {
		log.Printf("Connection '%s' not found in the store.\n", name)
		http.Error(w, fmt.Sprintf("Connection '%s' not found in the store", name), http.StatusNotFound)
		return
	}

	client, err := sshclient.ConnectSSH(credential)
	if err != nil {
		log.Println("SSH dial error:", err)
		http.Error(w, "SSH connection error", http.StatusInternalServerError)
		return
	}
	defer client.Close()

	sess, err := client.NewSession()
	if err != nil {
		log.Println("SSH session error:", err)
		http.Error(w, "Failed to create SSH session", http.StatusInternalServerError)
		return
	}
	defer sess.Close()

	stdin, _ := sess.StdinPipe()
	stdout, _ := sess.StdoutPipe()
	sess.Stderr = sess.Stdout

	if err := sshclient.RequestTTY(sess); err != nil {
		log.Println("SSH session error:", err)
		return
	}

	if err := sess.Shell(); err != nil {
		log.Println("failed to start shell:", err)
		return
	}

	// Leitor do SSH para WebSocket
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := stdout.Read(buf)
			if err != nil {
				break
			}
			conn.WriteMessage(websocket.TextMessage, buf[:n])
		}
	}()

	// Escritor do WebSocket para SSH
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		stdin.Write(msg)
	}
}
