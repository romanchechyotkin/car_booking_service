package car

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	car "github.com/romanchechyotkin/car_booking_service/internal/car/model"
	car2 "github.com/romanchechyotkin/car_booking_service/internal/car/storage/cars_storage"
	"github.com/romanchechyotkin/car_booking_service/internal/car/storage/images_storage"
	"github.com/romanchechyotkin/car_booking_service/pkg/jwt"
	"net/http"
	"strings"
)

var (
	EmptyString                       = errors.New("empty string")
	WrongCarNumbersLen                = errors.New("invalid car numbers length")
	WrongSymbolCarNumbers             = errors.New("no - ")
	WrongNumbersPartCarNumbers        = errors.New("invalid car numbers in numbers part")
	WrongLettersPartEnteredCarNumbers = errors.New("invalid car numbers in letters part")
	WrongRegionEnteredCarNumbers      = errors.New("invalid car numbers region")
)

type handler struct {
	carRepository   *car2.Repository
	imageRepository *images_storage.Repository
}

func NewHandler(carRep *car2.Repository, imgRep *images_storage.Repository) *handler {
	return &handler{
		carRepository:   carRep,
		imageRepository: imgRep,
	}
}

// TODO: finish create car route with jwt auth (which user car)

func (h *handler) Register(router *gin.Engine) {
	router.POST("/cars", jwt.Middleware(h.CreateCar))
	//router.POST("/cars", h.CreateCar)
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
