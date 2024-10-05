package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dsn := "username:password@tcp(localhost:3306)/"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	createDBSQL := "CREATE DATABASE IF NOT EXISTS users"
	_, err = db.Exec(createDBSQL)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Database created successfully")
}
