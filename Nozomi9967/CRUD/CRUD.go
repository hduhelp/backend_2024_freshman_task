package CRUD

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github/piexlMax/web/comment"
	"github/piexlMax/web/gorm"
	"github/piexlMax/web/post"
	"net/http"
	"strconv"
)

func SearchPostByindex(c *gin.Context) {
	indexstring := c.Param("index")
	index, error := strconv.Atoi(indexstring)
	if error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "index类型转换失败"})
		return
	}
	var post post.Post
	//order:按照id来升序排序，offset:跳过前index-1条数据，limit:查找index-1后的一条数据，find:赋值给post
	gorm.GLOBAL_DB.Table("t_post").Order("id").Offset(index - 1).Limit(1).Find(&post)
	fmt.Println("第", index, "条数据：", post)
	c.JSON(http.StatusOK, gin.H{"PostId": post.ID})
}

func DeletePost(c *gin.Context) {

	//json绑定前端帖子id信息
	var Postid struct {
		PostID uint
	}
	error := c.BindJSON(&Postid)
	if error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "绑定json失败"})
		return
	}
	fmt.Println("删除帖子函数接收到的postid为", Postid.PostID)

	//数据库总寻找对应id帖子
	var post post.Post
	result := gorm.GLOBAL_DB.Table("t_post").First(&post, Postid.PostID)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "没有找到对应帖子"})
		return
	}
	fmt.Println("待删除的post为", post)

	//删除对应帖子的评论
	var comments []comment.Comment
	gorm.GLOBAL_DB.Table("t_comment").Where("post_id = ?", post.ID).Find(&comments)
	fmt.Println("帖子对应评论：", comments)
	if len(comments) > 0 {
		gorm.GLOBAL_DB.Table("t_comment").Delete(&comments)
	}

	//删除帖子
	gorm.GLOBAL_DB.Delete(&post)
	c.JSON(200, gin.H{"msg": "删除成功"})
}

func GetModifInfo(c *gin.Context) {
	//string转int
	var postid = c.Param("postid")
	var PostIdInt, error = strconv.Atoi(postid)
	if error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "postid转换类型失败"})
		return
	}

	//数据库查找帖子，帖子数据绑定到post
	var post post.Post
	result := gorm.GLOBAL_DB.Table("t_post").First(&post, PostIdInt)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "帖子不存在"})
		return
	}

	//只需返回标题和内容
	c.JSON(http.StatusOK, gin.H{"PostHeadline": post.Headline, "PostContent": post.Content})
}

func ModifyPost(c *gin.Context) {

	//处理前端发来的信息
	var PostId = c.Param("postid")
	var PostIdInt, error = strconv.Atoi(PostId)
	if error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "postid转换类型失败"})
		return
	}
	var Postinfo struct {
		Headline string `json:"headline"`
		Content  string `json:"content"`
	}
	error = c.BindJSON(&Postinfo)
	if error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "json数据绑定失败"})
	}

	var post post.Post
	result := gorm.GLOBAL_DB.Table("t_post").First(&post, PostIdInt)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "帖子不存在"})
		return
	}

	post.Headline = Postinfo.Headline
	post.Content = Postinfo.Content
	fmt.Println("修改后的post：", post)
	gorm.GLOBAL_DB.Save(&post)

	c.JSON(200, gin.H{"msg": "修改成功"})
}
