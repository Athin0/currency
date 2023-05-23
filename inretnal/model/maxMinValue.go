package model

import "time"

type ResponseCurrency struct {
	Name  string
	Value string
	Date  time.Time
}

type ResponseCurrencyAvg struct {
	Name  string
	Value string
}
