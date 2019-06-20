package requests

type Order struct {
	Origin []string `json:"origin"`
	Destination []string `json:"destination"`
} 