package reservation

import (
	car "github.com/romanchechyotkin/car_booking_service/internal/car/model"
	"time"
)

type TimeDto struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type TimeFromDB struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

type Dto struct {
	CustomerId string  `json:"customer_id"`
	Car        car.Car `json:"car"`
	CarOwnerId string  `json:"car_owner_id"`
	StartDate  string  `json:"start_date"`
	EndDate    string  `json:"end_date"`
	TotalPrice float64 `json:"total_price"`
}

type GetResDto struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}
