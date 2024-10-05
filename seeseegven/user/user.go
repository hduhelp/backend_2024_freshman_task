package user

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "unicode"
)

type User struct {
    ID     string `json:"ID" binding:"required,min=8,max=8"`
    Psd    string `json:"psd" binding:"required,min=8"`
    Type   string `json:"type"` // "login" or "register"
    Islogin bool   `json:"islogin"` // 登录状态
}

var db = make(map[string]string) // 临时的“数据库”
var LoginStatus = make(map[string]bool) // 登录状态的“数据库”

func isUserRegistered(ID string) bool {
    _, exists := db[ID]
    return exists
}

func registerUser(ID, psd string) {
    db[ID] = psd
    LoginStatus[ID] = false // 注册时设置登录状态为false
}

func containsUpper(s string) bool {
    for _, v := range s {
        if unicode.IsUpper(v) {
            return true
        }
    }
    return false
}

func containsLower(s string) bool {
    for _, v := range s {
        if unicode.IsLower(v) {
            return true
        }
    }
    return false
}

func containsDigit(s string) bool {
    for _, v := range s {
        if unicode.IsDigit(v) {
            return true
        }
    }
    return false
}

func HandleLoginRegister(c *gin.Context) {
    var user User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if user.Type == "register" {
        if isUserRegistered(user.ID) {
            c.JSON(http.StatusBadRequest, gin.H{"error": "用户已存在"})
            return
        }
        if len(user.ID) != 8 {
            c.JSON(http.StatusBadRequest, gin.H{"error": "用户名必须是8位学工号"})
            return
        }
        if len(user.Psd) < 8 || !containsUpper(user.Psd) || !containsLower(user.Psd) || !containsDigit(user.Psd) {
            c.JSON(http.StatusBadRequest, gin.H{"error": "密码至少8位，包括大小写字母和数字"})
            return
        }
        registerUser(user.ID, user.Psd)
        c.JSON(http.StatusOK, gin.H{"message": "注册成功", "greeting": "你好，" + user.ID + "请登录。", "islogin": false})
    } else if user.Type == "login" {
        if isUserRegistered(user.ID) && db[user.ID] == user.Psd {
            LoginStatus[user.ID] = true // 登录成功设置登录状态为true
            c.SetCookie("userID", user.ID, 3600, "/", "", false, true) // 设置cookie
            c.JSON(http.StatusOK, gin.H{"message": "登录成功", "greeting": "你好，" + user.ID + "!", "islogin": true})
        } else {
            LoginStatus[user.ID] = false // 登录失败设置登录状态为false
            c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误", "islogin": false})
        }
    } else if user.Type == "logout" {
            userID := c.MustGet("userID").(string)
            if _, exists := db[userID]; exists {
                LoginStatus[userID] = false // 设置登录状态为false
                c.JSON(http.StatusOK, gin.H{"message": "登出成功", "islogin": false})
                c.SetCookie("userID", "", -1, "/", "", false, true) // 清除cookie
            }else{
                c.JSON(http.StatusBadRequest, gin.H{"error": "用户未登录"})
            }
    }else {
        c.JSON(http.StatusBadRequest, gin.H{"error": "请求类型未指定或无效"})
    }
}