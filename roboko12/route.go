package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Deleted_Questions struct {
	deleted map[string]Question
	mux     sync.Mutex
}

var deleted_questions Deleted_Questions = Deleted_Questions{deleted: make(map[string]Question)}

// 检查是否登录
func CheckLogin(c *gin.Context) interface{} {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID == nil {
		return nil
	}
	return userID
}

// 初始化答案的属性，只初始化了部分属性，意味着诸如创建时间的属性不会被初始化
// 只在删除答案中可以使用

func (answer *Answer) Init(c *gin.Context, user_id uint) {

	answer_id := c.Query("answer_id")
	question_id := c.Query("question_id")
	var err error
	var int_answer_id int
	var int_question_id int
	if int_answer_id, err = strconv.Atoi(answer_id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if int_question_id, err = strconv.Atoi(question_id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	answer.UserID = user_id
	answer.ID = uint(int_answer_id)
	answer.QuestionID = uint(int_question_id)
}

// 初始化评论的属性
func (comment *Comment) Init(c *gin.Context) {

	comment_id := c.Query("comment_id")
	answer_id := c.Query("answer_id")
	var int_comment_id int
	var int_answer_id int
	var err error
	if int_comment_id, err = strconv.Atoi(comment_id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if int_answer_id, err = strconv.Atoi(answer_id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	comment.AnswerID = uint(int_answer_id)
	comment.ID = uint(int_comment_id)

}

// 申请管理员权限
func ApplyManager(c *gin.Context) {
	userID := CheckLogin(c)
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录"})
		return
	}
	id := userID.(uint)
	var answer Answer
	var answer_sum int64
	var totalAgreeCount int64
	if err := db.Model(&answer).Where("user_id =?", id).Count(&answer_sum).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if answer_sum < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "您的回答数量不足"})
		return
	}
	if err := db.Model(&Answer{}).Where("user_id =?", id).Select("SUM(agree_num)").Scan(&totalAgreeCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if totalAgreeCount < 1000 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "您的赞同数量不足"})
		return
	}
	// 更新用户为管理员
	if err := db.Model(&User{}).Where("id =?", id).Update("is_manager", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "申请管理员成功"})
}

// 检查用户是否是管理员
func CheckManager(user_id uint) bool {
	var user User
	if err := db.Where("id =?", user_id).First(&user).Error; err != nil {
		return false
	}
	return user.IsManager

}

// 注册
func Register(c *gin.Context) {
	var user User
	user.IsManager = false
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Println("error:", err.Error())
		return
	}
	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "用户注册成功"})
}

// 登录
func Login(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var storedUser User
	if err := db.Where("username =?", user.Username).First(&storedUser).Error; err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if storedUser.Password != user.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}
	session := sessions.Default(c)
	session.Set("user_id", storedUser.ID)
	session.Save()
	c.JSON(http.StatusOK, gin.H{"message": "登录成功"})
}

// 创建问题
func CreateQuestion(c *gin.Context) {
	userID := CheckLogin(c)
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录"})
		return
	}
	var question Question
	if err := c.ShouldBindJSON(&question); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	question.UserID = userID.(uint)
	question.Answers = []Answer{}
	question.IsDeleted = false

	//将被删除时间设为一个未来的时间点，代表该问题未被删除
	question.Deleted_Time = time.Date(3000, 1, 1, 1, 1, 1, 1, time.UTC)
	if err := db.Create(&question).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "问题提出成功", "question": question})
}

// 创建答案
func CreateAnswer(c *gin.Context) {
	userID := CheckLogin(c)
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录"})
		return
	}
	var answer Answer
	var err error
	var question_id int
	if question_id, err = strconv.Atoi(c.Param("question_id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	answer.QuestionID = uint(question_id)
	if err := c.ShouldBindJSON(&answer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var question Question
	if err := db.Where("id =?", answer.QuestionID).First(&question).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if question.IsDeleted {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该问题已被删除"})
		return
	}

	// 初始化答案的属性
	answer.UserID = userID.(uint)
	answer.CreatedAt = time.Now()
	answer.AgreeNum = 0
	answer.IsDeleted = false
	answer.DisagreeNum = 0
	answer.Comments = []Comment{}

	if err := db.Create(&answer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "回答创建失败", "error": err.Error()})
		return
	}

	//将回答添加到问题中
	question.Answers = append(question.Answers, answer)
	if err := db.Save(&question).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "回答创建成功", "answer": answer})
}

// 赞同答案
func AgreeAnswer(c *gin.Context) {
	userID := CheckLogin(c)
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录"})
		return
	}

	var answer Answer
	answer.Init(c, userID.(uint))

	if err := db.Where("id =?", answer.ID).First(&answer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	answer.AgreeNum++
	if err := db.Save(&answer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "赞同成功"})
}

// 反对答案
func DisagreeAnswer(c *gin.Context) {
	userID := CheckLogin(c)
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录"})
		return
	}

	var answer Answer
	answer.Init(c, userID.(uint))

	if err := db.Where("id =?", answer.ID).First(&answer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	answer.DisagreeNum++
	if err := db.Save(&answer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "反对成功"})
}

