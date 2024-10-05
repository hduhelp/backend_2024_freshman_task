package main

import (
	"QASystem/dao"
	"QASystem/router"
	"fmt"
)

func main() {
	isConnect := dao.InitDB()
	if isConnect != nil {
		fmt.Println("数据库连接失败")
	}

	r := router.Router()
	err := r.Run(":8080")
	if err != nil {
		panic(err)
	}
}
