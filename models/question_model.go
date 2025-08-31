package models

import (
	"database/sql"
	"questionbasket/frame"
	"time"

	"github.com/google/uuid"
)

type Question struct {
	QID        string         `db:"qid"`
	AID        string         `db:"aid"`
	Content    string         `db:"content"`
	IsNSFW     bool           `db:"is_nsfw"`
	CreatedAt  time.Time      `db:"created_at"`
	Answer     sql.NullString `db:"answer"`
	ShareRange sql.NullInt64  `db:"share_range"`
	AnsweredAt sql.NullTime   `db:"answered_at"`
}

type QuestionModel struct {
	db frame.DatabaseConnector
}

func (m *QuestionModel) DatabaseConnect(dbc frame.DatabaseConnector) {
	m.db = dbc
}

func (m *QuestionModel) GetModelInfo() frame.ModelInfo {
	return frame.ModelInfo{Name: "QuestionModel"}
}

// InsertQuestion creates a new question record in the database.
func (m *QuestionModel) InsertQuestion(content string, isNSFW bool) (string, string, error) {
	newQID := uuid.New().String()
	newAID := uuid.New().String()
	query := `INSERT INTO questions (qid, aid, content, is_nsfw, created_at) VALUES (:qid, :aid, :content, :is_nsfw, :created_at)`

	newQuestion := Question{
		QID:       newQID,
		AID:       newAID,
		Content:   content,
		IsNSFW:    isNSFW,
		CreatedAt: time.Now(),
	}

	_, err := m.db.RunPrepared(query, newQuestion)
	if err != nil {
		return "", "", err
	}

	return newQID, newAID, nil
}

func (m *QuestionModel) AddAnswerByAID(aid string, content string, shareRange int) (string, error) {
	tx, err := m.db.Connector.Beginx()
	if err != nil {
		return "", err
	}
	defer tx.Rollback() // Rollback on failure

	var qidFromDB string
	err = tx.Get(&qidFromDB, "SELECT qid FROM questions WHERE aid = ?", aid)
	if err != nil {
		return "", err
	}

	updateQuery := "UPDATE questions SET answer = ?, share_range = ?, answered_at = ? WHERE aid = ?"
	_, err = tx.Exec(updateQuery, content, shareRange, time.Now(), aid)
	if err != nil {
		return "", err
	}

	err = tx.Commit()
	if err != nil {
		return "", err
	}

	return qidFromDB, nil
}

func (m *QuestionModel) GetQuestion(qid string) (*Question, error) {
	query := "SELECT * FROM questions WHERE qid = :qid"

	params := map[string]interface{}{
		"qid": qid,
	}

	row, err := m.db.RunPreparedRow(query, params)
	if err != nil {
		return nil, err
	}

	var question Question
	err = row.StructScan(&question)
	if err != nil {
		return nil, err
	}

	return &question, nil
}

func (m *QuestionModel) GetRecentQuestions() ([]Question, error) {
	query := "SELECT * FROM questions ORDER BY created_at DESC LIMIT 20"
	var questions []Question
	err := m.db.RunPreparedSelect(&questions, query, nil)
	if err != nil {
		return nil, err
	}
	return questions, nil
}
