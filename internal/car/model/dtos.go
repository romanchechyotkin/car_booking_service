package car

type CreateAutoDto struct {
	Id          string  `json:"id" binding:"required"`
	Brand       string  `json:"brand" binding:"required"`
	Model       string  `json:"model" binding:"required"`
	PricePerDay float64 `json:"pricePerDay" binding:"required"`
	Year        int     `json:"year" binding:"required"`
}
