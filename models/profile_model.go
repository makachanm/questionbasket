package models

import "questionbasket/frame"

type Profile struct {
	Name        string `db:"name"`
	Description string `db:"description"`
	BasedOn     int    `db:"based_on"`
}

type ProfileModel struct {
	db frame.DatabaseConnector
}

func (m *ProfileModel) DatabaseConnect(dbc frame.DatabaseConnector) {
	m.db = dbc
}

func (m *ProfileModel) GetModelInfo() frame.ModelInfo {
	return frame.ModelInfo{Name: "ProfileModel"}
}

func (m *ProfileModel) GetProfile() (*Profile, error) {
	profile := Profile{}
	// Assuming a single row in the profile table
	err := m.db.Connector.Get(&profile, "SELECT * FROM profile LIMIT 1")
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (m *ProfileModel) InsertProfile(name, description string, basedOn int) error {
	query := "INSERT INTO profile (name, description, based_on) VALUES (?, ?, ?)"
	_, err := m.db.Connector.Exec(query, name, description, basedOn)
	return err
}
