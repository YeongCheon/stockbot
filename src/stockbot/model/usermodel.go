package model

type User struct {
	Email string
	Name  string
}

type UserStock struct {
	id        uint
	userEmail string
	stockCode string
}
