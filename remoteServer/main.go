package main

import (
	"context"
	"fmt"
	"log"
	"music_service/config"
	"music_service/swaggerAPI/database"
	srv "music_service/swaggerAPI/go"
	"net/http"
	"os"
	"os/signal"

	"sync"
	"syscall"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	cfg := config.NewRemoteServerConfig()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)
	go database.Broadcast(ctx, &wg)

	log := srv.SetupLogger()

	router := srv.NewRouter()
	log.Info(fmt.Sprintf("server listening at %v", cfg.Address))

	go func() {
		log.Debug(http.ListenAndServe(cfg.Address, router).Error())
	}()

	wg.Wait()
	log.Info("server stopped")
}
