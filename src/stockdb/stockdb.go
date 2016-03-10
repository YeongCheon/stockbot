package stockdb

import (
	. "model"

	_ "github.com/go-sql-driver/mysql"

	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"runtime"
	"strings"
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
	_, filename, _, _ := runtime.Caller(1)
	currentPath := path.Dir(filename)

	var configPath string
	if strings.HasSuffix(currentPath, "src") {
		configPath = "./dbServerInfo.json"
	} else {
		configPath = "../dbServerInfo.json"
	}
	fmt.Println(path.Dir(filename))
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(file, &dbInfo)
}

func getConn() *sql.DB {
	conn, err := sql.Open("mysql", dbInfo.Username+":"+dbInfo.Password+"@"+dbInfo.Protocol+"("+dbInfo.Address+":"+dbInfo.Port+")/"+dbInfo.Dbname)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return conn
}

func InsertStock(symbol StockSymbol) {
	conn := getConn()
	if conn != nil {
		defer conn.Close()
	}
}

func InsertUser(user User) {
	conn := getConn()
	if conn != nil {
		defer conn.Close()
	}
	statment, err := conn.Prepare("INSERT INTO user(`email`, `name`) VALUES(?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	result, err := statment.Exec(user.Email, user.Name)
	if err != nil {
		log.Fatal(err)
		log.Fatal(result)
	}
}

func SelectUser(email string) User {
	conn := getConn()
	if conn != nil {
		defer conn.Close()
	}

	statment, err := conn.Prepare("SELECT `email`, `name` FROM `user` WHERE email = ?")
	rows, err := statment.Query(email)
	if err != nil {
		log.Fatal(err)
	}

	user := User{}
	for rows.Next() {
		rows.Scan(&user.Email, &user.Name)
	}

	return user
}

func InsertStockLog(stockLog StockLog) int64 {
	conn := getConn()
	if conn != nil {
		defer conn.Close()
	}

	statment, err := conn.Prepare("INSERT INTO stock_log(`stock_code`, `stock_market`, `ask`, `bid`) VALUES(?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	result, err := statment.Exec(stockLog.Code, stockLog.Market, stockLog.Ask, stockLog.Bid)
	if err != nil {
		log.Fatal(err)
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	return lastInsertId
}

func SelectStock() []StockSymbol {
	symbols := []StockSymbol{}

	conn := getConn()
	if conn != nil {
		defer conn.Close()
	}

	statment, err := conn.Prepare("SELECT `code`, `market`, `name` FROM `stock`")
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
