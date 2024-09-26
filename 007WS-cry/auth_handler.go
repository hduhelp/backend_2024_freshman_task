package main

import(
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

var sessionMap = make(map[string]*User)

type User struct {
    ID int
	Username string
    Password string
}

func registerHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

	    var req struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
			return
		}

		// TODO: 检查用户名是否已存在
		var existingUser User
		err := db.QueryRow("SELECT id FROM users WHERE username = ?", req.Username).Scan(&existingUser.ID)
		if err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "用户名已存在"})
			return
		}

		// TODO: 创建新用户
		_, err = db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", req.Username, req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建用户失败"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "用户创建成功"})
	}
}

func loginHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

	    var req struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
			return
		}

		// TODO: 验证用户名和密码
		var user User
		err := db.QueryRow("SELECT id FROM users WHERE username = ? AND password = ?", req.Username, req.Password).Scan(&user.ID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		}

		//session建立
		sessionID := "session_" + req.Username
		sessionMap[sessionID] = &user
		c.SetCookie("session_id", sessionID, 3600, "/", "localhost", false, true)

		c.JSON(http.StatusOK, gin.H{"message": "登录成功"})
	}
}