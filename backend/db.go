package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type DatabaseManager struct {
	*sql.DB
}
type PlayerData struct {
	Username string `json:"username"`
	Score    int    `json:"score"`
}
type User struct {
	username string
	password string
	score    string
}
type Feedback struct {
	id         int
	likes      string
	suggestion string
}

func InitialiseDB(dsn string) (*DatabaseManager, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	fmt.Println("Database connceted successfully !")
	return &DatabaseManager{db}, nil
}

// Ping the database
func (*DatabaseManager) PingDB(db *DatabaseManager) error {
	return db.Ping()
}

// InserUser registers a new user into the db
func (db *DatabaseManager) InsertUser(username string, password string) error {
	user, _ := db.GetUserByUsernameAndPassword(username, password)
	if user != nil {
		return ErrUserAlreadyExists
	}
	query := "INSERT INTO users (username,password,score) VALUES (?, ?, ?)"
	_, err := db.Exec(query, username, password, 0)
	return err
}

// Update user's highest score
func (*DatabaseManager) UpdateUserScore(db *DatabaseManager, user *User, score int) error {
	query := "SELECT score FROM users WHERE username = ? AND password = ?;"
	var highestScore int
	err := db.QueryRow(query, user.username, user.password).Scan(&highestScore)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}

	if score <= highestScore {
		return nil
	}

	query = "UPDATE users SET score = ? WHERE username = ? AND password = ?"
	_, err = db.Exec(query, score, user.username, user.password)
	if err != nil {
		return err
	}

	fmt.Println("Update made successfully")
	return nil
}

// GetUserByUsername retrieves a user from the database by username
func (db *DatabaseManager) GetUserByUsernameAndPassword(username string, password string) (*User, error) {
	query := "SELECT username, password, score FROM users WHERE username = ? AND password= ?"
	row := db.QueryRow(query, username, password)

	var user User
	err := row.Scan(&user.username, &user.password, &user.score)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUsers method to retrieve users from the database
func (db *DatabaseManager) GetAllUsers() ([]User, error) {
	query := "SELECT username,password,score FROM users" //possible comma error
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.username, &user.password, &user.score); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}
func (db *DatabaseManager) GetLeaderboardData() ([]PlayerData, error) {
	rows, err := db.Query("SELECT username, score FROM users ORDER BY score DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var leaderboard []PlayerData
	for rows.Next() {
		var player PlayerData
		if err := rows.Scan(&player.Username, &player.Score); err != nil {
			return nil, err
		}
		leaderboard = append(leaderboard, player)
	}
	return leaderboard, nil
}
func (db *DatabaseManager) InsertFeedback(likes string, suggestion string) error {
	query := "INSERT INTO feedback (likes, suggestion) VALUES (?, ?)"
	_, err := db.Exec(query, likes, suggestion)
	return err
}

// GetAllFeedback retrieves all feedback entries from the database
func (db *DatabaseManager) GetAllFeedback() ([]Feedback, error) {
	query := "SELECT id, likes, suggestion FROM feedback"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feedbacks []Feedback
	for rows.Next() {
		var feedback Feedback
		if err := rows.Scan(&feedback.id, &feedback.likes, &feedback.suggestion); err != nil {
			return nil, err
		}
		feedbacks = append(feedbacks, feedback)
	}

	return feedbacks, nil
}
