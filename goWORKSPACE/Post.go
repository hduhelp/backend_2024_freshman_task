package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// ReplyPost 回帖
func ReplyPost(c *gin.Context) {
	postID := c.Param("postID")
	ID, _ := strconv.ParseUint(postID, 10, 64)
	username := c.Param("username")
	var comment Comment
	if err := c.BindJSON(&comment); err != nil {
		c.JSON(400, gin.H{"err": "无效的请求"})
	} else {
		value, ok := UserStore.Load(username)
		if !ok {
			c.JSON(401, gin.H{"error": "未读取到用户信息！"})
		}
		if value == nil {
			c.JSON(401, gin.H{"error": "用户信息为空！"})
		}
		var post Post
		if err = DB.Where("ID=?", ID).Find(&post).Error; err != nil {
			c.JSON(404, gin.H{"error": "评论的帖子不存在！"})
		} else {
			comment.AuthorID = value.(UserInfo).ID
			comment.PostID = postID
			comment.CreatTime = time.Now()
			comment.Post = post
			post.Comments = append(post.Comments, comment)
			if err = DB.Create(&comment).Error; err != nil {
				c.JSON(500, gin.H{"error": "创建评论失败！"})
			} else {
				c.JSON(200, gin.H{"data": comment})
				c.JSON(200, gin.H{"post": post})
			}
		}

	}
}

// DeletePost 删帖
func DeletePost(c *gin.Context) {
	postID := c.Param("postID")
	if err := DB.First(&Post{}, postID).Error; err != nil {
		c.JSON(404, gin.H{"error": "贴子不存在！"})
	}
	result := DB.Delete(&Post{}, postID)
	if result.Error != nil {
		c.JSON(400, gin.H{"error": "删除帖子失败！"})
	} else {
		c.JSON(200, gin.H{"msg": "删除帖子成功！"})
	}
}

// UpdatePost 更新帖子
func UpdatePost(c *gin.Context) {
	PostID := c.Param("postID")
	var updatePost Post
	if err := c.BindJSON(&updatePost); err != nil {
		c.JSON(400, gin.H{"error": "无效的请求数据！"})
	}
	updatePost.UpdateTime = time.Now()
	result := DB.Model(&Post{}).Where("id = ?", PostID).Updates(updatePost)
	if result.RowsAffected == 0 {
		c.JSON(404, gin.H{"error": "帖子不存在!"})
	} else if result.Error != nil {
		c.JSON(400, gin.H{"error": "更新失败！"})
	} else {
		c.JSON(200, gin.H{"msg": "更新成功！"})
	}
}

// SearchPost 查帖子*****
func SearchPost(c *gin.Context) {
	query := c.Query("query")
	var posts []Post
	if query != "" {
		DB.Where("Title LINK? OR content LIKE ?", "%"+query+"%", "%"+query+"%").Find(&posts)
	} else {
		SortByLike := c.Query("SortByLikes")
		if SortByLike == "true" || SortByLike == "" {
			if err := DB.Order("likes desc").Find(&posts).Error; err != nil {
				c.JSON(500, gin.H{"error": "查找帖子失败！"})
			}
		}
	}
	c.JSON(200, posts)
}

// Like 点赞帖子
func Like(c *gin.Context) {
	postID := c.Param("postID")
	var post Post
	result := DB.First(&post, postID)
	if result.RowsAffected == 0 {
		c.JSON(404, gin.H{"error": "帖子不存在！"})
	} else {
		post.Likes += 1
		DB.Save(&post)
		c.JSON(http.StatusOK, gin.H{"msg": "点赞成功！"})
	}
}

// CreatePost 发帖
func CreatePost(c *gin.Context) {
	var post Post
	username := c.Param("username")
	if err := c.BindJSON(&post); err != nil {
		c.JSON(400, gin.H{
			"error": "无效请求",
		})
	} else {
		value, ok := UserStore.Load(username)
		if !ok {
			c.JSON(401, gin.H{"error": "未读取到用户信息！"})
		}
		post.AuthorID = value.(UserInfo).ID
		post.CreateTime = time.Now()
		post.UpdateTime = time.Now()
		if err = DB.Create(&post).Error; err != nil {
			c.JSON(400, gin.H{"error": "创建帖子失败！"})
		}
		c.JSON(200, post)
	}
}

// DeleteComment 删除评论
func DeleteComment(c *gin.Context) {
	commentID := c.Param("commentID")
	PostID := c.Param("postID")
	var comment Comment
	if err := DB.First(&comment, commentID).Error; err != nil {
		c.JSON(404, gin.H{"error": "评论不存在！"})
	}
	var post Post
	if err := DB.First(&post, PostID).Error; err != nil {
		c.JSON(404, gin.H{"error": "未找到对应帖子！"})
	}
	var newComments []Comment
	for _, cmt := range post.Comments {
		if strconv.FormatUint(cmt.ID, 10) != commentID {
			newComments = append(newComments, cmt)
		}
	}
	post.Comments = newComments
	DB.Delete(&comment)
}
