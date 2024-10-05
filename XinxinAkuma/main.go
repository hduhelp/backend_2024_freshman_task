package main

import (
	"Akuma/database1"
	"Akuma/database2"
	"Akuma/problem"
	"Akuma/register"
	"Akuma/router"
)

func main() {
	database1.InitDB()
	database1.AutoMigrate(&register.User{})
	database2.InitDB()
	database2.AutoMigrate(&problem.Problem{}, &problem.Submission{})

	r := router.InitRouter()
	err := r.Run(":8080")
	if err != nil {
		return
	}
}
