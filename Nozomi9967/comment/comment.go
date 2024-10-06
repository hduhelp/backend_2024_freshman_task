package comment

import (
	"fmt"
	"github/piexlMax/web/gorm"
)

type Comment struct {
	ID            uint   `gorm:"primary_key"`
	PostID        uint   `gorm:"index;not null"`
	UserID        uint   `gorm:"index"` //发布评论的用户的id
	UserName      string `gorm:"type:text"`
	ParentID      uint   `gorm:"index"` //父评论的id
	Content       string `gorm:"type:text"`
	FormattedTime string `gorm:"type:text"`
	//Post          post.Post `gorm:"foreignkey:PostID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Replies []Comment `gorm:"foreignKey:ParentID;references:ID"` // 关联关系定义
}

type TopcommentJson struct {
	ID            uint             `json:"id"`
	Username      string           `json:"username"`
	Content       string           `json:"content"`
	FormattedTime string           `json:"formatted_time"`
	Replies       []TopcommentJson `json:"replies"`
}

func CreateCommentTable() {
	gorm.GLOBAL_DB.AutoMigrate(&Comment{})
}

func GetAllComment(Postid uint) ([]TopcommentJson, error) {
	//var SummitComments struct {
	//	ID            uint      `json:"id"`
	//	UserID        uint      `json:"user_id"`
	//	Content       string    `json:"content"`
	//	CreatedAt     time.Time `json:"created_at"`
	//	FormattedTime string    `json:"formatted_time"`
	//}

	//fmt.Println("GetAllComment执行")

	//添加回复
	var TopComments []Comment
	gorm.GLOBAL_DB.Where("post_id = ? AND parent_id=0", Postid).Find(&TopComments)

	//帖子没有评论直接返回
	if TopComments == nil {

		return nil, nil
	}

	var EnrichedTopComments []Comment
	for _, TopComment := range TopComments {

		var Replies []Comment

		gorm.GLOBAL_DB.Where("parent_id=?", TopComment.ID).Find(&Replies)
		fmt.Println("ID=", TopComment.ID)
		fmt.Println("找到的Replies为", Replies)
		TopComment.Replies = Replies
		EnrichedTopComments = append(EnrichedTopComments, TopComment)
	}

	//fmt.Println("未处理之前的TopComments为", TopComments)

	//转json
	var topcommentsJson []TopcommentJson
	for _, EnrichedTopComment := range EnrichedTopComments {

		//先把父评论的回复转为json
		topcommentRepliesJson := make([]TopcommentJson, len(EnrichedTopComment.Replies))
		for i, Reply := range EnrichedTopComment.Replies {
			topcommentRepliesJson[i] = TopcommentJson{
				ID:            Reply.ID,
				Username:      Reply.UserName,
				Content:       Reply.Content,
				FormattedTime: Reply.FormattedTime,
			}
		}

		//fmt.Println("转json执行")
		//fmt.Println("ID为", EnrichedTopComment.ID)
		//fmt.Println(topcommentRepliesJson)

		//将父评论的其它变量转为json
		topcommentsJson = append(topcommentsJson, TopcommentJson{
			ID:            EnrichedTopComment.ID,
			Username:      EnrichedTopComment.UserName,
			Content:       EnrichedTopComment.Content,
			FormattedTime: EnrichedTopComment.FormattedTime,
			Replies:       topcommentRepliesJson,
		})

	}

	//JsonData, err := json.Marshal(topcommentsJson)
	//if err != nil {
	//	fmt.Println(err)
	//	return nil, err
	//}
	return topcommentsJson, nil
}
