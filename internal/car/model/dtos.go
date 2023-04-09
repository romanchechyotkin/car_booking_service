package car

type CreateCarFormDto struct {
	Id          string  `form:"id" json:"id"`
	Brand       string  `form:"brand" json:"brand"`
	Model       string  `form:"model" json:"model"`
	PricePerDay float64 `form:"price" json:"price"`
	Year        int     `form:"year" json:"year"`
}

type GetCarDto struct {
	Car    Car    `json:"car"`
	UserId string `json:"user_id"`
}

type ReservationTimeDto struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type ReservationDto struct {
	CustomerId string  `json:"customer_id"`
	Car        Car     `json:"car"`
	CarOwnerId string  `json:"car_owner_id"`
	StartDate  string  `json:"start_date"`
	EndDate    string  `json:"end_date"`
	TotalPrice float64 `json:"total_price"`
}
