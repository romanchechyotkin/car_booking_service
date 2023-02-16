package main

import (
	"context"
	"github.com/romanchechyotkin/car_booking-service/internal/config"
	"github.com/romanchechyotkin/car_booking-service/pkg/client/postgresql"

	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	log.Println("config init")
	cfg := config.GetConfig()

	log.Println("postgresql config init")
	pgConfig := postgresql.NewPgConfig(
		cfg.PostgresQLStorage.Username,
		cfg.PostgresQLStorage.Password,
		cfg.PostgresQLStorage.Host,
		cfg.PostgresQLStorage.Port,
		cfg.PostgresQLStorage.Database,
	)
	_ = postgresql.NewClient(context.Background(), pgConfig)

	log.Println("http server init")
	port := fmt.Sprintf(":%s", cfg.Listen.Port)
	server := http.Server{
		Addr:         port,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	log.Printf("server running http://localhost:%s/", cfg.Listen.Port)
	log.Fatal(server.ListenAndServe())
}
