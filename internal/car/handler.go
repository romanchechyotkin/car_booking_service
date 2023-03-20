package car

import (
	"errors"
	"github.com/gin-gonic/gin"
	car "github.com/romanchechyotkin/car_booking_service/internal/car/model"
	"github.com/romanchechyotkin/car_booking_service/pkg/jwt"
	"net/http"
)

var (
	WrongCarNumbersLen                = errors.New("invalid car numbers length")
	WrongSymbolCarNumbers             = errors.New("no - ")
	WrongNumbersPartCarNumbers        = errors.New("invalid car numbers in numbers part")
	WrongLettersPartEnteredCarNumbers = errors.New("invalid car numbers in letters part")
	WrongRegionEnteredCarNumbers      = errors.New("invalid car numbers region")
)

type handler struct{}

func NewHandler() *handler {
	return &handler{}
}

func (h *handler) Register(router *gin.Engine) {
	router.POST("/cars", jwt.Middleware(h.CreateCar))
}

func (h *handler) CreateCar(ctx *gin.Context) {
	var dto car.CreateAutoDto

	err := ctx.ShouldBindJSON(&dto)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = ValidateCarNumbers(dto.Id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, "ok")
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
