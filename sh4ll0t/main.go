package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
	"time"
)

type Answer struct {
	AnswerID   int    `json:"answer_id"`
	AnswerText string `json:"answer_text"`
	Respondent string `json:"respondent"`
	LikesCount int    `json:"likes_count"`
}

type Question struct {
	ID           int      `json:"id"`
	QuestionText string   `json:"question_text"`
	TotalLikes   int      `json:"total_likes"`
	Answers      []Answer `json:"answers"`
	Questioner   string   `json:"questioner"`
}

var (
	store = cookie.NewStore([]byte("shallot"))
)

var (
	Db     *sql.DB
	err    error
	dbHost = "127.0.0.1"
	dbPort = "8888"
	dbUser = "root"
	dbName = "users"
	dbPwd  = "123456"
	sqlStr = dbUser + ":" + dbPwd + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName
)

var question string
var id int

func init() {
	Db, err = sql.Open("mysql", sqlStr)
	if err != nil {
		panic(err)
	}
	err = Db.Ping()
	if err != nil {
		panic("数据库未连接成功: " + err.Error())
	}

	Db.SetConnMaxLifetime(time.Minute * 3)
	Db.SetMaxOpenConns(10)
	Db.SetMaxIdleConns(10)
}
func register(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码处理失败"})
		return
	}

	_, err = Db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, hashedPassword)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			c.JSON(http.StatusConflict, gin.H{"error": "用户名已存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "注册成功"})
}

func ask(c *gin.Context) {
	session, _ := store.Get(c.Request, "session-name")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	username := session.Values["username"].(string)
	question := c.PostForm("question")
	_, err := Db.Exec("INSERT INTO questions (question_text, questioner) VALUES (?, ?)", question, username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var id int
	err = Db.QueryRow("SELECT id FROM questions ORDER BY id DESC LIMIT 1").Scan(&id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	answer_text, err := generateAIAnswer(question)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法生成答案"})
		return
	}
	_, err = Db.Exec("INSERT INTO answers (question_id, answer_text,respondent) VALUES (?, ?,?)", id, answer_text, "ai")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "提问成功"})
}

func generateAIAnswer(question string) (string, error) {
	client := arkruntime.NewClientWithApiKey(
		"ec989009-e9d3-4600-bb9a-68457a8f5e1b",
	)

	ctx := context.Background()

	fmt.Println("----- standard request -----")
	req := model.ChatCompletionRequest{
		Model: "ep-20241003011435-spbz4",
		Messages: []*model.ChatCompletionMessage{
			{
				Role: model.ChatMessageRoleSystem,
				Content: &model.ChatCompletionMessageContent{
					StringValue: volcengine.String("请根据以下问题生成回答："),
				},
			},
			{
				Role: model.ChatMessageRoleUser,
				Content: &model.ChatCompletionMessageContent{
					StringValue: volcengine.String(question),
				},
			},
		},
	}

	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		fmt.Printf("standard chat error: %v\n", err)
		return "", err
	}
	answer := *resp.Choices[0].Message.Content.StringValue
	return answer, nil
}
func ShowQuestionAndAnswer(c *gin.Context) {
	session, _ := store.Get(c.Request, "session-name")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	query := `
		SELECT q.id, q.question_text, q.likes_count, a.id, a.answer_text, q.questioner, a.respondent, a.likes_count
		FROM questions q
		LEFT JOIN answers a ON q.id = a.question_id AND a.check=1
		WHERE q.check = 1 
		ORDER BY q.id
	`
	rows, err := Db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	questionsMap := make(map[int]*Question)
	for rows.Next() {
		var id, answerID, likesCount, answerLikesCount int
		var questionText, answerText, questioner, respondent string

		// 修改 Scan 接收所有字段
		if err := rows.Scan(&id, &questionText, &likesCount, &answerID, &answerText, &questioner, &respondent, &answerLikesCount); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if _, exists := questionsMap[id]; !exists {
			questionsMap[id] = &Question{
				ID:           id,
				QuestionText: questionText,
				TotalLikes:   likesCount, // 从数据库中获取的总点赞数
				Answers:      []Answer{},
				Questioner:   questioner,
			}
		}

		if answerText != "" {
			questionsMap[id].Answers = append(questionsMap[id].Answers, Answer{
				AnswerID:   answerID,
				AnswerText: answerText,
				Respondent: respondent,
				LikesCount: answerLikesCount, // 使用从数据库中获取的答案点赞数
			})
		}
	}

	questionsList := make([]Question, 0, len(questionsMap))
	for _, question := range questionsMap {
		questionsList = append(questionsList, *question)
	}

	c.JSON(http.StatusOK, questionsList)
}

