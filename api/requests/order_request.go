package requests

// struct for create order request body
type CreateOrderRequest struct {
	Origin      []string `json:"origin"`
	Destination []string `json:"destination"`
}


// struct for get order query strings
type GetOrderRequest struct {
	Page  int `form:"page"`
	Limit int `form:"limit"`
}

// struct for take order request body
type TakeOrderRequest struct {
	Status string `json:"status"`
}
