package model

import (
	"time"
)

type TradeLog struct {
	Id          uint
	UserEmail   string
	StockCode   string
	StockMarket string
	TradeType   string
	Price       int64
	Count       uint
	Timestamp   time.Time
}