// 对答案评论
func CreateComment(c *gin.Context) {
	UserID := CheckLogin(c)
	if UserID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录"})
		return
	}

}

//对别人的评论进行评论

// 查看某用户提出的问题
func GetQuestions(c *gin.Context) {
	var questions []Question
	user_id := c.Param("user_id")
	sort_way := c.DefaultQuery("sort_way", "created_at desc")
	if err := db.Where("is_deleted = false AND user_id = ?", user_id).Order(sort_way).Find(&questions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"questions": questions})
}

// 查看某问题答案
func GetAnswers(c *gin.Context) {
	var answers []Answer
	question_id := c.Param("question_id")
	sort_way := c.DefaultQuery("sort_way", "agree_num desc")
	if err := db.Where("question_id =?", question_id).Find(&answers).Order(sort_way).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"answers": answers})
}

func GetDeletedQuestions(c *gin.Context) {
	userID := CheckLogin(c)
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录"})
		return
	}
	var questions []Question
	if err := db.Where("is_deleted = true").Find(&questions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"questions": questions})
}

// 更新问题内容
func UpdateQuestionContent(c *gin.Context) {
	userID := CheckLogin(c)
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录"})
		return
	}
	id := c.Param("id")
	var question Question
	if err := db.Where("id =?", id).First(&question).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := c.ShouldBindJSON(&question); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.Save(&question).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "问题修改成功", "question": question})
}

// 恢复被逻辑删除的问题
func RestoreQuestion(c *gin.Context) {
	userID := CheckLogin(c)
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录"})
		return
	}
	id := c.Param("id")
	var question Question
	if err := db.Where("id =?", id).First(&question).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	question.IsDeleted = false
	if err := db.Save(&question).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "问题恢复成功", "question": question})
}

// 逻辑删除问题
func DeleteQuestion(c *gin.Context) {
	userID := CheckLogin(c)
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录"})
		return
	}
	id := c.Param("id")
	delete_forever := c.Query("delete_forever")
	var question Question
	if err := db.Where("id =?", id).First(&question).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if question.UserID != userID.(uint) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "您没有权限删除该问题"})
		return
	}

	if delete_forever == "false" {
		//逻辑删除问题的所有答案
		for _, answer := range question.Answers {
			answer.delete_answer(&question)
		}
		//保存被删除的问题到map中
		deleted_questions.deleted[id] = question
		question.Deleted_Time = time.Now()
		question.IsDeleted = true
		if err := db.Save(&question).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "问题删除成功, 该问题已被逻辑删除，可自行选择是否永久删除"})
	}
	if delete_forever == "true" {
		//永久删除问题的所有答案
		for _, answer := range question.Answers {
			if answer.delete_answer_forever(&question) {
				c.JSON(http.StatusOK, gin.H{"message": "答案删除成功"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "答案删除失败"})
		}
		if err := db.Unscoped().Delete(&question).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "问题删除成功, 该问题已被永久删除"})
		return
	}

}

// 永久删除问题
func DeleteQuestionForever(c *gin.Context) {
	userID := CheckLogin(c)
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录"})
		return
	}

	user_id := userID.(uint)

	id := c.Param("id")
	var question Question
	if err := db.Where("id =?", id).First(&question).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if question.UserID != user_id && CheckManager(user_id) == false {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "您没有权限删除该问题"})
		return
	}

	//删除问题的所有答案
	for _, answer := range question.Answers {
		if answer.delete_answer_forever(&question) {
			c.JSON(http.StatusOK, gin.H{"message": "答案删除成功"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "答案删除失败"})
	}

	//然后删除问题
	if err := db.Unscoped().Delete(&question).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "问题删除成功"})

}

// 永久删除答案
func (answer *Answer) delete_answer_forever(question *Question) bool {
	//将答案从所属问题中删除
	len := len(question.Answers)
	question.Answers = append(question.Answers[:len], question.Answers[len:]...)

	//删除所属的所有评论

	if err := db.Where("id =?", answer.ID).First(&answer).Error; err != nil {
		log.Output(2, err.Error())
		return false
	}
	if err := db.Unscoped().Delete(&answer).Error; err != nil {
		log.Output(2, err.Error())
		return false
	}
	return true
}

func DeleteAnswerForever(c *gin.Context) {
	userID := CheckLogin(c)
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录"})
		return
	}

	var answer Answer
	answer.Init(c, userID.(uint))

	if err := db.Where("id =?", answer.ID).First(&answer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if answer.UserID != userID.(uint) && !CheckManager(userID.(uint)) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "您没有权限删除该答案"})
		return
	}

	var question Question

	if err := db.Where("id =?", answer.QuestionID).First(&question).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if answer.delete_answer_forever(&question) {
		c.JSON(http.StatusOK, gin.H{"message": "答案删除成功"})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": "答案删除失败"})
}

//逻辑删除答案

