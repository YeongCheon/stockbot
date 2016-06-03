package model

type DbServerInfo struct {
	Username string `json:username`
	Password string `json:password`
	Protocol string `json:protocol`
	Address  string `json:address`
	Port     string `json:port`
	Dbname   string `json:dbname`
}
