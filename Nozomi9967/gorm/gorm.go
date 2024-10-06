package gorm

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"net/http"
	"time"
)

var GLOBAL_DB *gorm.DB

func Gorm() {
	db, _ := gorm.Open(mysql.New(mysql.Config{
		DSN:               "root:123456@tcp(127.0.0.1:3306)/my_database?charset=utf8mb4&parseTime=True&loc=Local",
		DefaultStringSize: 171,
	}), &gorm.Config{
		SkipDefaultTransaction: false,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "t_", // table name prefix, table for `User` would be `t_users`
			SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	if db.Error != nil {
		panic(db.Error)
	}

	sqldb, _ := db.DB()
	sqldb.SetMaxIdleConns(10)
	sqldb.SetMaxOpenConns(100)
	sqldb.SetConnMaxLifetime(time.Hour)

	GLOBAL_DB = db

}

func SelfInfoHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "selfinfo", nil)
}
