package main

import (
	stockdb "stockbot/db"
	"stockbot/model"

	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func GetCrawlUrlList() []string {
	result := []string{}
	for _, symbol := range stockdb.SelectStock() {
		s := "s=" + symbol.Code + ".KS"
		f := "f=" + "snabt1" // format: symbol, name, ask, bid, last trade time
		url := "http://finance.yahoo.com/d/quotes.csv?" + s + "&" + f
		result = append(result, url)
	}

	return result
}

func GetStockInfo(stockCode string) model.StockLog {
	s := "s=" + stockCode + ".KS"
	f := "f=" + "snabt1" // format: symbol, name, ask, bid, last trade time
	url := "http://finance.yahoo.com/d/quotes.csv?" + s + "&" + f

	res, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}

	result, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}

	parseResults := ParseYahooCSV(result)

	return parseResults[0]
}

func ParseYahooCSV(target []byte) (logs []model.StockLog) {
	reader := csv.NewReader(strings.NewReader(string(target)))
	for {
		stockLog := model.StockLog{}

		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
		}

		stockLog.Code = record[0]
		ask, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			log.Println(err)
		}
		bid, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			log.Println(err)
		}

		stockLog.Ask = ask
		stockLog.Bid = bid

		logs = append(logs, stockLog)
		//go stockdb.InsertStockLog(stockLog)
	}

	return logs
}

func CollectStockLog(controller chan string) {
	const MARKETSTARTTIME int = 9
	const MARKETENDTIME int = 15

	const maxGoroutineCnt int = 100
	limitCh := make(chan bool, maxGoroutineCnt)

	/*
		timezone, err := time.LoadLocation("Asia/Seoul")
		if err != nil {
			log.Println(err)
		}
	*/

	urls := GetCrawlUrlList()

	for {

		for _, url := range urls {
			now := time.Now()

			if now.Weekday() == 0 || now.Weekday() == 6 {
				//fmt.Println("today is weekend")
			} else if now.Hour() < MARKETSTARTTIME || now.Hour() >= MARKETENDTIME {
				//fmt.Println("market was closed")
			} else {
				limitCh <- true

				go func() {
					defer func() {
						<-limitCh
					}()

					res, err := http.Get(url)
					if err != nil {
						log.Println(err)
					}

					result, err := ioutil.ReadAll(res.Body)
					if err != nil {
						log.Println(err)
					}

					stockLogs := ParseYahooCSV(result)
					for _, stockLog := range stockLogs {
						stockdb.InsertStockLog(stockLog)
					}
				}()

			}
			isStop := false
			select {
			case param := <-controller:
				if param == "stop" {
					fmt.Println("break")
					isStop = true
				}
			default:
			}

			if isStop {
				break
			}
		}
	}
}

func tradeStock(userEmail, tradeType, code string, cnt uint) bool {
	if tradeType != "ask" && tradeType != "bid" {
		log.Println("tradeType parameter value was wrong.")
		return false
	}
	symbol := GetStockInfo(code)

	tradeLog := model.TradeLog{
		UserEmail:   "kyc1682@gmail.com",
		StockCode:   symbol.Code,
		StockMarket: symbol.Market,
		TradeType:   tradeType,
		Price:       symbol.Ask,
		Count:       cnt,
	}
	_, err := stockdb.InsertTradeLog(tradeLog)
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
		switch r.Method {
		case "POST":
			if tradeStock("kyc1682@gmail.com", "bid", "005930", 10) {
				io.WriteString(w, "trade success")
			} else {
				io.WriteString(w, "trade fail")
			}

		}
	})

	http.HandleFunc("/exit", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "exit")
		os.Exit(1)
	})

	log.Println(http.ListenAndServe(":9999", nil))
}

func main() {
	var commandCh chan string
	commandCh = make(chan string)

	go CollectStockLog(commandCh)

	commanderForWeb(commandCh)
}
