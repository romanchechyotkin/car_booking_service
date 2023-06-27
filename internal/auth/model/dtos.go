package auth

import user "github.com/romanchechyotkin/car_booking_service/internal/user/model"

type RegistrationDto struct {
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=6"`
	FullName        string `json:"full_name" binding:"required"`
	TelephoneNumber string `json:"telephone_number" binding:"required,e164"`
	City            string `json:"city"`
}

type LoginDto struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResDto struct {
	AccessToken string           `json:"access_token"`
	User        user.GetUsersDto `json:"user"`
}
