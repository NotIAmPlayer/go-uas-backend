package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func test() {
	password := "asdfghjkl"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return
	}

	fmt.Println(password + " " + string(hashedPassword))
}
