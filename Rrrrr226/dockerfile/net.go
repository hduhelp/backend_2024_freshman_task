package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

/*func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello,Golang"))
	})
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		panic(err)
		return
	}
	fmt.Println("Succeed!")
}*/

func main() {
	router := gin.Default()

	router.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"1": "welcome",
		})
		c.String(200, "hello world")
	})

	fmt.Println(router.Run(":8888"))
	err := router.Run(":8888")
	if err != nil {
		fmt.Println(err)
	}
}
