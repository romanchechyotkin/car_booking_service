package request

import "github.com/romanchechyotkin/car_booking_service/internal/httpsrv/types"

type CreateRestaurantReq struct {
	Name       string              `form:"name"`
	Payments   types.PaymentArray  `form:"payments" binding:"required,dive,oneof=cash credit-card visa mastercard mir"`
	Deliveries types.DeliveryArray `form:"deliveries" binding:"required,dive,oneof=pick-up waiter"`
}

type CreateCategoryReq struct {
	Name string `form:"name"`
}

type CreateDishReq struct {
	Name          string  `form:"name"`
	Description   string  `form:"description"`
	Price         float64 `form:"price"`
	Weight        int     `form:"weight"`
	Calories      int     `form:"calories"`
	Proteins      int     `form:"proteins"`
	Fats          int     `form:"fats"`
	Carbohydrates int     `form:"carbohydrates"`
}

type CreateOrderReq struct {
	Positions []int64 `json:"positions"`
}

type GetOrdersReq struct {
	IDs []int64 `json:"orders"`
}
