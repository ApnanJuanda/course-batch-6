package handler

import (
	"exercise/internal/app/domain"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type QuestionHandler struct {
	db *gorm.DB
}

func NewQuestionHandler(db *gorm.DB) *QuestionHandler {
	return &QuestionHandler{db: db}
}

//create new Question
func (qh QuestionHandler) CreateNewQuestion(c *gin.Context) {
	var addQuestion domain.AddQuestion
	if err := c.ShouldBind(&addQuestion); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"message" : "invalid body",
		})
		return
	}
	idString := c.Param("exerciseId")
	id, err := strconv.Atoi(idString)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"message": "invalid id",
		})
		return
	}

	//findById exercise
	var exercise domain.Exercise
	err = qh.db.Where("id = ?", id).Preload("Questions").Take(&exercise).Error
	if err != nil {
		c.JSON(http.StatusNotFound, map[string]string{
			"message": "exercise not found",
		})
		return
	}

	//add new question to exercise
	userID := c.Request.Context().Value("user_id").(int)
	newQuestion, err := domain.NewQuestion(addQuestion.Body, addQuestion.OptionA, addQuestion.OptionB, addQuestion.OptionC, addQuestion.OptionD, addQuestion.CorrectAnswer, userID, exercise.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
		return
	}
	
	//save newQuestion
	if err := qh.db.Create(newQuestion).Error; err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, map[string]string{
		"message" : "Success add question",
	})
} 
