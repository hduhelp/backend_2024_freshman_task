package question

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "strconv"
    "time"
    "myproject/user"
)

type Question struct {
    ID         uint      `json:"id"`
    Title      string    `json:"title" binding:"required"`
    Content    string    `json:"content" binding:"required"`
    Partition  string    `json:"partition" binding:"required"` // "life", "learning"
    CreatedAt  time.Time `json:"created_at"`             // 记录提问时间
    Answers    []Answer  `json:"answers"`
    CreatorID  string    `json:"creator_id"` // 问题的创建者ID
}

type Answer struct {
    ID      uint      `json:"id"`
    Content string    `json:"content" binding:"required"`
    QuestionID uint     `json:"question_id"`
    CreatedAt time.Time `json:"created_at"` // 记录回答时间
}

var db1 = make(map[uint]Question) // 临时的“数据库”
var nextID = 1

func InitializeRoutes(r *gin.Engine) {
    r.POST("/questions", CreateQuestion)
    r.GET("/questions", GetQuestions)
    r.GET("/questions/:id", GetQuestion)
    r.POST("/questions/:id/answers", CreateAnswer)
    r.DELETE("/questions/:id", DeleteQuestion) // 添加删除问题的路由
}

func CreateQuestion(c *gin.Context) {
    userID, err := c.Cookie("userID")
    if err != nil || userID == "" || !user.LoginStatus[userID]{
        c.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录或注册"})
        return
    }
    var question Question
    if err := c.ShouldBindJSON(&question); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    question.CreatedAt = time.Now()
    question.ID = uint(nextID)
    question.CreatorID = userID
    db1[question.ID] = question
    nextID++

    c.JSON(http.StatusOK, gin.H{"message": "提问成功", "question": question})
}

func GetQuestions(c *gin.Context) {
    partition := c.DefaultQuery("partition", "all")
    questions := make([]Question, 0, len(db1))
    for _, question := range db1 {
        if partition == "all" || question.Partition == partition {
            questions = append(questions, question)
        }
    }
    c.JSON(http.StatusOK, questions)
}

func GetQuestion(c *gin.Context) {
    idStr := c.Param("id")
    if idStr == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID去哪了？"})
        return
    }

    id, err := strconv.Atoi(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "无效ID"})
        return
    }

    question, exists := db1[uint(id)]
    if !exists {
        c.JSON(http.StatusNotFound, gin.H{"error": "没找到该问题"})
        return
    }
    c.JSON(http.StatusOK, question)
}

func CreateAnswer(c *gin.Context) {
    userID, err := c.Cookie("userID")
    if err != nil || userID == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录或注册"})
        return
    }

    idStr := c.Param("id") // 从URL参数中获取问题ID
    id, err := strconv.Atoi(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "无效ID"})
        return
    }

    question, exists := db1[uint(id)]
    if !exists {
        c.JSON(http.StatusNotFound, gin.H{"error": "没找到该问题"})
        return
    }

    var answer Answer
    answer.CreatedAt = time.Now()
    if err := c.ShouldBindJSON(&answer); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    answer.QuestionID = uint(id)
    question.Answers = append(question.Answers, answer)
    db1[uint(id)] = question // 更新问题
    c.JSON(http.StatusOK, gin.H{"message": "回答成功，感谢您的帮助", "answer": answer})
}

func DeleteQuestion(c *gin.Context) {
    userID, err := c.Cookie("userID")
    if err != nil || userID == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录或注册"})
        return
    }

    idStr := c.Param("id")
    if idStr == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID去哪了？"})
        return
    }

    id, err := strconv.Atoi(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "无效ID"})
        return
    }

    question, exists := db1[uint(id)]
    if !exists {
        c.JSON(http.StatusNotFound, gin.H{"error": "没找到该问题"})
        return
    }

    // 验证用户身份，确保只有问题的创建者可以删除问题
    if question.CreatorID != userID {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "只有问题的创建者可以删除问题"})
        return
    }

    //删除问题
delete(db1, uint(id))

c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}