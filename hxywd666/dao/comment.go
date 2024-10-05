package dao

import "QASystem/models"

func CountCommentsByPostID(postID int64) (int64, error) {
	var count int64
	return count, DB.Where("post_id=?", postID).Find(&models.Comment{}).Count(&count).Error
}

func CreateComment(comment *models.Comment) error {
	return DB.Create(&comment).Error
}

func DeleteComment(commentID int64) error {
	return DB.Where("id = ?", commentID).Delete(&models.Comment{}).Error
}

func GetComment(postID int64, page int, pageSize int, orderBy int) ([]*models.Comment, error) {
	var comments []*models.Comment
	if orderBy == 1 {
		return comments, DB.Where("post_id = ?", postID).Offset((page - 1) * pageSize).Limit(pageSize).Order("create_time desc").Find(&comments).Error
	} else {
		return comments, DB.Where("post_id = ?", postID).Offset((page - 1) * pageSize).Limit(pageSize).Order("likes desc").Find(&comments).Error
	}
}
