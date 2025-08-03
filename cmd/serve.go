/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"goxterm-cli/internal/api"
	"goxterm-cli/internal/websocket"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the GoXterm web server",
	Run: func(cmd *cobra.Command, args []string) {
		serve(port)
	},
}

func init() {
	defaultPort := 8080
	serveCmd.Flags().IntVarP(&port, "port", "p", defaultPort, fmt.Sprintf("Port to run the server on (default is %d)", defaultPort))
	rootCmd.AddCommand(serveCmd)
}

func serve(port int) {
	fmt.Println("Starting GoXterm web server...")

	http.HandleFunc("/api/credentials", api.GetListCredentials)

	http.HandleFunc("/ssh", websocket.SshWebSocketHandler)

	http.Handle("/", http.FileServer(http.Dir("/usr/share/goxterm")))

	log.Printf("Server started at http://localhost:%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
