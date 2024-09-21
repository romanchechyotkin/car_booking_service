package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"sync"
	"time"

	carModel "github.com/romanchechyotkin/car_booking_service/internal/car/model"

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
	address := fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.Port)
	server := http.Server{
		Handler:      router,
		Addr:         address,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	log.Printf("server running http://%s/health\n", address)
	log.Printf("docs http://%s/swagger/index.html\n", address)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %s", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := FillData()
		if err != nil {
			log.Printf("Error in FillData: %v", err)
		}
	}()

	wg.Wait()
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

func FillData() error {
	requests := []struct {
		car   carModel.CreateCarFormDto
		image string
	}{
		{
			car: carModel.CreateCarFormDto{
				Id:          "1111AA-1",
				Brand:       "BMW",
				Model:       "X3",
				PricePerDay: 123,
				Year:        2024,
			},
			image: "./data/image1.png",
		},
		{
			car: carModel.CreateCarFormDto{
				Id:          "1111AA-2",
				Brand:       "Mercedes",
				Model:       "wasd",
				PricePerDay: 1,
				Year:        2000,
			},
			image: "./data/image2.png",
		},
		{
			car: carModel.CreateCarFormDto{
				Id:          "1111AA-3",
				Brand:       "BMW",
				Model:       "X5",
				PricePerDay: 1230,
				Year:        2015,
			},
			image: "./data/image3.png",
		},
		{
			car: carModel.CreateCarFormDto{
				Id:          "1111AA-4",
				Brand:       "Mercedes",
				Model:       "wasd",
				PricePerDay: 1000,
				Year:        2005,
			},
			image: "./data/image4.png",
		},
		{
			car: carModel.CreateCarFormDto{
				Id:          "1111BA-1",
				Brand:       "BMW",
				Model:       "X3",
				PricePerDay: 123,
				Year:        2024,
			},
			image: "./data/image5.png",
		},
		{
			car: carModel.CreateCarFormDto{
				Id:          "1111AA-5",
				Brand:       "Mercedes",
				Model:       "Lol",
				PricePerDay: 11,
				Year:        1999,
			},
			image: "./data/image6.png",
		},
		{
			car: carModel.CreateCarFormDto{
				Id:          "1111AA-6",
				Brand:       "LAda",
				Model:       "X5",
				PricePerDay: 111,
				Year:        2015,
			},
			image: "./data/image7.png",
		},
		{
			car: carModel.CreateCarFormDto{
				Id:          "1111AA-7",
				Brand:       "Honda",
				Model:       "wasd",
				PricePerDay: 120,
				Year:        2019,
			},
			image: "./data/image8.png",
		},
	}

	for _, req := range requests {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		mp3File, err := os.Open(req.image)
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer mp3File.Close()

		imageFile, err := os.Open(req.image)
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer imageFile.Close()

		mp3Part, err := writer.CreateFormFile("image", req.image)
		if err != nil {
			fmt.Println(err)
			continue
		}

		_, err = io.Copy(mp3Part, mp3File)
		if err != nil {
			fmt.Println(err)
			continue
		}

		writer.WriteField("id", req.car.Id)
		writer.WriteField("brand", req.car.Brand)
		writer.WriteField("model", req.car.Model)
		writer.WriteField("price", fmt.Sprintf("%.2f", req.car.PricePerDay))
		writer.WriteField("year", fmt.Sprintf("%d", req.car.Year))

		err = writer.Close()
		if err != nil {
			fmt.Println(err)
			continue
		}

		client := &http.Client{}

		var loginBody = []byte(fmt.Sprintf(`{"email": "%s", "password": "%s"}`, "admin@gmail.com", "admin"))
		r, err := http.NewRequest("POST", "http://localhost:8000/auth/login", bytes.NewReader(loginBody))
		if err != nil {
			fmt.Println(err)
			continue
		}

		response, err := client.Do(r)
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer response.Body.Close()

		responseBody, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Println(err)
			continue
		}

		var p struct {
			AccessToken string `json:"access_token"`
		}

		if err := json.Unmarshal(responseBody, &p); err != nil {
			fmt.Println(err)
			continue
		}

		r, err = http.NewRequest("POST", "http://localhost:8000/cars", body)
		if err != nil {
			fmt.Println(err)
			continue
		}

		r.Header.Set("Content-Type", writer.FormDataContentType())
		r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.AccessToken))

		response, err = client.Do(r)
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer response.Body.Close()

		responseBody, err = io.ReadAll(response.Body)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println(string(responseBody))
	}

	return nil
}
