package main

import (
	_ "github.com/go-sql-driver/mysql"
	"hduhelp_text/router"
)

func main() {
	router.Run()
}
