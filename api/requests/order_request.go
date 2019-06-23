package requests

type CreateOrderRequest struct {
	Origin []string `json:"origin"`
	Destination []string `json:"destination"`
}

type GetOrderRequest struct {
	Page int `form:"page"`
	Limit int `form:"limit"`
}

type TakeOrderRequest struct {
	Status string `json:"status"`
}