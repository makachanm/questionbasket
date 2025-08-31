package api

import (
	"questionbasket/api/profile"
	"questionbasket/api/question"
	"questionbasket/frame"
	"questionbasket/middlewares"
)

type API struct {
	Router frame.APIRouter

	QuestionAPI question.QuestionAPI
	ProfileAPI  profile.ProfileAPI
}

func NewAPI() API {
	return API{
		Router:      frame.NewAPIRouter(),
		QuestionAPI: question.NewQuestionAPI(),
		ProfileAPI:  profile.NewProfileAPI(),
	}
}

func (ap *API) RegisterPath() {
	ap.Router.SetPrefix("api")

	ap.Router.RegisterMidddleware(middlewares.NewCheckContentTypeMiddleware())

	// Routes for profile
	ap.Router.GET("profile", ap.ProfileAPI.GetProfile, []string{})

	// Routes for questions
	ap.Router.POST("question/make", ap.QuestionAPI.MakeQuestion, []string{"ContentTypeCheck"})
	ap.Router.GET("questions", ap.QuestionAPI.GetRecentQuestions, []string{})
	ap.Router.GET("questions/{id}", ap.QuestionAPI.GetQuestion, []string{})
	ap.Router.POST("questions/answer", ap.QuestionAPI.AnswerQuestion, []string{"ContentTypeCheck"})
}
