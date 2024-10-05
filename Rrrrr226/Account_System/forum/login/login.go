package login

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"net/http"
)

type User struct {
	gorm.Model
	Username string
	Password string
}

var db *gorm.DB

/*func InitDB() {
	var err error
	db, err = gorm.Open("mysql", "user:password@tcp(127.0.0.1:24306)/main_db?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&User{})
	/*u0 := User{Username: "Admin", Password: "hdu123"}
	result := db.Create(&u0)
	if result.Error != nil {
		fmt.Println(result.Error)
	} else {
		fmt.Printf("User created with ID: %v\n", u0.ID)
	}
}*/

// 检查是否匹配
func findUser(Username string) *User {
	var user User
	db.Where("username=?", Username).First(&user)
	if user.Username != "" {
		return &user
	}
	return nil
}

func authenticate(Username, Password string) bool {
	user := findUser(Username)
	return user != nil && user.Password == Password
}

func LoginHandler(c *gin.Context) {
	if c.Request.Method != "POST" {
		c.AbortWithStatusJSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
		return
	}
	username := c.PostForm("Username")
	password := c.PostForm("Password")
	if username == "" || password == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Username or password is empty"})
		return
	}
	if authenticate(username, password) {
		//http.Redirect(w, r, "/forum/:username", http.StatusFound)
		_, _ = fmt.Println("Welcome to forum!")
		//c.Redirect(http.StatusFound, "/forum/"+username)
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
	}
}

/*func sayHello(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadFile("./Hello.txt")
	if err != nil {
		//成功
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, string(b))
}*/

func Registerhandler(c *gin.Context) {
	if c.Request.Method != "POST" {
		c.AbortWithStatusJSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
		return
	}
	username := c.PostForm("Username")
	password := c.PostForm("Password")
	if username == "" || password == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Username or password is empty"})
		return
	}
	if findUser(username) != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
		return
	}

	newuser := User{Username: username, Password: password}
	result := db.Create(&newuser)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully", "username": newuser.Username})
}
