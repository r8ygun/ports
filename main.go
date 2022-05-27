package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"ports/database"
	"ports/reader"
	"ports/uploader"
	"syscall"
)

const (
	DbHost     = "DB_HOST"
	DbUser     = "DB_USER"
	DbPassword = "DB_PASSWORD"
)

var (
	ErrDbHostNotSet     = errors.New("DB_HOST environment variable not set")
	ErrDbUserNotSet     = errors.New("DB_USER environment variable not set")
	ErrDbPasswordNotSet = errors.New("DB_PASSWORD environment variable not set")
)

func main() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-sigCh
		cancel()
	}()

	file, err := os.Open("./ports.json")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()
	defer cancel()

	logger := log.New(os.Stdout, "LOG: ", log.Ldate|log.Ltime)

	portReader := reader.Read(file)

	host, found := os.LookupEnv(DbHost)
	if !found {
		logger.Fatal(ErrDbHostNotSet)
	}

	user, found := os.LookupEnv(DbUser)
	if !found {
		logger.Fatal(ErrDbUserNotSet)
	}

	password, found := os.LookupEnv(DbPassword)
	if !found {
		logger.Fatal(ErrDbPasswordNotSet)
	}

	//I'll put the database/collection string as env vars later
	dbClient, err := database.New(ctx, host, user, password, "ports", "ports")
	defer dbClient.Stop(ctx)

	u := uploader.New(portReader, dbClient, logger)

	numUploaded := u.Upload(ctx)

	fmt.Println(numUploaded)
}
