package car

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	car "github.com/romanchechyotkin/car_booking_service/internal/car/model"
	paymentproducer "github.com/romanchechyotkin/car_booking_service/internal/car/producer"
	car2 "github.com/romanchechyotkin/car_booking_service/internal/car/storage/cars_storage"
	"github.com/romanchechyotkin/car_booking_service/internal/car/storage/images_storage"
	"github.com/romanchechyotkin/car_booking_service/internal/reservation/model"
	res2 "github.com/romanchechyotkin/car_booking_service/internal/reservation/storage"
	user3 "github.com/romanchechyotkin/car_booking_service/internal/user"
	user2 "github.com/romanchechyotkin/car_booking_service/internal/user/model"
	user "github.com/romanchechyotkin/car_booking_service/internal/user/storage"
	"github.com/romanchechyotkin/car_booking_service/pkg/jwt"
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
	reservationRep  *res2.Repository
	userRep         *user.Repository
}

func NewHandler(carRep *car2.Repository, imgRep *images_storage.Repository, pp *paymentproducer.PaymentPlacer, resRep *res2.Repository, up *user.Repository) *handler {
	return &handler{
		carRepository:   carRep,
		imageRepository: imgRep,
		paymentPlacer:   pp,
		reservationRep:  resRep,
		userRep:         up,
	}
}

func (h *handler) Register(router *gin.Engine) {
	router.POST("/cars", jwt.Middleware(h.CreateCar))
	router.GET("/cars", h.GetAllCars)
	router.GET("/cars/:id", h.GetCar)
	router.POST("/cars/:id/rent", jwt.Middleware(h.RentCar))
	router.POST("/cars/:id/rate", jwt.Middleware(h.RateCar))
	router.GET("/cars/:id/rate", h.GetAllCarRatings)
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

// TODO: transaction for reservation db and change availability

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

	customerId, err := token.GetIssuer()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
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

	dates, err := h.reservationRep.GetReservationDates(ctx, c.Id)
	fmt.Println(dates)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	for _, v := range dates {
		if startDate == v.StartDate || endDate == v.EndDate {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"msg": "car is busy this time",
			})
			return
		}
	}

	sub := endDate.Sub(startDate)
	days := sub.Hours() / 24

	price := c.PricePerDay * days

	carOwner, err := h.carRepository.GetCarOwner(ctx, carId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var reservation = reservation.Dto{
		Car:        c,
		CustomerId: customerId,
		CarOwnerId: carOwner,
		StartDate:  rtd.StartDate,
		EndDate:    rtd.EndDate,
		TotalPrice: price,
	}

	err = h.reservationRep.CreateReservation(ctx, reservation)
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

	//err = h.carRepository.ChangeIsAvailable(ctx, c.Id)
	//log.Printf("error due change availability %v", err)

	ctx.JSON(http.StatusOK, gin.H{
		"reservation": reservation,
	})
}

func (h *handler) RateCar(ctx *gin.Context) {
	carId := ctx.Param("id")

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

	user, err := h.userRep.GetOneUserById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// TODO check validation for rate your own car

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

	err = h.carRepository.CreateRating(ctx, dto, carId, user.Id)
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
		"user": user.FullName,
	})
}

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
