package handler

import (
	"exercise/internal/app/domain"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ExerciseHandler struct {
	db *gorm.DB
}

func NewExerciseHandler(db *gorm.DB) *ExerciseHandler {
	return &ExerciseHandler{db: db}
}

//create new Exercise
func (eh ExerciseHandler) CreateNewExercise(c *gin.Context) {
	var addExercise domain.AddExercise
	if err := c.ShouldBind(&addExercise); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"message": "invalid body",
		})
		return
	}
	newExercise, err := domain.NewExercise(addExercise.Title, addExercise.Description)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
		return
	}
	if err := eh.db.Create(newExercise).Error; err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, newExercise)
}

func (eh ExerciseHandler) GetExerciseByID(c *gin.Context) {
	idString := c.Param("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"message": "invalid id",
		})
		return
	}

	var exercise domain.Exercise
	err = eh.db.Where("id = ?", id).Preload("Questions").Take(&exercise).Error
	if err != nil {
		c.JSON(http.StatusNotFound, map[string]string{
			"message": "exercise not found",
		})
		return
	}
	c.JSON(http.StatusOK, exercise)
}

func (eh ExerciseHandler) GetScore(c *gin.Context) {
	idString := c.Param("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"message": "invalid id",
		})
		return
	}

	var exercise domain.Exercise
	err = eh.db.Where("id = ?", id).Preload("Questions").Take(&exercise).Error
	if err != nil {
		c.JSON(http.StatusNotFound, map[string]string{
			"message": "exercise not found",
		})
		return
	}
	userID := c.Request.Context().Value("user_id").(int)

	var answers []domain.Answer
	err = eh.db.Where("exercise_id = ? AND user_id = ?", id, userID).Find(&answers).Error
	if err != nil {
		c.JSON(http.StatusNotFound, map[string]string{
			"message": "not answere yet",
		})
		return
	}

	mapQA := make(map[int]domain.Answer)
	for _, answer := range answers {
		mapQA[answer.QuestionID] = answer
	}

	var score Score
	wg := new(sync.WaitGroup)
	for _, question := range exercise.Questions {
		wg.Add(1)
		go func(question domain.Question) {
			defer wg.Done()
			if strings.EqualFold(question.CorrectAnswer, mapQA[question.ID].Answer) {
				score.Inc(question.Score)
			}
		}(question)
	}

	wg.Wait()

	c.JSON(http.StatusOK, map[string]int{
		"score": score.total,
	})
}

type Score struct {
	total int
	mu    sync.Mutex
}

func (s *Score) Inc(value int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.total += value
}
