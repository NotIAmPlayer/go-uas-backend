package config

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func ConnectDatabase() {
	var err error
	DB, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/rapat")

	if err != nil {
		log.Fatal(err)
	}
}
