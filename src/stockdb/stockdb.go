package stockdb

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"encoding/json"
	"fmt"
	"io/ioutil"
	. "stockmodel"
)

type DbServerInfo struct {
	Username string `json:username`
	Password string `json:passowrd`
	Protocol string `json:protocol`
	Address  string `json:address`
	Port     string `json:port`
	Dbname   string `json:dbname`
}

var dbInfo DbServerInfo

func init() {
	file, err := ioutil.ReadFile("./dbServerInfo.json")
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(file, &dbInfo)
}

func InsertHistory(history StockHistory) {
	conn, err := sql.Open("mysql", dbInfo.Username+":"+dbInfo.Password+"@"+dbInfo.Protocol+"("+dbInfo.Address+":"+dbInfo.Port+")/"+dbInfo.Dbname)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	statment, err := conn.Prepare("INSERT INTO stockbot_kospi_history(history_code, history_ask, history_bid) VALUES(?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	result, err := statment.Exec(history.Code, history.Ask, history.Bid)
	if err != nil {
		log.Fatal(err)
		log.Fatal(result)
	}
}

func SelectKospiSymbols() []StockSymbol {
	var symbols []StockSymbol
	symbols = []StockSymbol{}

	conn, err := sql.Open("mysql", dbInfo.Username+":"+dbInfo.Password+"@"+dbInfo.Protocol+"("+dbInfo.Address+":"+dbInfo.Port+")/"+dbInfo.Dbname)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer conn.Close()

	statment, err := conn.Prepare("SELECT kospi_code, kospi_name FROM stockbot_kospi")
	if err != nil {
		log.Fatal(err)
	}

	rows, err := statment.Query()
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		symbol := StockSymbol{}
		rows.Scan(&symbol.Code, &symbol.Name)
		symbols = append(symbols, symbol)
	}

	return symbols
}
