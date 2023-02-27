package main

import (
	"github.com/julienschmidt/httprouter"
	_ "github.com/romanchechyotkin/car_booking-service/docs"
	"github.com/romanchechyotkin/car_booking-service/internal/config"
	user2 "github.com/romanchechyotkin/car_booking-service/internal/user"
	user "github.com/romanchechyotkin/car_booking-service/internal/user/storage"
	"github.com/romanchechyotkin/car_booking-service/pkg/client/postgresql"
	"github.com/romanchechyotkin/car_booking-service/pkg/metrics"
	httpSwagger "github.com/swaggo/http-swagger"

	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

// @title           Car Booking Service API
// @version         1.0
// @host      localhost:5000
func main() {
	ctx := context.Background()

	log.Println("router init")
	router := httprouter.New()

	log.Println("swagger init")
	router.Handler(http.MethodGet, "/swagger", http.RedirectHandler("/swagger/index.html", http.StatusMovedPermanently))
	router.Handler(http.MethodGet, "/swagger/*any", httpSwagger.WrapHandler)

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
	storage := user2.NewService(repository)
	handler := user2.NewHandler(storage)
	handler.Register(router)

	go func() {
		log.Fatal(metrics.ListenMetrics("127.0.0.1:5001"))
	}()

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
