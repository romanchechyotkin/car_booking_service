package car

import "time"

type Car struct {
	Id          string     `json:"id"`
	Brand       string     `json:"brand"`
	Model       string     `json:"model"`
	PricePerDay float64    `json:"pricePerDay"`
	Year        int        `json:"year"`
	IsAvailable bool       `json:"-"`
	Rating      float32    `json:"rating"`
    Location    string     `json:"location"`
    Seats       int        `json:"seats"`
    IsAutomatic bool       `json:"is_automatic"`
	Images      []string   `json:"images"`
	CreatedAt   *time.Time `json:"createdAt"`
}

type Image struct {
	Id    int    `json:"id"`
	Url   string `json:"url"`
	CarId string `json:"car_id"`
}
