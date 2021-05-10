package models

import (
	"time"
)

type Order struct {
	Id           int
	Offers       []int
	ClientId     int
	ManagerLogin string
	OrderDate    time.Time
	Quantity     int
	Status       string
}
