package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/romanchechyotkin/car_booking_service/pkg/jwt"

	auth "github.com/romanchechyotkin/car_booking_service/internal/auth/model"

	"errors"
	"fmt"
	"net/http"
	"strings"
)

var WrongEnteredPasswordError = errors.New("wrong entered password")
var EmptyFullNameError = errors.New("empty full name")
var WrongTelephoneNumberError = errors.New("wrong belarusian telephone number")

type handler struct {
	service *service
}

func NewHandler(s *service) *handler {
	return &handler{
		service: s,
	}
}

func (h *handler) Register(router *gin.Engine) {
	router.Handle(http.MethodPost, "/auth/registration", h.Registration)
	router.Handle(http.MethodPost, "/auth/login", h.Login)
	router.Handle(http.MethodGet, "/auth/logout", h.Logout)
	router.Handle(http.MethodGet, "/auth/refresh", h.RefreshToken)
}

func (h *handler) Registration(ctx *gin.Context) {
	var body auth.RegistrationDto
	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = ValidateForEmptyPasswordAndFullName(body.Password, body.FullName)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = ValidateBelarusTelephoneNumber(body.TelephoneNumber)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.service.Registration(ctx, body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "successful registration"})
}

func (h *handler) Login(ctx *gin.Context) {
	var body auth.LoginDto
	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("wrong entered data").Error()})
		return
	}

	user, token, err := h.service.Login(ctx, body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	refreshToken, err := jwt.GenerateRefreshToken(user.Id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	ctx.SetCookie("refresh_token", refreshToken, 24*60*60*1000, "/", "localhost", false, true) // 1 day

	ctx.JSON(http.StatusOK, gin.H{
		"access_token": token,
		"user":         user,
	})
}

func (h *handler) RefreshToken(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie("refresh_token")
	if err != nil {
		fmt.Println("cookie")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	if refreshToken == "" {
		fmt.Println("empty")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	token, err := jwt.ParseRefreshTokenToken(refreshToken)
	if err != nil {
		fmt.Println("hz")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	id, err := token.GetIssuer()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "server error",
		})
		return
	}

	u, err := h.service.repository.GetOneUserById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "server error",
		})
		return
	}

	generateAccessToken, err := jwt.GenerateAccessToken(u)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "server error",
		})
		return
	}

	generateRefreshToken, err := jwt.GenerateRefreshToken(u.Id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "server error",
		})
		return
	}

	ctx.SetCookie("refresh_token", generateRefreshToken, 24*60*60*1000, "/", "localhost", false, true)

	ctx.JSON(http.StatusOK, gin.H{
		"access_token": generateAccessToken,
	})
}

func (h *handler) Logout(ctx *gin.Context) {
	ctx.SetCookie("access_token", "", -1, "/", "127.0.01", false, false)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "log out",
	})
}

func ValidateForEmptyPasswordAndFullName(password, fullName string) error {
	if strings.Contains(password, " ") {
		return WrongEnteredPasswordError
	}

	password = strings.Trim(password, " ")
	fullName = strings.Trim(fullName, " ")

	if len(fullName) == 0 {
		return EmptyFullNameError
	}

	return nil
}

func ValidateBelarusTelephoneNumber(telephoneNumber string) error {
	if len(telephoneNumber) != 13 {
		return WrongTelephoneNumberError
	}

	if telephoneNumber[:4] != "+375" {
		return WrongTelephoneNumberError

	}

	codes := "24 25 29 33 44"
	code := telephoneNumber[4:6]

	if !strings.Contains(codes, code) {
		return WrongTelephoneNumberError
	}

	return nil
}
