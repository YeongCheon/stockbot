package db

import (
	. "stockbot/model"

	_ "github.com/go-sql-driver/mysql"

	"database/sql"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
)

var dbInfo DbServerInfo

func init() {
	file, err := ioutil.ReadFile("./dbServerInfo.json")
	if err != nil {
		log.Println(err)
	}

	json.Unmarshal(file, &dbInfo)
}

func getConn() *sql.DB {
	conn, err := sql.Open("mysql", dbInfo.Username+":"+dbInfo.Password+"@"+dbInfo.Protocol+"("+dbInfo.Address+":"+dbInfo.Port+")/"+dbInfo.Dbname)
	if err != nil {
		log.Println("connect error : ", err)
		return nil
	}

	conn.SetMaxOpenConns(10000)

	return conn
}

func InsertStock(symbol StockSymbol) {
	conn := getConn()
	if conn != nil {
		defer conn.Close()
	}
}

func InsertUser(user User) error {
	conn := getConn()
	if conn == nil {
		return errors.New("connection is null")
	}
	defer conn.Close()

	statment, err := conn.Prepare("INSERT INTO member(`email`, `name`) VALUES(?, ?)")
	if err != nil {
		log.Println(err)
		return err
	}

	result, err := statment.Exec(user.Email, user.Name)
	if err != nil {
		log.Println(err)
		log.Println(result)
		return err
	}

	return nil
}

func DeleteUser(user User) error {
	conn := getConn()
	if conn == nil {
		return errors.New("connection is null")
	}
	defer conn.Close()

	statment, err := conn.Prepare("DELETE FROM member WHERE email = ?")
	if err != nil {
		log.Println(err)
		return err
	}

	result, err := statment.Exec(user.Email)
	if err != nil {
		log.Println(err)
		log.Println(result)
		return err
	}

	return nil
}

func SelectUser(email string) *User {
	conn := getConn()
	if conn != nil {
		defer conn.Close()
	}

	statment, err := conn.Prepare("SELECT `email`, `name` FROM `member` WHERE email = ?")
	if err != nil {
		log.Println(err)
		return nil
	}
	rows, err := statment.Query(email)
	if err != nil {
		log.Println(err)
		return nil
	}

	user := User{}
	for rows.Next() {
		rows.Scan(&user.Email, &user.Name)
	}

	return &user
}

func InsertStockLog(stockLog StockLog) (int64, error) {
	conn := getConn()
	if conn != nil {
		defer conn.Close()
	}

	statment, err := conn.Prepare("INSERT INTO stock_log(`stock_code`, `stock_market`, `ask`, `bid`) VALUES(?, ?, ?, ?)")
	if err != nil {
		log.Println(err)
	}

	result, err := statment.Exec(stockLog.Code, stockLog.Market, stockLog.Ask, stockLog.Bid)
	if err != nil {
		log.Println(err)
		return -1, errors.New("insert statment execute fail")
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		log.Println(err)
		return 01, errors.New("get last insert id fail")
	}
	return lastInsertId, nil
}

func GetAllStockList() []StockSymbol {
	symbols := []StockSymbol{}

	conn := getConn()
	if conn != nil {
		defer conn.Close()
	}

	statment, err := conn.Prepare("SELECT `code`, `market`, `name` FROM `stock`")
	if err != nil {
		log.Println(err)
	}

	rows, err := statment.Query()
	if err != nil {
		log.Println(err)
	} else {
		for rows.Next() {
			symbol := StockSymbol{}
			rows.Scan(&symbol.Code, &symbol.Market, &symbol.Name)
			symbols = append(symbols, symbol)
		}
	}

	return symbols
}
