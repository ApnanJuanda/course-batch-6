package handler

import (
	"exercise/internal/app/domain"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AnswerHandler struct {
	db *gorm.DB
}

func NewAnswerHandler(db *gorm.DB) *AnswerHandler {
	return &AnswerHandler{db: db}
}

//create new answer
func (ah AnswerHandler) CreateNewAnswer(c *gin.Context) {
	var addAnswer domain.AddAnswer
	if err := c.ShouldBind(&addAnswer); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"message": "invalid body",
		})
		return
	}

	//convert questionId and exerciseId
	questionIdString := c.Param("questionId")
	questionId, err := strconv.Atoi(questionIdString)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"message": "invalid questionId",
		})
		return
	}

	exerciseIdString := c.Param("exerciseId")
	exerciseId, err := strconv.Atoi(exerciseIdString)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"message": "invalid id",
		})
		return
	}

	//check apakah ada question dengan exerciseId dan questionId
	var question domain.Question
	err = ah.db.Where("id = ? AND exercise_id = ?", questionId, exerciseId).Find(&question).Error
	if err != nil {
		c.JSON(http.StatusNotFound, map[string]string{
			"message": "not found question",
		})
		return
	}

	//check apakah sudah ada answer dengan exercise_id, question_id, userID
	userID := c.Request.Context().Value("user_id").(int)
	var answer domain.Answer
	err = ah.db.Where("exercise_id = ? AND question_id = ? AND user_id = ?", exerciseId, questionId, userID).Find(&answer).Error
	if err == nil {
		if answer.Answer != "" {
			//update oldAnswer
			var answer2 domain.Answer
			ah.db.Model(&answer2).Where("exercise_id = ? AND question_id = ? AND user_id = ?", exerciseId, questionId, userID).Update("answer", addAnswer.Answer)
			c.JSON(http.StatusCreated, map[string]string{
				"message": "Success update answer",
			})
			return
		}
		//save newAnswer
		newAnswer, err := domain.NewAnswer(addAnswer.Answer, exerciseId, questionId, userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, map[string]string{
				"message": err.Error(),
			})
			return
		}
		if err := ah.db.Create(newAnswer).Error; err != nil {
			c.JSON(http.StatusBadRequest, map[string]string{
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusCreated, map[string]string{
			"message": "Success add answer",
		})
	} else {
		c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
		return
	}

}
