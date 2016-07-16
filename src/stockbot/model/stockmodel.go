package model

import (
	"time"
)

type StockSymbol struct {
	Code         string
	Name         string
	Market       string
	CategoryCode string
	Category     string
	StockTotal   uint
	Capital      float64
	FaceValue    float64
	Currency     string
	Tel          string
	Address      string
	TotalCount   uint
}

type StockLog struct {
	Id       uint
	Code     string
	Market   string
	Ask      float64 //매도가
	Bid      float64 //매수가
	Datetime time.Time
}
