package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/romanchechyotkin/car_booking_service/docs"
	"github.com/romanchechyotkin/car_booking_service/internal/auth"
	emailproducer "github.com/romanchechyotkin/car_booking_service/internal/auth/producer"
	"github.com/romanchechyotkin/car_booking_service/internal/car"
	car2 "github.com/romanchechyotkin/car_booking_service/internal/car/storage/cars_storage"
	"github.com/romanchechyotkin/car_booking_service/internal/car/storage/images_storage"
	"github.com/romanchechyotkin/car_booking_service/internal/config"
	reservation "github.com/romanchechyotkin/car_booking_service/internal/reservation/storage"
	user2 "github.com/romanchechyotkin/car_booking_service/internal/user"
	user "github.com/romanchechyotkin/car_booking_service/internal/user/storage"
	"github.com/romanchechyotkin/car_booking_service/pkg/client/postgresql"
	grpc "github.com/romanchechyotkin/car_booking_service/pkg/grpc/client"
	"github.com/romanchechyotkin/car_booking_service/pkg/metrics"
	min "github.com/romanchechyotkin/car_booking_service/pkg/minio"
)

// TODO: IP feature (new device)
// TODO: new tables for postgresql for cars brands, models

// @title           Car Booking Service API
// @version         1.0
// @description  	P2P service for renting and booking cars
// @host      		localhost:5000
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	ctx := context.Background()

	client := min.New()
	log.Println(client)

	log.Println("gin init")
	router := gin.Default()
	router.Use(cors.Default())

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

	// TODO config for kafka
	kafkaConfig := &kafka.ConfigMap{
		"bootstrap.servers": cfg.Kafka.Port,
		"client.id":         "client",
		"acks":              "all",
	}
	producer, err := kafka.NewProducer(kafkaConfig)
	if err != nil {
		log.Fatalf("failed to connect to kafka %v", err)
	}
	defer producer.Close()

	emailPlacer := emailproducer.NewEmailPlacer(producer, cfg.Kafka.EmailTopic)
	authService := auth.NewService(repository, emailPlacer)
	authH := auth.NewHandler(authService)
	authH.Register(router)

	grpcClient := grpc.NewCarsManagementClient()
	carRepository := car2.NewRepository(pgClient)
	imgRep := images_storage.NewRepository(pgClient)
	reservationRep := reservation.NewRepository(pgClient)
	carHandler := car.NewHandler(carRepository, imgRep, reservationRep, repository, grpcClient, client)
	carHandler.Register(router)

	go func() {
		log.Fatal(metrics.ListenMetrics("127.0.0.1:5001"))
	}()

	router.GET("/health", health)

	log.Println("http server init")
	port := fmt.Sprintf(":%s", cfg.Listen.Port)
	server := http.Server{
		Handler:      router,
		Addr:         port,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	log.Printf("server running http://localhost:%s/health", cfg.Listen.Port)
	log.Println("docs http://localhost:5000/swagger/index.html")
	log.Fatal(server.ListenAndServe())
}

// @Summary Health Check
// @Description Checking health of backend
// @Produce application/json
// @Success 200 {string} health
// @Router /health [get]
func health(ctx *gin.Context) {
	ctx.String(http.StatusOK, "health")
}
