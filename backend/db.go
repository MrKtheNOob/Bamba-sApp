package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type DatabaseManager struct {
	*sql.DB
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
	return &DatabaseManager{db}, nil
}

func (*DatabaseManager) PingDB(db *DatabaseManager) error {
	return db.Ping()
}
func (db *DatabaseManager) InsertUser(username string, password string) error {
	query := "INSERT INTO users (username,password,score) VALUES (?, ?, ?)"
	_, err := db.Exec(query, username, password)
	return err
}

// GetUserByUsername retrieves a user from the database by username
func (db *DatabaseManager) GetUserByUsername(username string) (*User, error) {
	query := "SELECT username, password, score FROM users WHERE username = ?"
	row := db.QueryRow(query, username)

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
