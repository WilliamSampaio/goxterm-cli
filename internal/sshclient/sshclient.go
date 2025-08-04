package sshclient

import (
	"fmt"
	"goxterm-cli/internal/store"
	"os"
	"os/signal"

	"github.com/charmbracelet/x/term"
	"golang.org/x/crypto/ssh"
)

// ConnectAndRun initializes an SSH connection, starts an interactive shell, and manages the session lifecycle.
func ConnectAndRun(credential store.SshSession) error {
	client, err := ConnectSSH(credential)
	if err != nil {
		return err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()

	oldState, err := setupTerminal()
	if err != nil {
		return err
	}
	defer term.Restore(os.Stdin.Fd(), oldState)

	handleInterrupt(oldState)

	if err := setupSessionIO(session); err != nil {
		return err
	}

	if err := RequestTTY(session); err != nil {
		return err
	}

	if err := session.Shell(); err != nil {
		return fmt.Errorf("failed to start shell: %v", err)
	}

	return session.Wait()
}

// connectSSH creates and returns an SSH client using the provided credentials.
func ConnectSSH(credential store.SshSession) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User: credential.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(credential.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", credential.Host, credential.Port), config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %v", err)
	}

	return client, nil
}

// setupTerminal puts the terminal into raw mode and returns the old state for restoration.
func setupTerminal() (*term.State, error) {
	oldState, err := term.MakeRaw(os.Stdin.Fd())
	if err != nil {
		return nil, fmt.Errorf("failed to enter raw mode: %v", err)
	}
	return oldState, nil
}

// handleInterrupt listens for Ctrl+C and restores the terminal before exiting.
func handleInterrupt(oldState *term.State) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		<-c
		term.Restore(os.Stdin.Fd(), oldState)
		os.Exit(0)
	}()
}

// setupSessionIO attaches the standard IO streams to the SSH session.
func setupSessionIO(session *ssh.Session) error {
	session.Stdin = os.Stdin
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	return nil
}

// requestTTY requests an interactive terminal for the SSH session.
func RequestTTY(session *ssh.Session) error {
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	fd := uintptr(os.Stdin.Fd())
	width, height, err := term.GetSize(fd)
	if err != nil {
		// fallback for 80x24
		width = 80
		height = 24
	}

	if err := session.RequestPty("xterm-256color", height, width, modes); err != nil {
		return fmt.Errorf("failed to request TTY: %v", err)
	}
	return nil
}
