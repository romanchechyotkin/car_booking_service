package car

import (
	"context"
	"errors"
	"fmt"
	"github.com/minio/minio-go/v7"
	"log"
	"math"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	car "github.com/romanchechyotkin/car_booking_service/internal/car/model"
	car2 "github.com/romanchechyotkin/car_booking_service/internal/car/storage/cars_storage"
	"github.com/romanchechyotkin/car_booking_service/internal/car/storage/images_storage"
	"github.com/romanchechyotkin/car_booking_service/internal/reservation/model"
	res2 "github.com/romanchechyotkin/car_booking_service/internal/reservation/storage"
	user3 "github.com/romanchechyotkin/car_booking_service/internal/user"
	user2 "github.com/romanchechyotkin/car_booking_service/internal/user/model"
	user "github.com/romanchechyotkin/car_booking_service/internal/user/storage"
	"github.com/romanchechyotkin/car_booking_service/pkg/jwt"
	"github.com/romanchechyotkin/car_booking_service/proto/pb"
)

var (
	EmptyString                       = errors.New("empty string")
	WrongCarNumbersLen                = errors.New("invalid car numbers length")
	WrongSymbolCarNumbers             = errors.New("no - ")
	WrongNumbersPartCarNumbers        = errors.New("invalid car numbers in numbers part")
	WrongLettersPartEnteredCarNumbers = errors.New("invalid car numbers in letters part")
	WrongRegionEnteredCarNumbers      = errors.New("invalid car numbers region")
)

const (
	hhDDMMYYYY         = "15 02.01.2006"
	SORT_BY_ASC_PRICE  = "prc.a"
	SORT_BY_DESC_PRICE = "prc.d"
	SORT_BY_ASC_YEAR   = "y.a"
	SORT_BY_DESC_YEAR  = "y.d"
)

type handler struct {
	carRepository   car2.Storage
	imageRepository images_storage.ImageStorage
	reservationRep  res2.Storage
	userRep         user.Storage
	grpcClient      pb.CarsManagementClient
	minioClient     *minio.Client
}

func NewHandler(carRep car2.Storage, imgRep images_storage.ImageStorage, resRep res2.Storage, up user.Storage, grpcClient pb.CarsManagementClient, minioClient *minio.Client) *handler {
	return &handler{
		carRepository:   carRep,
		imageRepository: imgRep,
		reservationRep:  resRep,
		userRep:         up,
		grpcClient:      grpcClient,
		minioClient:     minioClient,
	}
}

func (h *handler) Register(router *gin.Engine) {
	router.POST("/cars", jwt.Middleware(h.CreateCar))
	router.GET("/cars", h.GetAllCars)
	router.GET("/cars/:id", h.GetCar)
	router.POST("/cars/:id/rent", jwt.Middleware(h.RentCar))
	router.POST("/cars/:id/rate", jwt.Middleware(h.RateCar))
	router.GET("/cars/:id/rate", h.GetAllCarRatings)
	router.GET("/cars/:id/reservations", h.GetAllCarReservations)
}

