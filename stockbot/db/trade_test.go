package db

import (
	. "stockbot/model"

	"fmt"
	"testing"
	"time"
)

func TestInsertTradeLog(t *testing.T) {
	tradeLog := TradeLog{
		UserEmail:   "kyc1682@gmail.com",
		StockCode:   "",
		StockMarket: "",
		TradeType:   "",
		Price:       9999,
		Timestamp:   time.Now(),
	}
	fmt.Println(InsertTradeLog(tradeLog))
}

func TestSelectTradeLog(t *testing.T) {
}
