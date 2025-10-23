package websocket

import (
	"goxterm-cli/internal/sshclient"
	"goxterm-cli/internal/store"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strconv"

	"github.com/creack/pty"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SshWebSocketHandler(w http.ResponseWriter, r *http.Request) {

	host := r.URL.Query().Get("host")

	port, err := strconv.Atoi(r.URL.Query().Get("port"))
	if err != nil {
		log.Println("Error get port:", err)
		return
	}

	user := r.URL.Query().Get("user")

	password := r.URL.Query().Get("password")

	credential := store.SshSession{
		Session: store.Session{
			Name: "",
		},
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
	}

	log.Printf("WebSocket connection request for: %s@%s:%d\n", user, host, port)

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer ws.Close()

	client, err := sshclient.ConnectSSH(credential)
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

	go func() {
		session.Wait()
		ws.Close()
	}()

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			break
		}
		stdinPipe.Write(msg)
	}
}

func ShellWebSocketHandler(w http.ResponseWriter, r *http.Request) {

	path := r.URL.Query().Get("path")

	log.Println("WebSocket connection request for:", path)

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer ws.Close()

	cmd := exec.Command(path)

	ptmx, err := pty.StartWithSize(cmd, &pty.Winsize{
		Rows: uint16(40),
		Cols: uint16(120),
	})
	if err != nil {
		log.Println("Erro ao iniciar PTY:", err)
		return
	}
	defer ptmx.Close()

	go func() {
		io.Copy(&wsWriter{ws}, ptmx)
	}()

	go func() {
		cmd.Wait()
		ws.Close()
	}()

	ptmx.Write([]byte("cd $HOME"))

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			break
		}
		ptmx.Write(msg)
	}
}

type wsWriter struct {
	ws *websocket.Conn
}

func (w *wsWriter) Write(p []byte) (int, error) {
	err := w.ws.WriteMessage(websocket.TextMessage, p)
	return len(p), err
}
