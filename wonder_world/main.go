package main

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type User struct { //用户
	gorm.Model `json:"gorm_._model"`
	Name       string `json:"name,omitempty"`
	Telephone  string `json:"telephone,omitempty"`
	Password   string `json:"password,omitempty"`
	ID         int    `json:"id,omitempty"`
}
type Ques struct { //问题
	gorm.Model
	Name  string
	Title string
	Put   string
	Key   string
}
type Anse struct { //我、回答
	Name  string `json:"name"`
	Text  string `json:"text"`
	Key   string `json:"key"`
	Title string `json:"title"`
	gorm.Model
}

var (
	db  *gorm.DB
	err error
)

func main() {
	dsn := "root:123789@tcp(127.0.0.1:3306)/demo?charset=utf8mb4&parseTime=True&loc=Local" //数据库登入
	db, _ = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	_ = db.AutoMigrate(&User{}, &Ques{}, &Anse{})
	r := gin.Default()
	v := r.Group("/v1")
	{
		v.POST("/register", register) //注册
		v.POST("/login", login)       //登入
	}
	q := r.Group("/problem_hall") //问题
	{
		q.POST("/put", question) //提问
		q.POST("delete", de)     //删除
		q.GET("/hall", hall)     //大厅（输出数据）
	}
	a := r.Group("/answer") //回答
	{
		a.POST("/put", anput)      //输入
		a.POST("/delete", andelet) //删除
		a.GET("/hall", anhall)     //大厅
	}
	r.Run(":8080")
}
func register(c *gin.Context) {
	name := c.PostForm("name")
	telephone := c.PostForm("telephone")
	password := c.PostForm("password")
	if len(name) == 0 {
		c.JSON(422, gin.H{
			"code":    422,
			"message": "用户名不能为空",
		})
		return
	}
	if len(telephone) != 11 {
		c.JSON(422, gin.H{
			"code":    422,
			"message": "手机号必须为11位",
		})
		return
	}
	if len(password) <= 6 {
		c.JSON(422, gin.H{
			"code":    422,
			"message": "密码不能少于6位",
		})
		return
	}
	var user User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 {
		c.JSON(422, gin.H{
			"code":    422,
			"message": "用户已存在",
		})
		return
	}
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    500,
			"message": "密码加密错误",
		})
		return
	}
	newUser := User{
		Name:      name,
		Telephone: telephone,
		Password:  string(hasedPassword),
	}
	db.Create(&newUser)
	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
	})

}
func login(c *gin.Context) {
	telephone := c.PostForm("telephone")
	password := c.PostForm("password")
	var user User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID == 0 {
		c.JSON(422, gin.H{
			"code":    422,
			"message": "用户不存在",
		})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		c.JSON(422, gin.H{
			"code":    422,
			"message": "密码错误",
		})
		return
	}
	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
	})
}
func question(c *gin.Context) {
	title := c.PostForm("title")
	description := c.PostForm("description")
	name := c.PostForm("name")
	key := c.PostForm("key")
	var use Ques
	db.Where("title = ?", title).First(&use)
	if use.ID != 0 {
		c.JSON(422, gin.H{
			"code":    422,
			"message": "问题标题重复",
		})
		return
	}
	var user User
	db.Where("name = ?", name).First(&user)
	if user.ID == 0 {
		c.JSON(422, gin.H{
			"code":    422,
			"message": "请输入正确的用户名",
		})
		return
	}
	newQuestion := Ques{
		Name:  name,
		Title: title,
		Put:   description,
		Key:   key,
	}
	db.Create(&newQuestion)
	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
	})

}
func de(c *gin.Context) {
	title := c.PostForm("title")
	key := c.PostForm("key")
	var use Ques
	db.Where("title = ?", title).First(&use)
	if use.ID == 0 {
		c.JSON(422, gin.H{
			"code":    422,
			"message": "问题不存在",
		})
		return
	}
	if use.ID != 0 {
		if key != use.Key {
			c.JSON(422, gin.H{
				"code":    422,
				"message": "密匙错误",
			})
			return
		}
		db.Delete(&use)
		c.JSON(200, gin.H{
			"code":    200,
			"message": "success",
		})
	}
}
func hall(c *gin.Context) {
	number := c.PostForm("first")
	number1 := c.PostForm("number")
	var num int
	num, _ = strconv.Atoi(number)
	num1, _ := strconv.Atoi(number1)
	var user []Ques
	db.Limit(num1).Offset(num).Find(&user) //num--编号num1--数量
	for _, qusetion := range user {
		c.JSON(200, gin.H{
			"title": qusetion.Title,
			"put":   qusetion.Put,
			"name":  qusetion.Name,
		})
	}
}
func anput(c *gin.Context) {
	description := c.PostForm("description")
	name := c.PostForm("name")
	key := c.PostForm("key")
	title := c.PostForm("title")
	var user User
	db.Where("name = ?", name).First(&user)
	if user.ID == 0 {
		c.JSON(422, gin.H{
			"code":    422,
			"message": "请输入正确的用户名",
		})
		return
	}
	var use Ques
	db.Where("title = ?", title).First(&use)
	if use.ID == 0 {
		c.JSON(422, gin.H{
			"code":    422,
			"message": "问题不存在",
		})
		return
	}
	newAnswer := Anse{
		Name:  name,
		Text:  description,
		Key:   key,
		Title: title,
	}
	db.Create(&newAnswer)
	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
	})

}
func andelet(c *gin.Context) {
	f := 1
	key := c.PostForm("key")
	title := c.PostForm("title")
	text := c.PostForm("text")
	name := c.PostForm("name")
	var user Ques
	db.Where("title = ?", title).First(&user)
	if user.ID == 0 {
		c.JSON(422, gin.H{
			"code":    422,
			"message": "问题不存在",
		})
		return
	}
	var use []Anse
	db.Where("name = ?", name).First(&use)
	if user.ID != 0 {
		for _, answer := range use {
			if answer.Text == text {
				f = 0
				if answer.Key == key {
					db.Delete(&answer)
				}
				if answer.Key != key {
					c.JSON(422, gin.H{
						"code":    422,
						"message": "密匙错误",
					})
					return
				}
				break
			}
			if f == 1 {
				c.JSON(422, gin.H{
					"code":    422,
					"message": "回答不存在",
				})
				return
			}
		}
		c.JSON(200, gin.H{
			"code":    200,
			"message": "success",
		})
	}
}
func anhall(c *gin.Context) {
	number := c.PostForm("first")
	number1 := c.PostForm("number")
	title := c.PostForm("title")
	var num, _ = strconv.Atoi(number)
	num1, _ := strconv.Atoi(number1)
	var user []Anse
	db.Limit(num1).Offset(num).Find(&user)
	for _, answer := range user {
		if answer.Title == title {
		}
		c.JSON(200, gin.H{
			"name": answer.Name,
			"text": answer.Text,
		})
	}
}