func (answer *Answer) delete_answer(question *Question) bool {
	answer.IsDeleted = true
	len := len(question.Answers)
	question.Answers = append(question.Answers[:len], question.Answers[len:]...)
	db.Delete(&answer)
	//逻辑删除答案的所有评论

	if err := db.Save(&answer).Error; err != nil {
		log.Output(2, err.Error())
		return false
	}
	return true
}

func DeleteAnswer(c *gin.Context) {
	userID := CheckLogin(c)
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录"})
		return
	}

	var answer Answer
	answer_id := c.Query("answer_id")
	if err := db.Where("id =?", answer_id).First(&answer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if answer.UserID != userID.(uint) && !CheckManager(userID.(uint)) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "您没有权限删除该问题"})
		return
	}
	var question Question

	if err := db.Where("id =?", answer.QuestionID).First(&question).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if answer.delete_answer(&question) {
		c.JSON(http.StatusOK, gin.H{"message": "答案删除成功,可以自行选择是否永久删除"})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": "答案删除失败"})
}

// 恢复答案
func RestoreAnswer(c *gin.Context) {
	userID := CheckLogin(c)
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录"})
		return
	}

	var answer Answer
	answer.Init(c, userID.(uint))

	answer.IsDeleted = false
	db.Unscoped().Model(&answer).Update("deleted_at", nil)
	db.Unscoped().Model(&answer).Update("IsDeleted", false)

	var question Question

	if err := db.Where("id =?", answer.QuestionID).First(&question).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	question.Answers = append(question.Answers, answer)

	c.JSON(http.StatusOK, gin.H{"message": "答案恢复成功"})
}

// 评论答案
func CreateCommentToAnswer(c *gin.Context) {
	userID := CheckLogin(c)
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录"})
		return
	}

	var answer Answer
	answer_id := c.Query("answer_id")
	//查找答案
	if err := db.Where("id =?", answer_id).First(&answer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var comment Comment
	if c.ShouldBindJSON(&comment) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数有误"})
		return
	}

	comment.AnswerID = answer.ID
	comment.UserID = userID.(uint)
	comment.AgreeNum = 0
	comment.DisagreeNum = 0
	comment.CreatedAt = time.Now()
	if err := db.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	answer.Comments = append(answer.Comments, comment)

	if err := db.Save(&answer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	//将评论存入数据库
	if err := db.Save(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "评论成功"})

}

// 软删除对答案的评论
func (comment *Comment) delete_comment_to_answer(answer *Answer) {
	db.Delete(comment)
	len := len(answer.Comments)
	answer.Comments = append(answer.Comments[:len], answer.Comments[len:]...)
}

func DeleteCommentToAnswer(c *gin.Context) {
	userID := CheckLogin(c)
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录"})
		return
	}
	var comment Comment
	comment.Init(c)

	if err := db.Where("id =?", comment.ID).First(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if comment.UserID != userID.(uint) && CheckManager(userID.(uint)) == false {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "您没有权限删除该问题"})
		return
	}

	var answer Answer
	if err := db.Where("id =?", comment.AnswerID).First(&answer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	comment.delete_comment_to_answer(&answer)

	c.JSON(http.StatusOK, gin.H{"message": "评论删除成功,可自行选择恢复"})

}

// 恢复对答案的评论
func RestoreCommentToAnswer(c *gin.Context) {
	userID := CheckLogin(c)
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录"})
		return
	}

	var comment Comment
	comment.Init(c)
	db.Unscoped().Model(&comment).Update("deleted_at", nil)

	var answer Answer
	if err := db.Where("id =?", comment.AnswerID).First(&answer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	answer.Comments = append(answer.Comments, comment)

	c.JSON(http.StatusOK, gin.H{"message": "评论恢复成功"})
}

//永久删除对答案的评论

func (comment *Comment) delete_comment_forever(answer *Answer) {
	db.Unscoped().Delete(comment)
	len := len(answer.Comments)
	answer.Comments = append(answer.Comments[:len], answer.Comments[len:]...)
}
func DeleteCommentForever(c *gin.Context) {
	userID := CheckLogin(c)
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录"})
		return
	}

	var comment Comment
	comment.Init(c)

	if err := db.Unscoped().Where("id =?", comment.ID).First(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if comment.UserID != userID.(uint) && CheckManager(userID.(uint)) == false {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "您没有权限删除该问题"})
		return
	}

	var answer Answer
	if err := db.Where("id =?", comment.AnswerID).First(&answer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	comment.delete_comment_forever(&answer)

	c.JSON(http.StatusOK, gin.H{"message": "评论删除成功"})

}

// 关键词搜索
func SearchQuestionsByKeyword(c *gin.Context) {
	var keyword KeyWord
	if err := c.ShouldBindJSON(&keyword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数有误"})
		return
	}
	if keyword.Keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请输入搜索关键词"})
		return
	}

	var questions []Question
	if err := db.Where("UPPER(title) LIKE UPPER(?) OR UPPER(content) LIKE UPPER(?)", "%"+keyword.Keyword+"%", "%"+keyword.Keyword+"%").Find(&questions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"questions": questions})
}
