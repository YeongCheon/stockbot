package main

import (
	"stockbot/bot"
	stockdb "stockbot/db"
	"stockbot/logger"
	"stockbot/model"

	// "context"
	"io"
	"net/http"
	"os"
)

var (
	stockLogger  *logger.Logger
	stockCrawler bot.StockCrawler
)

func init() {
	stockLogger = logger.NewLogger()
	stockCrawler = bot.StockCrawler{
		StockLogger: stockLogger,
	}
}

func tradeStock(userEmail, tradeType string, symbol model.StockSymbol, cnt uint) bool {
	if tradeType != "ask" && tradeType != "bid" {
		stockLogger.Println("tradeType parameter value was wrong.")
		return false
	}
	stockLogList, err := stockCrawler.GetStockLogFromYahoo([]model.StockSymbol{symbol})
	if err != nil {
		stockLogger.Println(err)
		return false
	}

	stockLog := stockLogList[0]

	tradeLog := model.TradeLog{
		UserEmail:   "kyc1682@gmail.com",
		StockCode:   stockLog.Code,
		StockMarket: stockLog.Market,
		TradeType:   tradeType,
		Price:       stockLog.Ask,
		Count:       cnt,
	}
	_, err = stockdb.InsertTradeLog(tradeLog)
	if err == nil {
		return true
	} else {
		return false
	}
}

func commanderForWeb(commandCh chan<- string) {

	http.HandleFunc("/trade", func(w http.ResponseWriter, r *http.Request) {
	})

	http.HandleFunc("/exit", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "exit")
		os.Exit(1)
	})

	http.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		commandCh <- "stop"
		io.WriteString(w, "stop")
	})

	http.HandleFunc("/start", func(w http.ResponseWriter, r *http.Request) {
		commandCh <- "start"
		io.WriteString(w, "start")
	})

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "test")
	})

	stockLogger.Println(http.ListenAndServe(":9999", nil))
}

func main() {
	commandCh := make(chan string)

	//go CollectStockLog(commandCh)
	go stockCrawler.CollectStockData(commandCh)

	commanderForWeb(commandCh)
}
