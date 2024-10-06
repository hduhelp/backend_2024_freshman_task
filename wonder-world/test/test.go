package test

import (
	"github.com/gin-gonic/gin"
	db2 "wonder-world/db"
)

type session struct {
	Name  string
	Value string
}

func Getcookie(c *gin.Context, name string) bool {

	f := true
	db := db2.Dbfrom()
	var use session
	cookievalue, err := c.Cookie(name)
	if err != nil {
		f = false
	}
	db.Where("name=?", cookievalue).First(&use)
	if use.Value != cookievalue {
		f = false
	}
	return f
}
