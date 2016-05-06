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
	Price       uint
	Timestamp   time.Time
}
