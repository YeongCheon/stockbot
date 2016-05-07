package db

import (
	"testing"
)

func TestSelectUser(t *testing.T) {
	email := "kyc1682@gmail.com"
	resultUser := SelectUser(email)
	if email != resultUser.Email {
		t.Fatalf("parameter email is %s, but result email is %s", email, resultUser.Email)
	}
}

func TestSelectStock(t *testing.T) {
	result := SelectStock()
}
