package question

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"l1nk4i/db"
	"net/http"
)

func Delete(c *gin.Context) {
	questionID := c.Param("question-id")

	// Verify user identity
	session := sessions.Default(c)
	userid, exists := session.Get("user_id").(string)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid session"})
		return
	}

	question, err := db.GetQuestionByQuestionID(questionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid question_id"})
		return
	}

	if question.UserID != userid {
		c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
		return
	}

	// Delete question
	err = db.DeleteQuestion(questionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "delete question error"})
		return
	}

	// Delete answers to the question
	err = deleteAnswers(questionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "delete answers error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "delete question successful!"})
}

func deleteAnswers(questionID string) error {
	answers, err := db.GetAnswersByQuestionID(questionID)
	if err != nil {
		return err
	}

	for _, answer := range *answers {
		err = db.DeleteAnswer(answer.AnswerID)
		if err != nil {
			return err
		}
	}

	return nil
}
