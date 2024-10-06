package main

import (
	"Akuma/database1"
	"Akuma/database2"
	"Akuma/problem"
	"Akuma/register"
	"Akuma/router"
	"log"
	"os"
)

func main() {
	database1.InitDB()
	database1.AutoMigrate(&register.User{})
	database2.InitDB()
	database2.AutoMigrate(&problem.Problem{}, &problem.Submission{})

	r := router.InitRouter()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on port %s", port)
	err := r.Run(":" + port)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
