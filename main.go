package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/romanchechyotkin/car_booking_service/docs"
	"github.com/romanchechyotkin/car_booking_service/internal/auth"
	"github.com/romanchechyotkin/car_booking_service/internal/car"
	car2 "github.com/romanchechyotkin/car_booking_service/internal/car/storage/cars_storage"
	"github.com/romanchechyotkin/car_booking_service/internal/car/storage/images_storage"
	reservation "github.com/romanchechyotkin/car_booking_service/internal/reservation/storage"
	user2 "github.com/romanchechyotkin/car_booking_service/internal/user"
	user "github.com/romanchechyotkin/car_booking_service/internal/user/storage"
	"github.com/romanchechyotkin/car_booking_service/pkg/client/postgresql"
	"github.com/romanchechyotkin/car_booking_service/pkg/config"

	// grpc "github.com/romanchechyotkin/car_booking_service/pkg/grpc/client"
	"github.com/romanchechyotkin/car_booking_service/pkg/metrics"
	min "github.com/romanchechyotkin/car_booking_service/pkg/minio"
)

func main() {
	ctx := context.Background()

	log.Println("gin init")
	router := gin.Default()
	router.Use(CORSMiddleware())

	log.Println("swagger init")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	log.Println("config init")
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("minio init")
	client := min.New(cfg)
	log.Println(client)

	log.Println("postgresql config init")
	pgConfig := postgresql.NewPgConfig(
		cfg.Postgresql.User,
		cfg.Postgresql.Password,
		cfg.Postgresql.Host,
		cfg.Postgresql.Port,
		cfg.Postgresql.Database,
	)
	pgClient := postgresql.NewClient(ctx, pgConfig)
	
	dbURL := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", cfg.Postgresql.User, cfg.Postgresql.Password, cfg.Postgresql.Host, cfg.Postgresql.Port, cfg.Postgresql.Database)
	postgresql.Migrate(dbURL)
	
	repository := user.NewRepository(pgClient)
	handler := user2.NewHandler(repository, client)
	handler.Register(router)

	//producer, err := kafka.NewProducer(kafkaConfig)
	//if err != nil {
	//	log.Fatalf("failed to connect to kafka %v", err)
	//}
	//defer producer.Close()
	//emailPlacer := emailproducer.NewEmailPlacer(producer, cfg.Kafka.EmailTopic)

	authService := auth.NewService(repository)
	authH := auth.NewHandler(authService)
	authH.Register(router)

	// grpcClient := grpc.NewCarsManagementClient(cfg.ElasticSearchMicroservice.Host, cfg.ElasticSearchMicroservice.Port)
	carRepository := car2.NewRepository(pgClient)
	imgRep := images_storage.NewRepository(pgClient)
	reservationRep := reservation.NewRepository(pgClient)
	carHandler := car.NewHandler(carRepository, imgRep, reservationRep, repository, client)
	carHandler.Register(router)

	go func() {
		log.Fatal(metrics.ListenMetrics("127.0.0.1:5001"))
	}()

	router.GET("/health", health)

	log.Println("http server init")
	address := fmt.Sprintf("%s:%s", cfg.HTTP.Host,cfg.HTTP.Port)
	server := http.Server{
		Handler:      router,
		Addr:         address,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	log.Printf("server running http://%s/health\n", address)
	log.Printf("docs http://%s/swagger/index.html\n", address)
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

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, DELETE, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
