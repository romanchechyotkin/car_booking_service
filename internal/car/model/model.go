package car

type Auto struct {
	Id          string  `json:"id"`
	Brand       string  `json:"brand"`
	Model       string  `json:"model"`
	PricePerDay float64 `json:"pricePerDay"`
	Year        int     `json:"year"`
	IsAvailable bool    `json:"isAvailable"`
	Rating      float32 `json:"rating"`
}
