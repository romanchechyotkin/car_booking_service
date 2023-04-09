package car

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	car "github.com/romanchechyotkin/car_booking_service/internal/car/model"
	paymentproducer "github.com/romanchechyotkin/car_booking_service/internal/car/producer"
	car2 "github.com/romanchechyotkin/car_booking_service/internal/car/storage/cars_storage"
	"github.com/romanchechyotkin/car_booking_service/internal/car/storage/images_storage"
	"github.com/romanchechyotkin/car_booking_service/pkg/jwt"
	"log"
	"net/http"
	"strings"
	"time"
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
	DDMMYYYY = "02.01.2006"
)

type handler struct {
	carRepository   *car2.Repository
	imageRepository *images_storage.Repository
	paymentPlacer   *paymentproducer.PaymentPlacer
}

func NewHandler(carRep *car2.Repository, imgRep *images_storage.Repository, pp *paymentproducer.PaymentPlacer) *handler {
	return &handler{
		carRepository:   carRep,
		imageRepository: imgRep,
		paymentPlacer:   pp,
	}
}

func (h *handler) Register(router *gin.Engine) {
	router.POST("/cars", jwt.Middleware(h.CreateCar))
	router.GET("/cars", h.GetAllCars)
	router.GET("/cars/:id", h.GetCar)
	router.POST("/cars/:id/rent", jwt.Middleware(h.RentCar))
}

func (h *handler) CreateCar(ctx *gin.Context) {
	var formDto car.CreateCarFormDto

	err := ctx.ShouldBind(&formDto)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
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

	form, err := ctx.MultipartForm()
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

	id, err := token.GetIssuer()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	err = h.carRepository.CreateCar(ctx, &formDto, id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	files := form.File["image"]
	for _, file := range files {
		name := uuid.NewString()
		file.Filename = name

		err = ctx.SaveUploadedFile(file, "static/"+file.Filename)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		err = h.imageRepository.SaveImageToDB(ctx, file.Filename, formDto.Id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	ctx.JSON(http.StatusOK, formDto)
}

func (h *handler) GetAllCars(ctx *gin.Context) {
	cars, err := h.carRepository.GetAllCars(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, cars)
}

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

func (h *handler) RentCar(ctx *gin.Context) {
	carId := ctx.Param("id")

	c, err := h.carRepository.GetCar(ctx, carId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "no such car",
		})
		return
	}

	if c.IsAvailable == false {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "car is not available",
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

	customerId, err := token.GetIssuer()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	var rtd car.ReservationTimeDto
	err = ctx.ShouldBindJSON(&rtd)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	startDate, err := time.Parse(DDMMYYYY, rtd.StartDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	endDate, err := time.Parse(DDMMYYYY, rtd.EndDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	sub := endDate.Sub(startDate)
	days := sub.Hours() / 24

	price := c.PricePerDay * days

	carOwner, err := h.carRepository.GetCarOwner(ctx, carId)
	fmt.Println(carOwner)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var reservation = car.ReservationDto{
		Car:        c,
		CustomerId: customerId,
		CarOwnerId: carOwner,
		StartDate:  rtd.StartDate,
		EndDate:    rtd.EndDate,
		TotalPrice: price,
	}

	marshal, _ := json.Marshal(&reservation)
	log.Printf("payload: %s goes to kafka", string(marshal))

	err = h.paymentPlacer.SendPayment(marshal)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "server error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"reservation": reservation,
	})
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
