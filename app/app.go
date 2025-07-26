package app

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/charmbracelet/x/term"
	"github.com/urfave/cli"
	"golang.org/x/crypto/ssh"
)

type Credential struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
}

func Application() *cli.App {
	app := cli.NewApp()
	app.Name = "CLI App"
	app.Usage = "A simple CLI application"
	app.Version = "1.0.0"

	flags := []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "Nome da conexão",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:   "open",
			Usage:  "Abre uma conexão com o servidor",
			Flags:  flags,
			Action: open,
		},
	}

	return app
}

func open(cli *cli.Context) {
	name := cli.String("name")

	var credentials []Credential

	file, err := os.Open("credentials.json")
	if err != nil {
		fmt.Println("Erro ao abrir o arquivo:", err)
		return
	}
	defer file.Close()

	bytes, err := os.ReadFile("credentials.json")
	if err != nil {
		fmt.Println("Erro ao ler o arquivo:", err)
		return
	}

	err = json.Unmarshal(bytes, &credentials)
	if err != nil {
		fmt.Println("Erro ao decodificar o JSON:", err)
		return
	}

	for _, credential := range credentials {
		// fmt.Println(i, credential.Name, credential.Host, credential.Port, credential.User)

		fmt.Println(name)

		if credential.Name == name {
			fmt.Printf("Conectando a %s (%s:%d) como %s...\n", credential.Name, credential.Host, credential.Port, credential.User)
			// Aqui você pode adicionar a lógica para abrir a conexão com o servidor
			// Por exemplo, usando uma biblioteca SSH ou outra de sua escolha

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
				log.Fatalf("Falha ao conectar: %v", err)
			}
			defer client.Close()

			// Cria nova sessão SSH
			session, err := client.NewSession()
			if err != nil {
				log.Fatalf("Falha ao criar sessão: %v", err)
			}
			defer session.Close()

			// Prepara o terminal local para modo raw (interativo)
			oldState, err := term.MakeRaw(os.Stdin.Fd())
			if err != nil {
				log.Fatalf("Falha ao entrar em modo raw: %v", err)
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
				log.Fatalf("Falha ao solicitar TTY: %v", err)
			}

			// Inicia shell interativo
			if err := session.Shell(); err != nil {
				log.Fatalf("Falha ao iniciar shell: %v", err)
			}

			// Espera até a sessão terminar
			session.Wait()
		}
	}
}
