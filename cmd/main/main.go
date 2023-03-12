package main

import (
	"github.com/gin-gonic/gin"
	"github.com/romanchechyotkin/car_booking_service/internal/auth"
	"github.com/romanchechyotkin/car_booking_service/internal/config"
	"github.com/romanchechyotkin/car_booking_service/pkg/client/postgresql"
	"github.com/romanchechyotkin/car_booking_service/pkg/metrics"

	_ "github.com/romanchechyotkin/car_booking_service/docs"
	user2 "github.com/romanchechyotkin/car_booking_service/internal/user"
	user "github.com/romanchechyotkin/car_booking_service/internal/user/storage"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

// TODO: http tests, swagger docs, add secret to config file, refactor code, add Makefile

// @title           Car Booking Service API
// @version         1.0
// @host      localhost:5000
func main() {
	ctx := context.Background()

	log.Println("gin init")
	router := gin.Default()

	log.Println("swagger init")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

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

	authService := auth.NewService(repository)
	authH := auth.NewHandler(authService)
	authH.Register(router)

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
