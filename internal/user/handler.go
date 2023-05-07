package user

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	user3 "github.com/romanchechyotkin/car_booking_service/internal/user/metrics"
	user2 "github.com/romanchechyotkin/car_booking_service/internal/user/model"
	user "github.com/romanchechyotkin/car_booking_service/internal/user/storage"
	"github.com/romanchechyotkin/car_booking_service/pkg/jwt"

	"net/http"
	"strings"
	"time"
)

var WrongRating = errors.New("wrong rating")

type handler struct {
	repository *user.Repository
}

func NewHandler(repository *user.Repository) *handler {
	return &handler{repository: repository}
}

func (h *handler) Register(router *gin.Engine) {
	router.Handle(http.MethodGet, "/users", h.GetALlUsers)
	router.Handle(http.MethodPost, "/users", h.CreateUser)
	router.Handle(http.MethodGet, "/users/:id", h.GetOneUserById)
	router.Handle(http.MethodPatch, "/users", jwt.Middleware(h.UpdateUser))
	router.Handle(http.MethodDelete, "/users", jwt.Middleware(h.DeleteUser))
	router.Handle(http.MethodGet, "/users/me", jwt.Middleware(h.GetMySelf))
	router.Handle(http.MethodPost, "/users/verify", jwt.Middleware(h.Verify))
	router.Handle(http.MethodGet, "/users/verify", jwt.Middleware(h.GetVerify))
	router.Handle(http.MethodPost, "/users/verify/:id", jwt.Middleware(h.VerifyUser))
	router.Handle(http.MethodPost, "/users/:id/rate", jwt.Middleware(h.RateUser))
	router.Handle(http.MethodGet, "/users/:id/rate", h.GetAllUserRates)
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
	ctx.JSON(http.StatusOK, users)
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

	ctx.JSON(http.StatusOK, userById)
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

	id := token["id"]

	var uu user2.UpdateUserDto
	err = ctx.ShouldBindJSON(&uu)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.repository.UpdateUser(ctx, fmt.Sprintf("%s", id), &uu)
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

	id := token["id"]

	err = h.repository.DeleteUserById(ctx, fmt.Sprintf("%s", id))
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

	id := token["id"]

	userById, err := h.repository.GetOneUserById(ctx, fmt.Sprintf("%s", id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	ctx.JSON(http.StatusOK, userById)
}

func (h *handler) Verify(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	headers := strings.Split(authHeader, " ")

	token, err := jwt.ParseAccessToken(headers[1])
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id := token["id"]

	form, err := ctx.FormFile("image")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	form.Filename = uuid.NewString()
	err = ctx.SaveUploadedFile(form, "static/users/"+form.Filename)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = h.repository.CreateApplication(ctx, fmt.Sprintf("%s", id), form.Filename)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "u have already sent",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "your application is successfully sent",
	})
}

func (h *handler) GetVerify(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	headers := strings.Split(authHeader, " ")

	token, err := jwt.ParseAccessToken(headers[1])
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	role := token["role"]
	if role != "ADMIN" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "no rights"})
		return
	}

	applications, err := h.repository.GetApplications(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, err)
		return
	}

	ctx.JSON(http.StatusOK, applications)
}

func (h *handler) VerifyUser(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	headers := strings.Split(authHeader, " ")

	token, err := jwt.ParseAccessToken(headers[1])
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userId := ctx.Param("id")
	role := token["role"]
	if role != "ADMIN" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "no rights"})
		return
	}

	err = h.repository.ChangeUserVerify(ctx, userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "server error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("user %s verified", userId),
	})
}

func (h *handler) RateUser(ctx *gin.Context) {
	ratedUserId := ctx.Param("id")

	authHeader := ctx.GetHeader("Authorization")
	headers := strings.Split(authHeader, " ")

	token, err := jwt.ParseAccessToken(headers[1])
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id := token["id"]

	user, err := h.repository.GetOneUserById(ctx, fmt.Sprintf("%s", id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	userForRating, err := h.repository.GetOneUserById(ctx, ratedUserId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	if user.Id == userForRating.Id {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "u cant rate yourself"})
		return
	}

	var dto user2.RateDto
	err = ctx.ShouldBindJSON(&dto)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = ValidateRating(dto.Rating)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = ValidateCommentLength(dto.Comment)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = h.repository.CreateRating(ctx, dto, userForRating.Id, user.Id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	amount, sum, err := h.repository.GetUserRatings(ctx, userForRating.Id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	rate := sum / amount

	err = h.repository.ChangeUserRating(ctx, userForRating.Id, rate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "server error",
		})
		return
	}

	userForRating.Rating = rate

	ctx.JSON(http.StatusOK, gin.H{
		"me":        user,
		"ratedUser": userForRating,
	})
}

func (h *handler) GetAllUserRates(ctx *gin.Context) {
	userId := ctx.Param("id")

	ratings, err := h.repository.GetAllUserRatings(ctx, userId)
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

func ValidateRating(rating float32) error {
	if rating > 5 || rating < 1 {
		return WrongRating
	}
	return nil
}

func ValidateCommentLength(comment string) error {
	if len(comment) > 250 {
		return errors.New("long comment")
	}
	return nil
}
