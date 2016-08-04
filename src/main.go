package main

import (
	stockdb "stockbot/db"
	"stockbot/logger"
	"stockbot/model"

	"encoding/csv"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var stockLogger *logger.Logger

func init() {
	stockLogger = logger.NewLogger()
	stockLogger.Println("hello world!!!!!!!!!!!!!")
}

func GetStockLogFromYahoo(stockSymbolList []model.StockSymbol) ([]model.StockLog, error) {
	var s string
	for _, stockSymbol := range stockSymbolList {
		s += stockSymbol.Code + ".KS+"
	}
	f := "f=" + "snabt1" // format: symbol, name, ask, bid, last trade time
	url := "http://finance.yahoo.com/d/quotes.csv?s=" + s + "&" + f
	res, err := http.Get(url)
	if err != nil {
		stockLogger.Println(err)
	}

	result, err := ioutil.ReadAll(res.Body)
	if err != nil {
		stockLogger.Println(err)
	}

	stockLogs := ParseYahooCSV(result)
	if len(stockLogs) <= 0 {
		return nil, errors.New("GetStockLogFromYahoo func's result is nil")
	}

	for idx, _ := range stockLogs {
		stockLogs[idx].Code = stockLogs[idx].Code[:6]
		stockLogs[idx].Market = stockSymbolList[idx].Market
	}

	return stockLogs, nil
}

func ParseYahooCSV(target []byte) []model.StockLog {
	var result []model.StockLog
	reader := csv.NewReader(strings.NewReader(string(target)))
	defer func() {
		if r := recover(); r != nil {
			stockLogger.Println(r)
		}
	}()
	for {
		stockLog := model.StockLog{}

		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			stockLogger.Println(err)
		}

		stockLog.Code = record[0]
		ask, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			stockLogger.Println(err)
		}
		bid, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			stockLogger.Println(err)
		}

		stockLog.Ask = ask
		stockLog.Bid = bid

		result = append(result, stockLog)
	}

	return result
}

func CollectStockLog(controller chan string) {
	stockLogger.Println("start stock log crawling...")
	const MARKETSTARTTIME int = 9
	const MARKETENDTIME int = 15

	const maxGoroutineCnt int = 50
	limitCh := make(chan bool, maxGoroutineCnt)

	/*
		timezone, err := time.LoadLocation("Asia/Seoul")
		if err != nil {
			stockLogger.Println(err)
		}
	*/

	stockList := stockdb.GetAllStockList()

	isStop := false
	for {
		select {
		case param := <-controller:
			switch param {
			case "start":
				isStop = false
			case "stop":
				isStop = true
			}
		default:
		}

		/*
			now := time.Now()

			if now.Weekday() == 0 || now.Weekday() == 6 {
				//fmt.Println("today is weekend")
			} else if now.Hour() < MARKETSTARTTIME || now.Hour() >= MARKETENDTIME {
				//fmt.Println("market was closed")
			} else {
		*/
		for i := 0; i < len(stockList); i += 10 {
			list := stockList[i : i+10]
			limitCh <- true

			if isStop {
				break
			}

			go func(list []model.StockSymbol) {
				defer func() {
					<-limitCh
				}()

				stockLogList, err := GetStockLogFromYahoo(list)
				if err != nil {
					stockLogger.Println(err)
				} else {
					for _, stockLog := range stockLogList {
						stockdb.InsertStockLog(stockLog)
					}
				}

			}(list)

		}
		//		}

	}
}

func tradeStock(userEmail, tradeType string, symbol model.StockSymbol, cnt uint) bool {
	if tradeType != "ask" && tradeType != "bid" {
		stockLogger.Println("tradeType parameter value was wrong.")
		return false
	}
	stockLogList, err := GetStockLogFromYahoo([]model.StockSymbol{symbol})
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
	/*
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			io.WriteString(w, "hello, stockbot!")
		})
	*/

	http.HandleFunc("/trade", func(w http.ResponseWriter, r *http.Request) {
		/*
			switch r.Method {
			case "POST":
				if tradeStock("kyc1682@gmail.com", "bid", "005930", 10) {
					io.WriteString(w, "trade success")
				} else {
					io.WriteString(w, "trade fail")
				}
			}
		*/
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

	go CollectStockLog(commandCh)

	commanderForWeb(commandCh)
}
