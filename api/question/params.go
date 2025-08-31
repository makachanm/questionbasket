package question

import "questionbasket/bot"

type MakeQuestionParams struct {
	Content string `json:"content"`
	IsNSFW  bool   `json:"is_nsfw"`
}

type QuestionAnswerParams struct {
	Content    string             `json:"content"`
	ShareRange bot.ShareRangeType `json:"share_range"`
}
