package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"bytes"
	"encoding/csv"
	"io"
	"model"
	"stockdb"
	"strings"
)

func GetSymbolParameter() string {
	var buffer bytes.Buffer
	for _, symbol := range stockdb.SelectStock() {
		buffer.WriteString(symbol.Code + ".KS+")
	}

	//return "005930.KS+155900.KS"
	return buffer.String()
}

func ParseYahooCSV(target string) (logs []model.StockLog) {
	reader := csv.NewReader(strings.NewReader(target))
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
		stockLog.Ask = record[2]
		stockLog.Bid = record[3]

		logs = append(logs, stockLog)
		go stockdb.InsertStockLog(stockLog)
	}

	return logs
}

func CollectStockLog() {
	s := "s=" + GetSymbolParameter()
	f := "f=" + "snabt1" // format: symbol, name, ask, bid, last trade time
	url := "http://finance.yahoo.com/d/quotes.csv?" + s + "&" + f
	fmt.Println("url : " + url)
	for {
		res, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		result, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		ParseYahooCSV(string(result))

		time.Sleep(time.Second * 1)
	}
}

func main() {
	CollectStockLog()
}
