package user

import (
	"github.com/gin-gonic/gin"
	"github.com/romanchechyotkin/car_booking_service/pkg/jwt"
	"strings"

	user3 "github.com/romanchechyotkin/car_booking_service/internal/user/metrics"
	user2 "github.com/romanchechyotkin/car_booking_service/internal/user/model"
	user "github.com/romanchechyotkin/car_booking_service/internal/user/storage"

	"net/http"
	"time"
)

type handler struct {
	repository *user.Repository
}

func NewHandler(repository *user.Repository) *handler {
	return &handler{repository: repository}
}

func (h *handler) Register(router *gin.Engine) {
	router.Handle(http.MethodGet, "/users", jwt.Middleware(h.GetALlUsers))
	router.Handle(http.MethodPost, "/users", h.CreateUser)
	router.Handle(http.MethodGet, "/users/:id", h.GetOneUserById)
	router.Handle(http.MethodPatch, "/users", jwt.Middleware(h.UpdateUser))
	router.Handle(http.MethodDelete, "/users", jwt.Middleware(h.DeleteUser))
	router.Handle(http.MethodGet, "/users/me", jwt.Middleware(h.GetMySelf))
}

func (h *handler) GetALlUsers(ctx *gin.Context) {
	start := time.Now()
	status := http.StatusOK
	defer func() {
		user3.GetAllUsersObserveRequest(time.Since(start), status)
	}()

	users, err := h.repository.GetAllUsers(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"users": users})
}

func (h *handler) GetOneUserById(ctx *gin.Context) {
	start := time.Now()
	status := http.StatusOK
	defer func() {
		user3.GetOneUserByIdObserveRequest(time.Since(start), status)
	}()

	id := ctx.Param("id")
	userById, err := h.repository.GetOneUserById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user": userById})
}

func (h *handler) CreateUser(ctx *gin.Context) {
	start := time.Now()
	status := http.StatusOK
	defer func() {
		user3.CreateUserObserveRequest(time.Since(start), status)
	}()

	var cu user2.CreateUserDto
	err := ctx.ShouldBindJSON(&cu)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.repository.CreateUser(ctx, &cu)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "created",
	})
}

func (h *handler) UpdateUser(ctx *gin.Context) {
	start := time.Now()
	status := http.StatusOK
	defer func() {
		user3.UpdateUserObserveRequest(time.Since(start), status)
	}()

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

	var uu user2.UpdateUserDto
	err = ctx.ShouldBindJSON(&uu)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.repository.UpdateUser(ctx, id, &uu)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "updated",
	})
}

func (h *handler) DeleteUser(ctx *gin.Context) {
	start := time.Now()
	status := http.StatusOK
	defer func() {
		user3.DeleteUserObserveRequest(time.Since(start), status)
	}()

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

	err = h.repository.DeleteUserById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "successfully deleted"})
}

func (h *handler) GetMySelf(ctx *gin.Context) {
	start := time.Now()
	status := http.StatusOK
	defer func() {
		user3.GetMySelfObserveRequest(time.Since(start), status)
	}()

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

	userById, err := h.repository.GetOneUserById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	ctx.JSON(http.StatusOK, userById)
}
