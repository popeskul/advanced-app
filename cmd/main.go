package main

import (
	"advanced-app/internal/server"
	"fmt"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	mw := io.MultiWriter(os.Stdout, file)
	log.SetOutput(mw)

	log.SetFormatter(&log.JSONFormatter{})
}

// @title           Advanced App
// @version         1.0
// @description     This is a sample server celler server.
// @host      localhost:8080
// @BasePath  /
func main() {
	envMap, err := godotenv.Read("./.env")
	if err != nil {
		fmt.Println("Error loading .env file")
		os.Exit(1)
	}

	host, port := envMap["HOST"], envMap["PORT"]
	addr := fmt.Sprintf("%s:%s", host, port)

	server.Start(addr)

	fmt.Println("Server started on:", addr)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	fmt.Println("Server shutting down...")
}
