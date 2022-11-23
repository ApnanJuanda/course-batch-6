package domain

import (
	"errors"
	"time"
)

type Exercise struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Questions   []Question `json:"quetions,omitempty"`
}

type AddExercise struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Question struct {
	ID            int       `json:"id"`
	ExerciseID    int       `json:"-"`
	Body          string    `json:"body"`
	OptionA       string    `json:"option_a"`
	OptionB       string    `json:"option_b"`
	OptionC       string    `json:"option_c"`
	OptionD       string    `json:"option_d"`
	CorrectAnswer string    `json:"-"`
	Score         int       `json:"score"`
	CreatorID     int       `json:"-"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type AddQuestion struct {
	Body          string `json:"body"`
	OptionA       string `json:"option_a"`
	OptionB       string `json:"option_b"`
	OptionC       string `json:"option_c"`
	OptionD       string `json:"option_d"`
	CorrectAnswer string `json:"correct_answer"`
}

type Answer struct {
	ID         int
	ExerciseID int
	QuestionID int
	UserID     int
	Answer     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type AddAnswer struct {
	Answer string `json:"answer"`
}

func NewExercise(title, description string) (*Exercise, error) {
	if title == "" {
		return nil, errors.New("title is required")
	}
	if description == "" {
		return nil, errors.New("description is required")
	}
	return &Exercise{
		Title:       title,
		Description: description,
	}, nil
}

func NewQuestion(body, option_a, option_b, option_c, option_d, correct_answer string, userId, exercise_id int) (*Question, error) {
	if body == "" {
		return nil, errors.New("body is required")
	}
	if option_a == "" {
		return nil, errors.New("option_a is required")
	}
	if option_b == "" {
		return nil, errors.New("option_b is required")
	}
	if option_a == "" {
		return nil, errors.New("option_c is required")
	}
	if option_a == "" {
		return nil, errors.New("option_d is required")
	}
	if correct_answer == "" {
		return nil, errors.New("correct_answer is required")
	}
	return &Question{
		ExerciseID:    exercise_id,
		Body:          body,
		OptionA:       option_a,
		OptionB:       option_b,
		OptionC:       option_c,
		OptionD:       option_d,
		CorrectAnswer: correct_answer,
		Score:         10,
		CreatorID:     userId,
	}, nil
}

func NewAnswer(answer string, exercise_id, question_id, userId int) (*Answer, error) {
	if answer == "" {
		return nil, errors.New("answer is required")
	}
	return &Answer{
		ExerciseID: exercise_id,
		QuestionID: question_id,
		UserID: userId,
		Answer: answer,

	}, nil
}
