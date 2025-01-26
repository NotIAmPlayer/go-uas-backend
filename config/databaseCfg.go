package config

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func ConnectDatabase() {
	var err error
	DB, err = sql.Open("mysql", "root:@tcp(autorack.proxy.rlwy.net:33157)/rapat")

	if err != nil {
		log.Fatal(err)
	}
}
