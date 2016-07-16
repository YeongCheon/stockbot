package main

import (
	stockdb "stockbot/db"
	"stockbot/model"

	"bufio"
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
		log.Fatal(err)
	}

	result, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
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
			log.Fatal(err)
		}

		stockLog.Code = record[0]
		ask, err := strconv.Atoi(record[2])
		if err != nil {
			log.Fatal(err)
		}
		bid, err := strconv.Atoi(record[3])
		if err != nil {
			log.Fatal(err)
		}

		stockLog.Ask = int64(ask)
		stockLog.Bid = int64(bid)

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
			log.Fatal(err)
		}
	*/

	for {
		urls := GetCrawlUrlList()

		for _, url := range urls {
			now := time.Now()

			if now.Weekday() == 0 || now.Weekday() == 6 {
				fmt.Println("today is weekend")
			} else if now.Hour() < MARKETSTARTTIME || now.Hour() >= MARKETENDTIME {
				fmt.Println("market was closed")
			} else {
				limitCh <- true

				go func() {
					defer func() {
						<-limitCh
					}()

					res, err := http.Get(url)
					if err != nil {
						log.Fatal(err)
					}

					result, err := ioutil.ReadAll(res.Body)
					if err != nil {
						log.Fatal(err)
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

func buyStock(code string, cnt uint) {
	symbol := GetStockInfo(code)

	tradeLog := model.TradeLog{
		UserEmail:   "kyc1682@gmail.com",
		StockCode:   symbol.Code,
		StockMarket: symbol.Market,
		TradeType:   "ask", //매수
		Price:       symbol.Ask,
		Count:       cnt,
	}
	stockdb.InsertTradeLog(tradeLog)
}

func commander(commandCh chan<- string) {
	consoleReader := bufio.NewReader(os.Stdin)

	for {
		inputStr, _ := consoleReader.ReadString('\n')
		inputStr = strings.Replace(inputStr, "\n", "", -1)
		args := strings.Split(inputStr, " ")
		command := args[0]

		switch command {
		case "exit":
			break
		case "start", "stop":
			commandCh <- inputStr
		case "buy":
			code := args[1]
			cnt, err := strconv.Atoi(args[2])
			if err != nil {
				log.Fatal(err)
			}

			buyStock(code, uint(cnt))
		}
	}
}

func main() {
	var commandCh chan string
	commandCh = make(chan string)

	go CollectStockLog(commandCh)

	commander(commandCh)
}
