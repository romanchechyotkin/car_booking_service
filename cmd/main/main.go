package main

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gin-gonic/gin"
	_ "github.com/romanchechyotkin/car_booking_service/docs"
	"github.com/romanchechyotkin/car_booking_service/internal/auth"
	"github.com/romanchechyotkin/car_booking_service/internal/auth/producer"
	"github.com/romanchechyotkin/car_booking_service/internal/car"
	paymentproducer "github.com/romanchechyotkin/car_booking_service/internal/car/producer"
	car2 "github.com/romanchechyotkin/car_booking_service/internal/car/storage/cars_storage"
	"github.com/romanchechyotkin/car_booking_service/internal/car/storage/images_storage"
	"github.com/romanchechyotkin/car_booking_service/internal/config"
	reservation "github.com/romanchechyotkin/car_booking_service/internal/reservation/storage"
	user2 "github.com/romanchechyotkin/car_booking_service/internal/user"
	user "github.com/romanchechyotkin/car_booking_service/internal/user/storage"
	"github.com/romanchechyotkin/car_booking_service/pkg/client/postgresql"
	"github.com/romanchechyotkin/car_booking_service/pkg/metrics"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"net/http"
	"time"
)

// TODO: http tests, swagger docs, add secret to config file, refactor code
// TODO: IP feature (new device)
// TODO: credit-cards microservice & payments microservice

// @title           Car Booking Service API
// @version         1.0
// @host      localhost:5000
func main() {
	ctx := context.Background()

	//conn, err := grpc.Dial("localhost:9000", grpc.WithInsecure())
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//defer conn.Close()

	//client := pb.NewImageSenderClient(conn)

	log.Println("gin init")
	router := gin.Default()
	router.Static("/static", "./static")

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

	carRepository := car2.NewRepository(pgClient)
	imgRep := images_storage.NewRepository(pgClient)
	reservationRep := reservation.NewRepository(pgClient)
	paymentPlacer := paymentproducer.NewPaymentPlacer(producer, cfg.Kafka.PaymentTopic)
	carHandler := car.NewHandler(carRepository, imgRep, paymentPlacer, reservationRep, repository)
	carHandler.Register(router)

	go func() {
		log.Fatal(metrics.ListenMetrics("127.0.0.1:5001"))
	}()

	router.GET("/health", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "health")
	})

	log.Println("http server init")
	port := fmt.Sprintf(":%s", cfg.Listen.Port)
	server := http.Server{
		Handler:      router,
		Addr:         port,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	log.Printf("server running http://localhost:%s/health", cfg.Listen.Port)
	log.Fatal(server.ListenAndServe())
}