func deleteQuestion(c *gin.Context) {
	session, _ := store.Get(c.Request, "session-name")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	var username string
	id := c.PostForm("id")
	err := Db.QueryRow("SELECT questioner FROM questions WHERE id=?", id).Scan(&username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if session.Values["username"] != username {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "只有本作者才可以删除！"})
		return
	}
	_, err = Db.Exec("DELETE FROM answers WHERE question_id=?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	_, err = Db.Exec("DELETE FROM questions WHERE id=?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

func deleteAnswer(c *gin.Context) {
	session, _ := store.Get(c.Request, "session-name")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	id := c.PostForm("id")
	var username string

	err := Db.QueryRow("SELECT respondent FROM answers WHERE id=?", id).Scan(&username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if session.Values["username"] != username {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "只有本作者才可以删除！"})
		return
	}

	_, err = Db.Exec("DELETE FROM answers WHERE id=?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

func changeQuestion(c *gin.Context) {
	session, _ := store.Get(c.Request, "session-name")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	id := c.PostForm("id")
	question := c.PostForm("question")
	var username string
	err := Db.QueryRow("SELECT questioner FROM questions WHERE id=?", id).Scan(&username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if session.Values["username"] != username {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "只有本作者才可以修改！"})
		return
	}

	_, err = Db.Exec("UPDATE questions SET question_text=? WHERE id=?", question, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "修改成功"})
}
func changeAnswer(c *gin.Context) {
	session, _ := store.Get(c.Request, "session-name")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	id := c.PostForm("id")
	answer := c.PostForm("answer")
	var username string

	err := Db.QueryRow("SELECT respondent FROM answers WHERE id=?", id).Scan(&username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if session.Values["username"] != username {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "只有本作者才可以修改！"})
		return
	}
	_, err = Db.Exec("UPDATE answers SET answer_text=? WHERE id=?", answer, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "修改成功"})
}

func Like(c *gin.Context) {
	session, _ := store.Get(c.Request, "session-name")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	id := c.PostForm("id")
	var count int
	var totalLikes int
	var questionID int
	err := Db.QueryRow("SELECT likes_count FROM answers WHERE id=?", id).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	count++
	_, err = Db.Exec("UPDATE answers SET likes_count=? WHERE id=?", count, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = Db.QueryRow("SELECT question_id FROM answers WHERE id=?", id).Scan(&questionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = Db.QueryRow("SELECT SUM(likes_count) FROM answers WHERE question_id=?", questionID).Scan(&totalLikes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	_, err = Db.Exec("UPDATE questions SET likes_count=? WHERE id=?", totalLikes, questionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"likes_count": count})
}
func like_sort(c *gin.Context) {
	session, _ := store.Get(c.Request, "session-name")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	query := `
		SELECT q.id, q.question_text, q.likes_count, 
		       a.id AS answer_id, a.answer_text, a.respondent, a.likes_count AS answer_likes_count
		FROM questions q
		LEFT JOIN answers a ON q.id = a.question_id AND a.check = 1 
		WHERE q.check = 1 
		ORDER BY q.likes_count DESC, answer_likes_count DESC
	`

	rows, err := Db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	type Answer struct {
		AnswerID   *int
		AnswerText *string
		Respondent *string
		LikesCount *int
	}

	type Question struct {
		ID           int
		QuestionText string
		TotalLikes   int
		Answers      []Answer
	}

	questionsMap := make(map[int]*Question)

	for rows.Next() {
		var qID int
		var questionText string
		var totalLikes int
		var aID *int
		var answerText *string
		var respondent *string
		var answerLikesCount *int

		if err := rows.Scan(&qID, &questionText, &totalLikes, &aID, &answerText, &respondent, &answerLikesCount); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if _, exists := questionsMap[qID]; !exists {
			questionsMap[qID] = &Question{
				ID:           qID,
				QuestionText: questionText,
				TotalLikes:   totalLikes,
				Answers:      []Answer{},
			}
		}

		if aID != nil && answerText != nil {
			questionsMap[qID].Answers = append(questionsMap[qID].Answers, Answer{
				AnswerID:   aID,
				AnswerText: answerText,
				Respondent: respondent,
				LikesCount: answerLikesCount,
			})
		}
	}

	questionsList := make([]Question, 0, len(questionsMap))
	for _, question := range questionsMap {
		questionsList = append(questionsList, *question)
	}

	c.JSON(http.StatusOK, questionsList)
}

func checkAnswer(c *gin.Context) {
	session, _ := store.Get(c.Request, "session-name")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	if session.Values["username"] != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "管理员才可以访问！"})
		return
	}
	id := c.PostForm("id")
	checkStatusStr := c.PostForm("check")
	checkStatus, err := strconv.Atoi(checkStatusStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的审核状态"})
		return
	}

	// 检查 ID 是否存在
	var exists bool
	err = Db.QueryRow("SELECT EXISTS(SELECT 1 FROM answers WHERE id=?)", id).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "ID 不存在"})
		return
	}

	_, err = Db.Exec("UPDATE answers SET `check`=? WHERE id=?", checkStatus, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "审核状态更新成功"})
}
func checkQuestion(c *gin.Context) {
	session, _ := store.Get(c.Request, "session-name")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	if session.Values["username"] != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "管理员才可以访问！"})
		return
	}
	checkStatusStr := c.PostForm("check")
	checkStatus, err := strconv.Atoi(checkStatusStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的审核状态"})
		return
	}
	id := c.PostForm("id")
	var exists bool
	err = Db.QueryRow("SELECT EXISTS(SELECT 1 FROM questions WHERE id=?)", id).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "ID 不存在"})
		return
	}

	_, err = Db.Exec("UPDATE questions SET `check`=? WHERE id=?", checkStatus, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "审核状态更新成功"})
}

