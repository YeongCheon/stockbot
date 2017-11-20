package db

import (
	"stockbot/model"

	"testing"
)

var testUser model.User

func init() {
	testUser = model.User{
		Email: "kyc1682@gmail.com",
		Name:  "kyc",
	}
}

func Test_InsertUser(t *testing.T) {
	err := DeleteUser(testUser)
	if err != nil {
		t.Fatal(err)
	}

	err = InsertUser(testUser)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_SelectUser(t *testing.T) {
	email := "kyc1682@gmail.com"
	resultUser := SelectUser(email)
	if resultUser == nil {
		t.Fatalf("resultUser is nil")
	} else if email != resultUser.Email {
		t.Fatalf("parameter email is %s, but result email is %s", email, resultUser.Email)
	}
}
