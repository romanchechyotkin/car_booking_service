package car

type CreateCarDto struct {
	Id          string  `form:"id" json:"id"`
	Brand       string  `form:"brand" json:"brand"`
	Model       string  `form:"model" json:"model"`
	PricePerDay float64 `form:"price" json:"price"`
	Year        int     `form:"year" json:"year"`
}
