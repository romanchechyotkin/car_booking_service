package car

import "time"

type CreateCarFormDto struct {
	Id          string  `form:"id" json:"id"`
	Brand       string  `form:"brand" json:"brand"`
	Model       string  `form:"model" json:"model"`
	PricePerDay float64 `form:"price" json:"price"`
	Year        int     `form:"year" json:"year"`
    Seats       int     `form:"seats" json:"seats"`
    IsAutomatic bool    `form:"is_automatic" json:"is_automatic"`   
}

type GetCarDto struct {
	Car    Car    `json:"car"`
	UserId string `json:"user_id"`
}

type GetAllCarRatingsDto struct {
	Rating    float32   `json:"rating"`
	Comment   string    `json:"comment"`
	User      string    `json:"user"`
	CreatedAt time.Time `json:"created_at"`
}
