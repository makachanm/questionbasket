package question

import (
	"database/sql"
	"log"
	"net/http"
	"questionbasket/bot"
	"questionbasket/config"
	"questionbasket/frame"
	"questionbasket/models"
)

type QuestionAPI struct {
	questionModel models.QuestionModel
	messageBot    bot.MessageBot
}

func NewQuestionAPI() QuestionAPI {
	qAPI := QuestionAPI{
		questionModel: models.QuestionModel{},
	}
	frame.DatabaseBind(&qAPI.questionModel)

	// Initialize bot from config
	cfg := config.Cfg
	if cfg.MastodonInstanceURL != "" && cfg.MastodonToken != "" && cfg.ServiceURL != "" {
		qAPI.messageBot = bot.NewBot(bot.MASTODON, cfg.MastodonInstanceURL, cfg.MastodonToken, cfg.ServiceURL)
		log.Println("Mastodon bot initialized.")
	} else {
		log.Println("Mastodon bot not configured. Skipping initialization.")
	}

	return qAPI
}

func (qAPI *QuestionAPI) MakeQuestion(ctx frame.APIContext) {
	params := new(MakeQuestionParams)
	err := ctx.GetContext(params)
	if err != nil {
		ctx.ReturnError("bad_request", "Invalid JSON format", http.StatusBadRequest)
		return
	}

	qid, aid, err := qAPI.questionModel.InsertQuestion(params.Content, params.IsNSFW)
	if err != nil {
		ctx.ReturnError("db_error", "Failed to create question", http.StatusInternalServerError)
		return
	}

	// Send notification in the background
	if qAPI.messageBot != nil {
		go func() {
			err := qAPI.messageBot.SendMessage(aid, params.IsNSFW, params.Content, bot.DIRECT)
			if err != nil {
				log.Printf("Failed to send message via bot: %v", err)
			}
		}()
	}

	resp := MakeQuestionResponse{
		QID: qid,
	}

	ctx.ReturnJSON(resp)
}

func (qAPI *QuestionAPI) GetQuestion(ctx frame.APIContext) {
	qid := ctx.GetPathParamValue("id")

	question, err := qAPI.questionModel.GetQuestion(qid)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.ReturnError("not_found", "Question not found", http.StatusNotFound)
		} else {
			ctx.ReturnError("db_error", "Failed to get question", http.StatusInternalServerError)
		}
		return
	}

	// Build the response
	resp := QuestionDetailResponse{
		ID:        question.QID,
		Content:   question.Content,
		CreatedAt: question.CreatedAt,
	}

	// Add answer if it exists
	if question.Answer.Valid {
		resp.Answer = &question.Answer.String
		shareRange := bot.ShareRangeType(question.ShareRange.Int64)
		resp.ShareRange = &shareRange
	}

	ctx.ReturnJSON(resp)
}

func (qAPI *QuestionAPI) GetRecentQuestions(ctx frame.APIContext) {
	questions, err := qAPI.questionModel.GetRecentQuestions()
	if err != nil {
		ctx.ReturnError("db_error", "Failed to get recent questions", http.StatusInternalServerError)
		return
	}

	if questions == nil {
		questions = []models.Question{}
	}

	// Build the response
	resp := make([]QuestionDetailResponse, len(questions))
	for i, q := range questions {
		item := QuestionDetailResponse{
			ID:        q.QID,
			Content:   q.Content,
			CreatedAt: q.CreatedAt,
		}
		// Add answer if it exists
		if q.Answer.Valid {
			item.Answer = &q.Answer.String
			shareRange := bot.ShareRangeType(q.ShareRange.Int64)
			item.ShareRange = &shareRange
		}
		resp[i] = item
	}

	ctx.ReturnJSON(resp)
}

func (qAPI *QuestionAPI) AnswerQuestion(ctx frame.APIContext) {
	// Define a local struct to extract 'aid' from query params
	type aidQueryParam struct {
		AID string `param:"aid"`
	}
	aidHolder := new(aidQueryParam)
	ctx.GetContext(aidHolder)
	aid := aidHolder.AID

	if aid == "" {
		ctx.ReturnError("bad_request", "Missing 'aid' query parameter", http.StatusBadRequest)
		return
	}

	// Get 'content' and 'share_range' from JSON body
	params := new(QuestionAnswerParams)
	err := ctx.GetContext(params)
	if err != nil {
		ctx.ReturnError("bad_request", "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Call the model function to add the answer
	qid, err := qAPI.questionModel.AddAnswerByAID(aid, params.Content, int(params.ShareRange))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.ReturnError("not_found", "Question not found for the given aid", http.StatusNotFound)
		} else {
			ctx.ReturnError("db_error", "Failed to add answer", http.StatusInternalServerError)
		}
		return
	}

	// Return the QID in the response
	resp := QuestionAnswerResponse{
		QID: qid,
	}

	ctx.ReturnJSON(resp)
}
