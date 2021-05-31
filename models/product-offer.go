package models

type Offer struct {
	Id         int     `json:"id"`
	ProviderId int     `json:"provider_id"`
	ProductId  int     `json:"product_id"`
	Cost       float32 `json:"cost"`
}
