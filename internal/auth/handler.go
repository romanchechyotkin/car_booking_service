package auth

import (
	"github.com/gin-gonic/gin"

	auth "github.com/romanchechyotkin/car_booking_service/internal/auth/model"

	"fmt"
	"net/http"
)

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
	//router.Handle(http.MethodPost, "/auth/registration", h.Logout)
	//router.Handle(http.MethodGet, "/auth/refresh", h.RefreshToken)
}

func (h *handler) Registration(ctx *gin.Context) {
	var body auth.RegistrationDto
	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("wrong entered data").Error()})
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

	err = h.service.Login(ctx, body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "successful login"})
}

func (h *handler) Logout(ctx *gin.Context) {

}

//func (h *handler) Registration(ctx *gin.Context) {
//
//}
