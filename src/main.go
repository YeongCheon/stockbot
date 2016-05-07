package main

import (
	stockdb "stockbot/db"
	"stockbot/model"

	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
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
		//go stockdb.InsertStockLog(stockLog)
	}

	return logs
}

func CollectStockLog(controller chan string) {
	const MARKETSTARTTIME int = 9
	const MARKETENDTIME int = 15
	const TICK time.Duration = 1 * time.Second

	s := "s=" + GetSymbolParameter()
	f := "f=" + "snabt1" // format: symbol, name, ask, bid, last trade time
	url := "http://finance.yahoo.com/d/quotes.csv?" + s + "&" + f
	fmt.Println("url : " + url)

	timezone, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		log.Fatal(err)
	}
	for {
		now := time.Now()

		if now.Weekday() == 0 || now.Weekday() == 6 {
			//fmt.Println("today is weekend")
			/*
				var leftDay int
				if 1-now.Weekday() < 0 {
					leftDay = 2
				} else {
					leftDay = 1
				}
				nextStartTime := time.Date(now.Year(), now.Month(), now.Day()+leftDay, 9, 0, 0, 0, timezone)
				time.Sleep(time.Millisecond*time.Duration(nextStartTime.UnixNano()/1000) - time.Duration(now.UnixNano()/1000))
			*/
		} else if now.Hour() < MARKETSTARTTIME || now.Hour() >= MARKETENDTIME {
			fmt.Println("market close")
			nextStartTime := time.Date(now.Year(), now.Month(), now.Day()+1, 9, 0, 0, 0, timezone)
			time.Sleep(time.Millisecond*time.Duration(nextStartTime.UnixNano()/1000) - time.Duration(now.UnixNano()/1000))
		} else {
			res, err := http.Get(url)
			if err != nil {
				log.Fatal(err)
			}
			result, err := ioutil.ReadAll(res.Body)
			if err != nil {
				log.Fatal(err)
			}
			ParseYahooCSV(string(result))
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

		fmt.Println("collecting...")
		time.Sleep(TICK)
	}
}

func main() {
	var controller chan string
	controller = make(chan string)
	go CollectStockLog(controller)

	consoleReader := bufio.NewReader(os.Stdin)
	for {
		tmp, _ := consoleReader.ReadString('\n')
		tmp = strings.Replace(tmp, "\n", "", -1)
		controller <- tmp
		if tmp == "stop" {
			break
		}
	}
}
