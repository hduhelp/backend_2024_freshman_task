package dao

import (
	"QASystem/models"
	"gorm.io/gorm"
)

func CreatePost(post *models.Post) error {
	return DB.Create(post).Error
}

func DeletePost(postId int64) error {
	return DB.Where("id = ?", postId).Delete(&models.Post{}).Error
}

func GetPostByID(postId int64) (*models.Post, error) {
	post := &models.Post{}
	return post, DB.Where("id = ?", postId).First(post).Error
}

func UpdatePost(post *models.Post) error {
	return DB.Model(&models.Post{}).Where("id = ?", post.ID).Updates(post).Error
}

func ViewPost(postId int64) error {
	return DB.Model(&models.Post{}).Where("id = ?", postId).UpdateColumn("views", gorm.Expr("views + ?", 1)).Error
}

func LikePost(postId int64, isLiked bool) error {
	if isLiked {
		return DB.Model(&models.Post{}).Where("id = ?", postId).UpdateColumn("likes", gorm.Expr("likes - ?", 1)).Error
	} else {
		return DB.Model(&models.Post{}).Where("id = ?", postId).UpdateColumn("likes", gorm.Expr("likes + ?", 1)).Error
	}
}

func PagePost(pageNum int, pageSize int) ([]*models.Post, error) {
	var posts []*models.Post
	return posts, DB.Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&posts).Error
}
