package model

type DbServerInfo struct {
	Username string `json:username`
	Password string `json:passowrd`
	Protocol string `json:protocol`
	Address  string `json:address`
	Port     string `json:port`
	Dbname   string `json:dbname`
}
