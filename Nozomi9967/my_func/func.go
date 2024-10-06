package my_func

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github/piexlMax/web/comment"
	"github/piexlMax/web/gorm"
	"github/piexlMax/web/post"
	"net/http"
	"reflect"
	"strconv"
	"time"
)

var SinglePageSize = 10

type CustomClaims struct {
	UserId uint `json:"user_id"`
	jwt.StandardClaims
}

type User struct {
	ID   uint `gorm:"primary_key"`
	Name string
	Psw  string
}

func Init() {
	CreateUserTable()
	CreatePostTable()
	comment.CreateCommentTable()
}

func CreateUserTable() {
	gorm.GLOBAL_DB.AutoMigrate(&User{})
}

func CreatePostTable() {
	gorm.GLOBAL_DB.AutoMigrate(&post.Post{})
}

func Menu(c *gin.Context) {
	c.HTML(http.StatusOK, "menu", nil)
}

func UploadPostInfo(c *gin.Context) {
	var post post.Post
	err := c.BindJSON(&post)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的json格式"})
		return
	}

	//关联发布者id
	UseridInterface, exist := c.Get("user_id")

	fmt.Println("UploadPostInfo中接收到的userid", UseridInterface)
	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{"error": "获取用户id失败"})
		return
	}
	post.UserID = UseridInterface.(uint)

	//加入上传时间
	post.PostTime = time.Now()
	//格式化发布时间
	post.FormattedTime = post.PostTime.Format("2006-01-02 15:04:05")

	gorm.GLOBAL_DB.Create(&post)

	c.JSON(http.StatusOK, gin.H{"msg": "帖子上传成功"})
}

func GetPosts(c *gin.Context) {

	//绑定搜索关键词
	var SearchKeyWord struct {
		SKW string
	}
	err := c.ShouldBindJSON(&SearchKeyWord)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "绑定搜索词失败"})
		return
	}
	fmt.Println("获取到的搜索词为", SearchKeyWord.SKW)

	//获取页码
	var PageIndex = c.Param("index")
	PageIndexInt, Err := strconv.Atoi(PageIndex)
	if Err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "转换类型失败"})
		return
	}
	fmt.Println("获取到的页码为：", PageIndexInt)

	var posts []post.Post
	if SearchKeyWord.SKW == "没有搜索词" {
		gorm.GLOBAL_DB.Table("t_post").Offset((PageIndexInt - 1) * SinglePageSize).Limit(10).Find(&posts)
	} else {
		gorm.GLOBAL_DB.Table("t_post").Where("headline LIKE ?", "%"+SearchKeyWord.SKW+"%").Offset((PageIndexInt - 1) * SinglePageSize).Limit(10).Find(&posts)
	}
	fmt.Println("获取到的：", posts)
	//if len(posts) > 0 {
	//	c.JSON(http.StatusOK, gin.H{"posts": posts})
	//} else {
	//	c.JSON(http.StatusNotFound, gin.H{"error": "还没有人发帖子"})
	//}

	c.JSON(http.StatusOK, gin.H{"posts": posts})

}

// 获取帖子初始化数据
func GetIdpostInfo(c *gin.Context) {

	var PostId struct {
		ID uint `json:"id"`
	}

	err := c.BindJSON(&PostId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "绑定json失败"})
		return
	}

	//根据id查找帖子数据
	fmt.Println("后端得到的postid为", PostId.ID)

	var post post.Post
	result := gorm.GLOBAL_DB.First(&post, "ID=?", PostId.ID)
	//fmt.Println("返回值信息为", result)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "没有对应id的帖子"})
		return
	}
	fmt.Println("帖子id：", post.ID)
	fmt.Println("发布帖子的用户id：", post.UserID)
	fmt.Println("帖子标题：", post.Headline)
	fmt.Println("帖子内容：", post.Content)

	//根据帖子id绑定用户
	var post_user User
	result = gorm.GLOBAL_DB.First(&post_user, "ID=?", post.UserID)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "此帖子的用户不存在"})
		return
	}

	//寻找该帖子下的评论
	TopcommentsJson, err := comment.GetAllComment(PostId.ID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "查找评论信息失败"})
	}

	//找到了，返回数据到前端去渲染
	c.JSON(http.StatusOK, gin.H{"post": post, "post_user": post_user, "comments": TopcommentsJson})
}

