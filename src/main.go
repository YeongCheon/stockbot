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

func ParseYahooCSV(target string) {
	reader := csv.NewReader(strings.NewReader(target))
	for {
		history := model.StockLog{}

		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		history.Code = record[0]
		history.Ask = record[2]
		history.Bid = record[3]

		go stockdb.InsertStockLog(history)

		fmt.Println(record)
	}
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
		//fmt.Println(string(result))

		time.Sleep(time.Minute * 1)
	}
}

func main() {
	CollectStockLog()
}
