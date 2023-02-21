package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/romanchechyotkin/car_booking-service/internal/config"
	user2 "github.com/romanchechyotkin/car_booking-service/internal/user"
	user "github.com/romanchechyotkin/car_booking-service/internal/user/storage"
	"github.com/romanchechyotkin/car_booking-service/pkg/client/postgresql"

	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	ctx := context.Background()

	log.Println("router init")
	router := httprouter.New()

	log.Println("config init")
	cfg := config.GetConfig()

	log.Println("postgresql config init")
	pgConfig := postgresql.NewPgConfig(
		cfg.PostgresStorage.Username,
		cfg.PostgresStorage.Password,
		cfg.PostgresStorage.Host,
		cfg.PostgresStorage.Port,
		cfg.PostgresStorage.Database,
	)
	pgClient := postgresql.NewClient(ctx, pgConfig)
	repository := user.NewRepository(pgClient)
	handler := user2.NewHandler(repository)
	handler.Register(router)

	log.Println("http server init")
	port := fmt.Sprintf(":%s", cfg.Listen.Port)
	server := http.Server{
		Handler:      router,
		Addr:         port,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	log.Printf("server running http://localhost:%s/", cfg.Listen.Port)
	log.Fatal(server.ListenAndServe())
}

func Print(str string) string {
	return fmt.Sprintf("string = %s", str)
}
