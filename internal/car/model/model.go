package car

type Car struct {
	Id          string  `json:"id"`
	Brand       string  `json:"brand"`
	Model       string  `json:"model"`
	PricePerDay float64 `json:"pricePerDay"`
	Year        int     `json:"year"`
	IsAvailable bool    `json:"isAvailable"`
	Rating      float32 `json:"rating"`
	Images      []Image `json:"images"`
}

type Image struct {
	Id    int    `json:"id"`
	Url   string `json:"url"`
	CarId string `json:"car_id"`
}
