package main

import (
	"bufio"
	"fmt"
	"os"
	"questionbasket/frame"
	"questionbasket/models"
	"strings"
)

func DoInitSetup() {
	fmt.Println("Making Database Tables...")
	db := frame.GetDB()

	questionsSchema := `
	CREATE TABLE IF NOT EXISTS questions (
		"qid" TEXT PRIMARY KEY,
		"aid" TEXT NOT NULL UNIQUE,
		"content" TEXT NOT NULL,
		"is_nsfw" INTEGER NOT NULL,
		"created_at" DATETIME NOT NULL,
		"answer" TEXT,
		"share_range" INTEGER,
		"answered_at" DATETIME
	);
	`

	profileSchema := `
	CREATE TABLE IF NOT EXISTS profile (
		"name" TEXT NOT NULL,
		"description" TEXT NOT NULL,
		"based_on" INTEGER NOT NULL
	);
	`

	_, err := db.Exec(questionsSchema)
	if err != nil {
		fmt.Printf("Error creating questions table: %v\n", err)
		return
	}

	_, err = db.Exec(profileSchema)
	if err != nil {
		fmt.Printf("Error creating profile table: %v\n", err)
		return
	}
	fmt.Println("Tables created successfully.")

	fmt.Println("\n--- Profile Setup ---")
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter your nickname: ")
	nickname, _ := reader.ReadString('\n')
	nickname = strings.TrimSpace(nickname)

	fmt.Print("Enter your description: ")
	description, _ := reader.ReadString('\n')
	description = strings.TrimSpace(description)

	profileModel := models.ProfileModel{}
	frame.DatabaseBind(&profileModel)

	// The 'based_on' field is '1' in the api.md example. Hardcoding this for now.
	err = profileModel.InsertProfile(nickname, description, 1)
	if err != nil {
		fmt.Printf("Error saving profile: %v\n", err)
		return
	}

	fmt.Println("\nProfile setup complete. You can now run the application normally.")
}