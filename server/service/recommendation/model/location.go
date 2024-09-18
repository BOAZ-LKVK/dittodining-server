package model

import "github.com/shopspring/decimal"

type Location struct {
	Latitude  decimal.Decimal `json:"latitude"`
	Longitude decimal.Decimal `json:"longitude"`
}
