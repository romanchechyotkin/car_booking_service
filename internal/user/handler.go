package user

import (
<<<<<<< HEAD
	"context"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	user2 "github.com/romanchechyotkin/car_booking-service/internal/user/model"
	"log"
=======
	"github.com/gin-gonic/gin"

	user3 "github.com/romanchechyotkin/car_booking_service/internal/user/metrics"
	user2 "github.com/romanchechyotkin/car_booking_service/internal/user/model"
	user "github.com/romanchechyotkin/car_booking_service/internal/user/storage"

>>>>>>> gin
	"net/http"
	"time"
)

type handler struct {
	service *Service
}

func NewHandler(service *Service) *handler {
	return &handler{service: service}
}

func (h *handler) Register(router *gin.Engine) {
	router.Handle(http.MethodGet, "/users", h.GetALlUsers)
	router.Handle(http.MethodPost, "/users", h.CreateUser)
<<<<<<< HEAD
	//router.Handle(http.MethodGet, "/users/:id", h.GetOneUserById)
	//router.Handle(http.MethodPatch, "/users", h.UpdateUser)
	//router.Handle(http.MethodDelete, "/users/:id", h.DeleteUserById)
=======
	router.Handle(http.MethodGet, "/users/:id", h.GetOneUserById)
	router.Handle(http.MethodPatch, "/users/:id", h.UpdateUser)
	router.Handle(http.MethodDelete, "/users/:id", h.DeleteUserById)
>>>>>>> gin
}

func (h *handler) GetALlUsers(ctx *gin.Context) {
	start := time.Now()
	status := http.StatusOK
	defer func() {
		user3.GetAllUsersObserveRequest(time.Since(start), status)
	}()

<<<<<<< HEAD
	users, err := h.service.FindAll()
=======
	users, err := h.repository.GetAllUsers(ctx)
>>>>>>> gin
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"users": users})
}

<<<<<<< HEAD
//	func (h *handler) GetOneUserById(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
//		id := params.ByName("id")
//
//		u, err := h.repository.GetOneUserById(context.Background(), id)
//		if err != nil {
//			w.WriteHeader(http.StatusNotFound)
//			http.NotFoundHandler()
//			log.Println(err)
//		}
//
//		marshal, err := json.Marshal(u)
//		w.WriteHeader(http.StatusOK)
//		w.Write(marshal)
//	}
func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	defer r.Body.Close()

	var u user2.CreateUserDto
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Printf("user: %s, %s, %s, %s", u.Email, u.Password, u.FullName, u.TelephoneNumber)

	err = h.service.CreateUser(context.Background(), &u)
=======
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
>>>>>>> gin
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "created",
	})
}

<<<<<<< HEAD
//func (h *handler) UpdateUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
//	defer r.Body.Close()
//
//	var body user2.UpdateUserDto
//	err := json.NewDecoder(r.Body).Decode(&body)
//	if err != nil {
//		log.Println(err)
//	}
//
//	err = h.repository.UpdateUser(context.Background(), &body)
//	if err != nil {
//		w.WriteHeader(http.StatusBadRequest)
//	}
//
//	marshal, _ := json.Marshal("updated")
//	w.Write(marshal)
//}
//
//func (h *handler) DeleteUserById(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
//	id := params.ByName("id")
//
//	err := h.repository.DeleteUserById(context.Background(), id)
//	if err != nil {
//		log.Println(err)
//	}
//
//	w.WriteHeader(http.StatusNoContent)
//	fmt.Fprintf(w, "deleted")
//}
=======
func (h *handler) UpdateUser(ctx *gin.Context) {
	start := time.Now()
	status := http.StatusOK
	defer func() {
		user3.UpdateUserObserveRequest(time.Since(start), status)
	}()

	id := ctx.Param("id")
	var uu user2.UpdateUserDto
	err := ctx.ShouldBindJSON(&uu)
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

func (h *handler) DeleteUserById(ctx *gin.Context) {
	start := time.Now()
	status := http.StatusOK
	defer func() {
		user3.DeleteUserObserveRequest(time.Since(start), status)
	}()

	id := ctx.Param("id")
	err := h.repository.DeleteUserById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "successfully deleted"})
}
>>>>>>> gin
