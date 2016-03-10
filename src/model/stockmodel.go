package model

import (
	"time"
)

type StockSymbol struct {
	Code         string
	Name         string
	CategoryCode string
	Category     string
	StockTotal   uint
	Capital      int64
	FaceValue    int64
	Currency     string
	Tel          string
	Address      string
	TotalCount   uint
}

type StockLog struct {
	Id       uint
	Code     string
	Market   string
	Ask      string //매도가
	Bid      string //매수가
	Datetime time.Time
}
