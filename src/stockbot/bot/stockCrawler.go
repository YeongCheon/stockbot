package bot

import (
	stockdb "stockbot/db"
	"stockbot/logger"
	"stockbot/model"

	"encoding/csv"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

const APIURL string = "http://finance.yahoo.com/d/quotes.csv"

type StockCrawler struct {
	stockLogger *logger.Logger
}

func (stockCrawler *StockCrawler) GetStockLogFromYahoo(stockSymbolList []model.StockSymbol) ([]model.StockLog, error) {
	var s string
	for _, stockSymbol := range stockSymbolList {
		s += stockSymbol.Code + ".KS+"
	}
	f := "f=" + "snabt1" // format: symbol, name, ask, bid, last trade time
	url := APIURL + "?s=" + s + "&" + f
	res, err := http.Get(url)
	if err != nil {
		stockCrawler.stockLogger.Println(err)
	}

	result, err := ioutil.ReadAll(res.Body)
	if err != nil {
		stockCrawler.stockLogger.Println(err)
	}

	stockLogs := stockCrawler.parseYahooCSV(result)
	if len(stockLogs) <= 0 {
		return nil, errors.New("GetStockLogFromYahoo func's result is nil")
	}

	for idx, _ := range stockLogs {
		stockLogs[idx].Code = stockLogs[idx].Code[:6]
		stockLogs[idx].Market = stockSymbolList[idx].Market
	}

	return stockLogs, nil
}

func (stockCrawler *StockCrawler) parseYahooCSV(target []byte) []model.StockLog {
	var result []model.StockLog
	reader := csv.NewReader(strings.NewReader(string(target)))
	defer func() {
		if r := recover(); r != nil {
			stockCrawler.stockLogger.Println(r)
		}
	}()
	for {
		stockLog := model.StockLog{}

		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			stockCrawler.stockLogger.Println(err)
		}

		stockLog.Code = record[0]
		ask, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			stockCrawler.stockLogger.Println(err)
		}
		bid, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			stockCrawler.stockLogger.Println(err)
		}

		stockLog.Ask = ask
		stockLog.Bid = bid

		result = append(result, stockLog)
	}

	return result
}

func (stockCrawler *StockCrawler) CollectStockData(controller chan string) {
	stockCrawler.stockLogger.Println("start stock log crawling...")
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

				stockLogList, err := stockCrawler.GetStockLogFromYahoo(list)
				if err != nil {
					stockCrawler.stockLogger.Println(err)
				} else {
					for _, stockLog := range stockLogList {
						stockdb.InsertStockLog(stockLog)
					}
				}
			}(list)
		}
	}
}