func Upload_comment(c *gin.Context) {

	fmt.Println("Upload_comment函数执行了")

	var comment_part struct {
		Content  string `form:"content"`
		PostId   uint   `form:"post_id"`
		ParentID uint   `form:"parent_id"`
	}
	//将前端接收的评论数据保存至comment
	c.BindJSON(&comment_part)

	//转移到Comment结构体中
	var comment comment.Comment
	comment.Content = comment_part.Content
	comment.PostID = comment_part.PostId
	comment.ParentID = comment_part.ParentID

	userinterface, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{"error": "获取发布评论者id失败"})
	}
	userid := userinterface.(uint)
	var user User
	result := gorm.GLOBAL_DB.First(&user, "ID=?", userid)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "没有找到对应id的用户名"})
		return
	}
	comment.UserName = user.Name
	comment.UserID = userid

	comment.FormattedTime = time.Now().Format("2006-01-02 15:04:05")

	fmt.Println(comment)

	gorm.GLOBAL_DB.Create(&comment)
	c.JSON(http.StatusOK, gin.H{"msg": "评论成功"})
}

func Main_page_re(c *gin.Context) {
	//c.Redirect(http.StatusFound, "/auth/main")
	fmt.Println("主界面")
	c.HTML(http.StatusOK, "main", nil)
}

func Regis_page(c *gin.Context) {
	c.HTML(http.StatusOK, "regis", nil)
}

func Replication(Name string) bool {
	if gorm.GLOBAL_DB.First(&User{}, "Name=?", Name).Error == nil {
		fmt.Println(Name)
		return true
	}
	return false
}

func Regis(c *gin.Context) {
	var user User
	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的jason格式"})
		return
	}
	if Replication(user.Name) {
		c.JSON(http.StatusConflict, gin.H{"error": "用户名已存在"})
		fmt.Println("冲突+", user.Name)
		return
	}
	gorm.GLOBAL_DB.Create(&user)
	c.JSON(http.StatusOK, gin.H{"msg": "注册成功"})
}

func Login_page(c *gin.Context) {
	c.HTML(http.StatusOK, "login", nil)
}

func Login(c *gin.Context) {
	var user User
	var loginData struct {
		ID   uint
		Name string
		Psw  string
	}
	err := c.BindJSON(&loginData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的jason格式"})
		return
	}
	error := gorm.GLOBAL_DB.First(&user, "Name=?", loginData.Name)
	if error.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该账号不存在"})
		return
	}
	if user.Psw != loginData.Psw {
		c.JSON(http.StatusBadRequest, gin.H{"error": "密码错误"})
		return
	}
	TokenString, err := GenerateToken(user.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "生成token失败"})
		return
	}

	c.Set("user_id", uint(user.ID))
	fmt.Println("login中的userid", uint(user.ID))
	//session := sessions.Default(c)
	//session.Set("user_id", user.ID)
	//session.Save()
	fmt.Println("login中的token")
	fmt.Println(TokenString)
	c.JSON(http.StatusOK, gin.H{"msg": "登录成功", "token": TokenString, "userid": user.ID})

}

func GenerateToken(userid uint) (string, error) {

	fmt.Println("编进token之前的userid类型为", reflect.TypeOf(userid))

	claims := &CustomClaims{
		UserId: userid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 2).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenstring, err := token.SignedString([]byte("nozomi"))
	if err != nil {
		return "", err
	}
	return tokenstring, nil
}

func AuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		tokenString := c.Request.Header.Get("Authorization")
		if tokenString == "" {
			//fmt.Println("未提供token")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "未提供token"})
			return
		}
		token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte("nozomi"), nil
		})
		//fmt.Println("中间件接收到的token")
		fmt.Println(token, err)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "无效的token"})
			return
		}
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			//fmt.Println("中间件在context中保存的userid", claims.UserId)
			//fmt.Println("类型为",type(claims["userid"]))

			//存在上下文中，方便以后函数使用
			c.Set("user_id", claims.UserId)
			//fmt.Println("next执行了")
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "无效的token"})
			return
		}

		//session实现失败
		//session := sessions.Default(c)
		//userid := session.Get("user_id")
		//if userid == nil {
		//	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		//	return
		//}
		//fmt.Println("中间件")
		//fmt.Println(userid)
		//c.Set("user_id", userid)
		//c.Next()

	}
}

//
//func Main_page_(c *gin.Context) {
//	c.Redirect(http.StatusFound, "/auth/main")
//	//c.HTML(http.StatusOK, "main", nil)
//}

func Post_question_page(c *gin.Context) {
	c.HTML(http.StatusOK, "postquestion", nil)
}
