package stockdb

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
	if len(result) != 1 {
		t.Fatalf("result size is wrong. your result size is %d", len(result))
	}

	if result[0].Code != "005930" {
		t.Fatalf("your code is wrong. your code is %s", result[0].Code)
	}
}
