package main

import (
	"log"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 用户结构体
type User struct {
	gorm.Model
	Username  string `gorm:"unique;not null"`
	Password  string `gorm:"not null"`
	IsManager bool
}

type KeyWord struct {
	Keyword string
}

// 问题结构体
type Question struct {
	gorm.Model
	Title        string
	Content      string
	UserID       uint
	Answers      []Answer `gorm:"foreignKey:QuestionID"`
	IsDeleted    bool
	Deleted_Time time.Time
}

type AgreeOrNot struct {
	AgreeNum    uint
	DisagreeNum uint
}

// 答案中的评论结构体
type Comment struct {
	gorm.Model
	AgreeOrNot
	Content  string
	UserID   uint
	AnswerID uint
}

// 回复评论结构体
type ReplyComment struct {
	gorm.Model
	CommentUserID  uint
	CommentID      uint
	ReplyCommentID uint
	ReplyUserID    uint
}

// 答案结构体
type Answer struct {
	gorm.Model
	AgreeOrNot
	Content    string
	UserID     uint
	QuestionID uint
	IsDeleted  bool
	// 评论
	Comments []Comment `gorm:"foreignKey:AnswerID"`
}

var db *gorm.DB

const max_question_num = 30

var resultCh chan Question = make(chan Question, max_question_num)

func main() {
	var err error
	// 使用mysql.Open创建Dialector实例
	dialector := mysql.Open("root:123Qbz123@tcp(localhost:3306)/go_web_1?charset=utf8mb4&parseTime=True&loc=Local")
	// 将Dialector实例传递给gorm.Open
	db, err = gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}
	db.AutoMigrate(&User{}, &Question{}, &Answer{}, &Comment{}, &ReplyComment{})

	r := gin.Default()

	// 设置会话
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	v_post := r.Group("/")
	{
		//注册
		v_post.POST("register", Register)
		//登录
		v_post.POST("login", Login)
		// 提出问题
		v_post.POST("questions", CreateQuestion)
		// 回答问题
		v_post.POST("questions/:question_id/answers", CreateAnswer)
		// 评论答案,需要查询参数，传入答案id和所属问题id
		v_post.POST("questions/answers/comments", CreateCommentToAnswer)
		//回复其他评论,需要传入查询参数，传入其他评论id和所属答案id
		v_post.POST("questions/answers/comments/comments", CreateReplyToComment)

	}

	v_get := r.Group("/")
	{
		// 查看某个用户提出的问题
		v_get.GET(":user_id/questions", GetQuestions)
		// 查看某问题的答案
		v_get.GET("questions/:question_id/answers", GetAnswers)
		// 查看被逻辑删除删除的问题
		v_get.GET("deleted_questions", GetDeletedQuestions)
		//按关键词搜索问题
		v_get.GET("questions/search", SearchQuestionsByKeyword)
		//查看对自己评论的回复
		v_get.GET("questions/answers/comments/comments", GetReplysToSelf)

	}

	v_put := r.Group("/")
	{
		//申请管理员权限
		v_put.PUT("users", ApplyManager)
		// 更新问题内容
		v_put.PUT("questions/:id", UpdateQuestionContent)
		// 恢复被逻辑删除的问题
		v_put.PUT("deleted_questions/:id", RestoreQuestion)
		// 赞同答案
		v_put.PUT("questions/answers/agree", AgreeAnswer)
		// 反对答案
		v_put.PUT("questions/answers/disagree", DisagreeAnswer)
		//恢复被删除的答案,需要传入查询参数，传入答案id和所属问题id
		v_put.PUT("questions/answers/restore", RestoreAnswer)
		//恢复对答案的评论,需要传入查询参数，传入评论id和所属答案id
		v_put.PUT("questions/answers/comments/restore", RestoreCommentToAnswer)
	}

	v_delete := r.Group("/")
	{
		// 逻辑删除问题
		v_delete.DELETE("questions/:id", DeleteQuestion)
		// 永久删除问题
		v_delete.DELETE("deleted_questions/:id", DeleteQuestionForever)
		// 永久删除答案
		v_delete.DELETE("questions/answers", DeleteAnswerForever)
		//逻辑删除答案,需要传入查询参数，传入答案id和所属问题id
		v_delete.DELETE("questions/deleted_answers", DeleteAnswer)
		//软删除评论,需要传入查询参数，传入评论id和所属答案id
		v_delete.DELETE("questions/answers/comments", DeleteCommentToAnswer)
		//永久删除评论，需要传入查询参数，传入评论id和所属答案id
		v_delete.DELETE("questions/answers/deleted_comments", DeleteCommentForever)
	}

	//统计每个被删除的删除天数，超过15天则永久删除
	go func() {
		flag := true
		for flag {
			time.Sleep(time.Second * 30) // 添加延迟
			deleted_questions.mux.Lock()
			for _, q := range deleted_questions.deleted {
				if q.IsDeleted && time.Since(q.Deleted_Time) > time.Second*10 {
					resultCh <- q
					db.Unscoped().Delete(&q)
				}
			}
			deleted_questions.mux.Unlock()
		}
	}()

	// 收集结果
	go func() {
		for {
			select {
			case q := <-resultCh:
				delete(deleted_questions.deleted, strconv.FormatUint(uint64(q.ID), 10))
			}
		}
	}()

	r.Run(":8080")
}
