package main

import (
	"goxterm/app"
	"log"
	"os"
)

func main() {
	application := app.Application()
	if err := application.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
