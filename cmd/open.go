/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"goxterm-cli/internal/config"
	"goxterm-cli/internal/store"
	"os"
	"os/signal"
	"strings"

	"github.com/charmbracelet/x/term"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

// openCmd represents the open command
var openCmd = &cobra.Command{
	Use:   "open",
	Short: "Open a SSH connection",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		if name == "" {
			fmt.Print("Name: ")
			nameFromUser, _ := reader.ReadString('\n')
			name = strings.TrimSpace(nameFromUser)
		}

		cfg, err := config.Load()
		if err != nil {
			fmt.Printf("Error loading configuration: %v\n", err)
			os.Exit(1)
		}

		if !store.Exists(cfg.StorePath) {
			fmt.Printf("Store does not exist or is not located: %v\n", err)
			os.Exit(1)
		}

		db, err := store.Load(cfg.StorePath)
		if err != nil {
			fmt.Printf("Error loading store: %v\n", err)
			os.Exit(1)
		}

		credential, exists := db.Credentials[name]
		if !exists {
			fmt.Printf("Connection '%s' not found in the store.\n", name)
			fmt.Println("Please run 'goxterm save' to add a new connection.")
			os.Exit(1)
		}

		// Configurações do cliente SSH
		config := &ssh.ClientConfig{
			User: credential.User, // substitua pelo seu usuário
			Auth: []ssh.AuthMethod{
				ssh.Password(credential.Password), // ou use chave privada
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(), // cuidado: desativa verificação da chave do host
		}

		// Conecta ao servidor SSH
		client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", credential.Host, credential.Port), config)
		if err != nil {
			fmt.Printf("Falha ao conectar: %v", err)
			os.Exit(1)
		}
		defer client.Close()

		// Cria nova sessão SSH
		session, err := client.NewSession()
		if err != nil {
			fmt.Printf("Falha ao criar sessão: %v", err)
			os.Exit(1)
		}
		defer session.Close()

		// Prepara o terminal local para modo raw (interativo)
		oldState, err := term.MakeRaw(os.Stdin.Fd())
		if err != nil {
			fmt.Printf("Falha ao entrar em modo raw: %v", err)
			os.Exit(1)
		}
		defer term.Restore(os.Stdin.Fd(), oldState)

		// Captura Ctrl+C para restaurar terminal
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		go func() {
			<-c
			term.Restore(os.Stdin.Fd(), oldState)
			os.Exit(0)
		}()

		// Redireciona IO (entrada/saída padrão)
		session.Stdin = os.Stdin
		session.Stdout = os.Stdout
		session.Stderr = os.Stderr

		// Solicita TTY (modo interativo)
		modes := ssh.TerminalModes{
			ssh.ECHO:          1,
			ssh.TTY_OP_ISPEED: 14400,
			ssh.TTY_OP_OSPEED: 14400,
		}

		if err := session.RequestPty("xterm-256color", 80, 40, modes); err != nil {
			fmt.Printf("Falha ao solicitar TTY: %v", err)
			os.Exit(1)
		}

		// Inicia shell interativo
		if err := session.Shell(); err != nil {
			fmt.Printf("Falha ao iniciar shell: %v", err)
			os.Exit(1)
		}

		// Espera até a sessão terminar
		session.Wait()
	},
}

func init() {
	openCmd.Flags().StringVarP(&name, "name", "n", "", "Enter the name of the connection")
	rootCmd.AddCommand(openCmd)
}
