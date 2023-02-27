package user

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	user2 "github.com/romanchechyotkin/car_booking-service/internal/user/model"
	"log"
	"net/http"
	"time"
)

type handler struct {
	service *Service
}

func NewHandler(service *Service) *handler {
	return &handler{service: service}
}

func (h *handler) Register(router *httprouter.Router) {
	router.Handle(http.MethodGet, "/users", h.GetALlUsers)
	router.Handle(http.MethodPost, "/users", h.CreateUser)
	//router.Handle(http.MethodGet, "/users/:id", h.GetOneUserById)
	//router.Handle(http.MethodPatch, "/users", h.UpdateUser)
	//router.Handle(http.MethodDelete, "/users/:id", h.DeleteUserById)
}

func (h *handler) GetALlUsers(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	start := time.Now()
	status := http.StatusOK
	defer func() {
		observeRequest(time.Since(start), status)
	}()

	users, err := h.service.FindAll()
	if err != nil {
		log.Println(err)
	}

	marshal, err := json.Marshal(users)
	w.WriteHeader(status)
	w.Write(marshal)
}

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

	fmt.Println(r.Body)

	var u user2.CreateUserDto
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		log.Println(err)
	}

	fmt.Printf("user: %s, %s, %s, %s", u.Email, u.Password, u.FullName, u.TelephoneNumber)

	err = h.service.CreateUser(context.Background(), &u)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "{error: %v}", err)
	} else {
		marshal, _ := json.Marshal(u)
		w.WriteHeader(http.StatusCreated)
		w.Write(marshal)
	}
}

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