func admin(c *gin.Context) {
	session, _ := store.Get(c.Request, "session-name")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	if session.Values["username"] != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "管理员才可以访问！"})
		return
	}
	query := `
		SELECT q.id, q.question_text, a.id, a.answer_text, q.questioner, a.respondent, a.likes_count
		FROM questions q
		LEFT JOIN answers a ON q.id = a.question_id
		ORDER BY q.id
	`
	rows, err := Db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	questionsMap := make(map[int]*Question)
	for rows.Next() {
		var id, answerID, likesCount int
		var questionText, answerText, questioner, respondent string

		if err := rows.Scan(&id, &questionText, &answerID, &answerText, &questioner, &respondent, &likesCount); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if _, exists := questionsMap[id]; !exists {
			questionsMap[id] = &Question{
				ID:           id,
				QuestionText: questionText,
				Answers:      []Answer{},
				Questioner:   questioner,
			}
		}
		if answerText != "" {
			questionsMap[id].Answers = append(questionsMap[id].Answers, Answer{
				AnswerID:   answerID,
				AnswerText: answerText,
				Respondent: respondent,
				LikesCount: likesCount,
			})
		}
	}

	questionsList := make([]Question, 0, len(questionsMap))
	for _, question := range questionsMap {
		questionsList = append(questionsList, *question)
	}

	c.JSON(http.StatusOK, questionsList)

}
func answer(c *gin.Context) {
	session, _ := store.Get(c.Request, "session-name")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	respondent := session.Values["username"].(string)
	id := c.PostForm("id")
	answer := c.PostForm("answer")

	_, err := Db.Exec("INSERT INTO answers (question_id, answer_text, respondent) VALUES (?, ?, ?)", id, answer, respondent)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "回答失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "回答成功"})
}

func login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	var storedHash string
	err := Db.QueryRow("SELECT password FROM users WHERE username=?", username).Scan(&storedHash)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	session, _ := store.Get(c.Request, "session-name")
	session.Values["username"] = username
	session.Values["authenticated"] = true
	session.Save(c.Request, c.Writer)
	c.JSON(http.StatusOK, gin.H{"message": "登录成功"})
}

func logout(c *gin.Context) {
	session, _ := store.Get(c.Request, "session-name")
	session.Options.MaxAge = -1
	session.Save(c.Request, c.Writer)
	c.JSON(http.StatusOK, gin.H{"message": "登出成功"})
}
func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/img", "./img")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{})
	})
	r.POST("/login", login)
	r.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", gin.H{})
	})
	r.POST("/register", register)
	r.GET("/logout", logout)
	r.POST("/ask", ask)
	r.POST("/answer", answer)
	r.GET("/show", ShowQuestionAndAnswer)
	r.POST("deleteAnswer", deleteAnswer)
	r.POST("/deleteQuestion", deleteQuestion)
	r.POST("/changeQuestion", changeQuestion)
	r.POST("/changeAnswer", changeAnswer)
	r.POST("like", Like)
	r.GET("like_sort", like_sort)
	r.GET("admin", admin)
	r.POST("checkAnswer", checkAnswer)
	r.POST("checkQuestion", checkQuestion)
	r.Run(":8000")
}
