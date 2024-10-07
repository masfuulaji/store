package models

type Order struct {
	Cart     Cart       `json:"cart"`
	CartItem []CartItem `json:"cart_item"`
}