// CreateCar godoc
// @Tags cars
// @Security BearerAuth
// @Summary CreateCar
// @Description Endpoint for creating car post
// @Produce application/json
// @ID image
// @Accept multipart/form-data
// @Param id formData string true "id"
// @Param brand formData string true "brand"
// @Param model formData string true "model"
// @Param price formData float64 true "price"
// @Param year formData int true "year"
// @Param image formData file true "Image file"
// @Success 201 {object} car.CreateCarFormDto{}
// @Router /cars [post]
func (h *handler) CreateCar(ctx *gin.Context) {
	var formDto car.CreateCarFormDto

	err := ctx.ShouldBind(&formDto)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	fmt.Println(formDto)

	form, err := ctx.MultipartForm()
	files := form.File["image"]
	if len(files) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "no photos",
		})
		return
	}

	err = ValidateCarNumbers(formDto.Id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	brand, model, err := ValidateForEmptyStrings(formDto.Brand, formDto.Model)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	formDto.Brand = brand
	formDto.Model = model

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	authHeader := ctx.GetHeader("Authorization")
	headers := strings.Split(authHeader, " ")

	token, err := jwt.ParseAccessToken(headers[1])
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userId := token["id"]

	err = h.carRepository.CreateCar(ctx, &formDto, fmt.Sprintf("%s", userId))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	for _, file := range files {
		name := uuid.NewString()
		file.Filename = name + ".png"

		err = ctx.SaveUploadedFile(file, "static/cars/"+file.Filename)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		err = saveToMinio(ctx, h.minioClient, file.Filename, "static/cars/"+file.Filename)
		if err != nil {
			log.Println(err)
		}

		err = h.imageRepository.SaveImageToDB(ctx, file.Filename, formDto.Id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		err = os.Remove("static/cars/" + file.Filename)
		if err != nil {
			log.Println(err)
		}
	}

	usr, err := h.userRep.GetOneUserById(ctx, userId.(string))
	if err != nil {
		log.Println(err)
	}

	grpcReq := &pb.CreateCarDocReq{
		CarId:       formDto.Id,
		Brand:       formDto.Brand,
		Model:       formDto.Model,
		Year:        uint32(formDto.Year),
		PricePerDay: float32(formDto.PricePerDay),
		IsAvailable: false,
		Rating:      0.0,
		Username:    usr.FullName,
		City:        usr.City,
	}

	res, err := h.grpcClient.CreateCarDoc(ctx, grpcReq)
	if err != nil {
		log.Println("GRPC ERROR", err)
	}
	log.Println("GRPC RESPONSE", res)

	ctx.JSON(http.StatusCreated, formDto)
}

// GetAllCars godoc
// @Tags cars
// @Summary GetAllCars
// @Description Endpoint for getting all cars posts
// @Produce application/json
// @Success 200 {object} []car.GetCarDto{}
// @Router /cars [get]
func (h *handler) GetAllCars(ctx *gin.Context) {
	value := ctx.Query("sort")
	log.Println(value)

	var cars []car.GetCarDto
	var err error

	switch value {
	case SORT_BY_ASC_PRICE:
		cars, err = h.carRepository.GetAllCars(ctx, value)
	case SORT_BY_DESC_PRICE:
		cars, err = h.carRepository.GetAllCars(ctx, value)
	case SORT_BY_ASC_YEAR:
		cars, err = h.carRepository.GetAllCars(ctx, value)
	case SORT_BY_DESC_YEAR:
		cars, err = h.carRepository.GetAllCars(ctx, value)
	default:
		cars, err = h.carRepository.GetAllCars(ctx)
	}

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, cars)
}

// GetCar godoc
// @Tags cars
// @Summary GetCar
// @Description Endpoint for getting one car info by its id
// @Produce application/json
// @Param id path string true "Car ID"
// @Success 200 {object} car.GetCarDto{}
// @Router /cars/{id} [get]
func (h *handler) GetCar(ctx *gin.Context) {
	id := ctx.Param("id")

	c, err := h.carRepository.GetCar(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "car not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, c)
}

// RentCar godoc
// @Tags cars
// @Security BearerAuth
// @Summary RentCar
// @Description Endpoint for renting cars
// @Produce application/json
// @Param id path string true "Car ID"
// @Param body body reservation.TimeDto{} true "Times"
// @Success 201 {object} car.CreateCarFormDto{}
// @Router /cars/{id}/rent [post]
func (h *handler) RentCar(ctx *gin.Context) {
	carId := ctx.Param("id")

	c, err := h.carRepository.GetCar(ctx, carId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "no such car",
		})
		return
	}

	authHeader := ctx.GetHeader("Authorization")
	headers := strings.Split(authHeader, " ")

	token, err := jwt.ParseAccessToken(headers[1])
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	customerId := token["id"]
	isVerified := token["is_verified"]
	if !isVerified.(bool) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "you are not verified customer for driving",
		})
		return
	}

	var rtd reservation.TimeDto
	err = ctx.ShouldBindJSON(&rtd)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	startDate, err := time.Parse(hhDDMMYYYY, rtd.StartDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	endDate, err := time.Parse(hhDDMMYYYY, rtd.EndDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	now := time.Now()
	if startDate.Before(now) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "u wanna to rent car in the past",
		})
		return
	}
	log.Println(startDate, endDate)

	dates, err := h.reservationRep.GetReservationDates(ctx, c.Id)
	log.Println(dates)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	for _, v := range dates {
		if (startDate.After(v.StartDate) && startDate.Before(v.EndDate)) || (endDate.After(v.StartDate) && endDate.Before(v.EndDate)) {
			log.Println(startDate, v.StartDate)
			log.Println(endDate, v.EndDate)
			log.Printf("equal dates %d-%d or %d-%d\n", startDate.Day(), v.StartDate.Day(), endDate.Day(), v.EndDate.Day())
			ctx.JSON(http.StatusBadRequest, gin.H{
				"msg": "car is busy this time",
			})
			return
		}
	}

	sub := endDate.Sub(startDate)
	days := sub.Hours() / 24
	if sub.Milliseconds() < 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "negative time",
		})
		return
	}
	price := math.Round(c.PricePerDay * days)

	carOwner, err := h.carRepository.GetCarOwner(ctx, carId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var rsv = reservation.Dto{
		Car:        c,
		CustomerId: fmt.Sprintf("%s", customerId),
		CarOwnerId: carOwner,
		StartDate:  rtd.StartDate,
		EndDate:    rtd.EndDate,
		TotalPrice: price,
	}

	err = h.reservationRep.CreateReservation(ctx, rsv)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "server error",
		})
		return
	}

	// kafka DONT TOUCH
	//marshal, _ := json.Marshal(&reservation)
	//log.Printf("payload: %s goes to kafka", string(marshal))
	//
	//err = h.paymentPlacer.SendPayment(marshal)
	//log.Println(err)

	err = h.carRepository.ChangeIsAvailable(ctx, c.Id)
	log.Printf("error due change availability %v", err)

	ctx.JSON(http.StatusOK, gin.H{
		"reservation": rsv,
	})
}

