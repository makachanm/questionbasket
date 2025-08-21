package question

import (
	"questionbasket/bot"
	"time"
)

type MakeQuestionResponse struct {
	QID string `json:"qid"`
}

// QuestionDetailResponse contains the detailed information of a question.
type QuestionDetailResponse struct {
	ID         string              `json:"id"`
	Content    string              `json:"content"`
	CreatedAt  time.Time           `json:"created_at"`
	Answer     *string             `json:"answer,omitempty"`
	ShareRange *bot.ShareRangeType `json:"share_range,omitempty"`
}

type QuestionAnswerResponse struct {
	QID string `json:"qid"`
}
