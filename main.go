package main

//go:generate go run github.com/99designs/gqlgen

import (
	"io"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/ngutzmann/wireguard-web-config/server"
)

func main() {
	godotenv.Load()

	if mode := os.Getenv("GIN_MODE"); mode == server.RELEASE {
		file := os.Getenv("LOG_FILE")
		if len(file) != 0 {
			f, err := os.OpenFile("/tmp/orders.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				log.Fatalf("error opening file: %v", err)
			}
			defer f.Close()
			wrt := io.MultiWriter(os.Stdout, f)
			log.SetOutput(wrt)
		}
	}

	server.Server()
}