// RateCar godoc
// @Tags cars
// @Security BearerAuth
// @Summary RateCar
// @Description Endpoint for rating cars
// @Produce application/json
// @Param id path string true "Car ID"
// @Param body body user2.RateDto{} true "Rate"
// @Success 201 {object} car.CreateCarFormDto{}
// @Router /cars/{id}/rate [post]
func (h *handler) RateCar(ctx *gin.Context) {
	carId := ctx.Param("id")

	authHeader := ctx.GetHeader("Authorization")
	headers := strings.Split(authHeader, " ")

	token, err := jwt.ParseAccessToken(headers[1])
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id := token["id"]

	usr, err := h.userRep.GetOneUserById(ctx, fmt.Sprintf("%s", id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	var dto user2.RateDto
	err = ctx.ShouldBindJSON(&dto)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = user3.ValidateRating(dto.Rating)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = user3.ValidateCommentLength(dto.Comment)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = h.carRepository.CreateRating(ctx, dto, carId, usr.Id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	amount, sum, err := h.carRepository.GetAmountCarRatings(ctx, carId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	rate := sum / amount

	err = h.carRepository.ChangeCarRating(ctx, carId, rate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "server error",
		})
		return
	}

	c, err := h.carRepository.GetCar(ctx, carId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "server error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"car":  c,
		"user": usr.FullName,
	})
}

// GetAllCarRatings godoc
// @Tags cars
// @Summary GetAllCarRatings
// @Description Endpoint for getting all cars rates
// @Produce application/json
// @Param id path string true "Car ID"
// @Success 200 {object} []car.GetAllCarRatingsDto{}
// @Router /cars/{id}/rate [get]
func (h *handler) GetAllCarRatings(ctx *gin.Context) {
	carId := ctx.Param("id")

	ratings, err := h.carRepository.GetAllCarRatings(ctx, carId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if len(ratings) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg": "no ratings",
		})
		return
	}

	ctx.JSON(http.StatusOK, ratings)
}

func (h *handler) GetAllCarReservations(ctx *gin.Context) {
	carId := ctx.Param("id")

	dates, err := h.reservationRep.GetReservationDates(ctx, carId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if len(dates) == 0 {
		ctx.JSON(http.StatusOK, "no reservations")
		return
	}

	var resDates = make([]reservation.GetResDto, 0, len(dates))

	for _, d := range dates {
		resDates = append(resDates, reservation.GetResDto{
			StartDate: d.StartDate.Format(hhDDMMYYYY),
			EndDate:   d.EndDate.Format(hhDDMMYYYY),
		})
	}

	ctx.JSON(http.StatusOK, resDates)
}

func ValidateCarNumbers(numbers string) error {
	if len(numbers) != 8 {
		return WrongCarNumbersLen
	}

	if string(numbers[6]) != "-" {
		return WrongSymbolCarNumbers
	}

	if numbers[7] < 49 || numbers[7] > 55 {
		return WrongRegionEnteredCarNumbers
	}

	for i := 0; i < 4; i++ {
		if numbers[i] > 57 || numbers[i] < 48 {
			return WrongNumbersPartCarNumbers
		}
	}

	for i := 4; i < 6; i++ {
		if numbers[i] > 90 || numbers[i] < 65 {
			return WrongLettersPartEnteredCarNumbers
		}
	}

	return nil
}

func ValidateForEmptyStrings(brand, model string) (string, string, error) {
	brandTrim := strings.Trim(brand, " ")
	if len(brandTrim) == 0 {
		return "", "", EmptyString
	}

	modelTrim := strings.Trim(model, " ")
	if len(modelTrim) == 0 {
		return "", "", EmptyString
	}

	return brandTrim, modelTrim, nil
}

func saveToMinio(ctx context.Context, client *minio.Client, fileName, filePath string) error {
	info, err := client.FPutObject(ctx, "Admin-bucket", fileName, filePath, minio.PutObjectOptions{ContentType: "image/png"})
	if err != nil {
		return err
	}

	log.Printf("Successfully uploaded %s of size %d\n", fileName, info.Size)

	return nil
}
