package main

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

//Question问题模型
type Question struct {
    ID int
	UserID int
	Title string
	Answer string
}

//提问函数
func askHandler(db *sql.DB) gin.HandlerFunc {

	return func(c *gin.Context) {
	    sessionID, err := c.Cookie("session_id")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
			return
		}

		user, ok := sessionMap[sessionID]
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
			return
		}

		var req struct {
			Title string `json:"title" binding:"required"`
		}
		if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
			return
		}

		_, err = db.Exec("INSERT INTO questions (user_id, title) VALUES (?, ?)", user.ID, req.Title)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库错误"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "问题已提交"})
	}
}

//修改问题函数
func updateHandler(db *sql.DB) gin.HandlerFunc {
    
	return func(c *gin.Context) {
        sessionID, err := c.Cookie("session_id")
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
            return
        }

        user, ok := sessionMap[sessionID]
        if !ok {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
            return
        }

        var req struct {
            QuestionID int    `json:"question_id" binding:"required"`
			NewTitle string `json:"new_title" binding:"required"`
        }
        if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
            return
        }

		_, err = db.Exec("UPDATE questions SET title = ? WHERE id = ? AND user_id = ?", req.NewTitle, req.QuestionID, user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库错误"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "问题已更新"})
    }
}

//回答问题函数
func answerHandler(db *sql.DB) gin.HandlerFunc {

	return func(c *gin.Context) {
        sessionID, err := c.Cookie("session_id")
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
            return
        }

        user, ok := sessionMap[sessionID]
        if !ok {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
            return
        }

		var req struct {
			QuestionID int    `json:"question_id" binding:"required"`
			Answer string `json:"answer" binding:"required"`
		}

		if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		}

		_, err = db.Exec("UPDATE questions SET answer = ? WHERE id = ? AND user_id = ?", req.Answer, req.QuestionID, user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库错误"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "回答已提交"})
	}
}

//搜索问题函数
func searchHandler(db *sql.DB) gin.HandlerFunc {

	return func(c *gin.Context) {
		sessionID, err := c.Cookie("session_id")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
			return
		}

		_, ok := sessionMap[sessionID]
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
			return
		}

		var req struct {
		    Kewword string `json:"keyword" binding:"required"`
		}

		if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
			return
		}

		rows, err := db.Query("SELECT * FROM questions WHERE title LIKE ?", "%"+req.Kewword+"%")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库错误"})
		}

		defer rows.Close()

		var questions []Question
		for rows.Next() {
		    var q Question
			err := rows.Scan(&q.ID, &q.UserID, &q.Title, &q.Answer)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库错误"})
				return
			}

			questions = append(questions, q)
		}

		c.JSON(http.StatusOK, gin.H{"questions": questions})
	}
}