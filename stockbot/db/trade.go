package db

import (
	. "stockbot/model"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"
)

func init() {
	file, err := ioutil.ReadFile("./dbServerInfo.json")
	if err != nil {
		log.Println(err)
	}

	json.Unmarshal(file, &dbInfo)
}

func InsertTradeLog(tradeLog TradeLog) (int64, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()

	conn := getConn()
	if conn != nil {
		defer conn.Close()
	}

	statment, err := conn.Prepare("INSERT INTO trade_log(`member_email`, `stock_code`, `stock_market`, `trade_type`, `price`) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		log.Println(err)
		return -1, err
	}

	result, err := statment.Exec(tradeLog.UserEmail, tradeLog.StockCode, tradeLog.StockMarket, tradeLog.TradeType, tradeLog.Price)

	if err != nil {
		log.Println(err)
		return -1, err
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		log.Println(err)
		return -1, err
	}

	return lastInsertId, nil
}

func SelectTradeLog(start, end time.Time) []TradeLog {
	tradeLogs := []TradeLog{}
	conn := getConn()
	if conn != nil {
		defer conn.Close()
	}

	statement, err := conn.Prepare("SELECT * FROM trade_log WHERE `trade_timestamp` BETWEEN ? AND ?")

	rows, err := statement.Query(start, end)
	if err != nil {
		log.Println(err)
	}
	for rows.Next() {
		tradeLog := TradeLog{}
		rows.Scan(&tradeLog.Id, &tradeLog.UserEmail, &tradeLog.StockCode, &tradeLog.StockMarket, &tradeLog.StockMarket, &tradeLog.TradeType, &tradeLog.Price, &tradeLog.Timestamp)
		tradeLogs = append(tradeLogs, tradeLog)
	}

	return tradeLogs
}
