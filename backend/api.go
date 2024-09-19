package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type QuizItem struct {
	Question string   `json:"question"`
	Options  []string `json:"options"`
	Answer   string   `json:"answer"`
	Meaning  string   `json:"meaning"`
}

var ErrUserAlreadyExists = errors.New("user already exists")

func loadPage(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Page loaded for:%v\n", r.RemoteAddr)
	tmpl := template.Must(template.ParseFiles(filepath.Join("..", "frontend", "index.html")))

	// Render the index.html template
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func vocabPage(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Vocab Page loaded for:%v\n", r.RemoteAddr)
	tmpl := template.Must(template.ParseFiles(filepath.Join("..", "frontend", "vocabulary.html")))

	// Render the index.html template
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func leaderboardPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("leaderboard page loaded")
	tmpl := template.Must(template.ParseFiles(filepath.Join("..", "frontend", "leaderboard.html")))

	// Render the index.html template
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func login(w http.ResponseWriter, r *http.Request, db *DatabaseManager) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	username := r.Form.Get("Name")
	password := r.Form.Get("Password")
	score := r.Form.Get("score")
	intscore, err := strconv.Atoi(score)
	if err != nil {
		fmt.Println(intscore)
		http.Error(w, "The score must be a number", http.StatusBadRequest)
		return
	}
	if username == "" || password == "" {
		http.Error(w, "Empty input field", http.StatusBadRequest)
	}

	user, err := db.GetUserByUsernameAndPassword(username, password)
	if err != nil {
		fmt.Printf("%v tried to login but failed finding himself in the DB\n", r.RemoteAddr)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if user == nil {
		fmt.Printf("%v tried to login with invalid credentials\n", r.RemoteAddr)
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Update user score
	err = db.UpdateUserScore(db, user, intscore)
	if err != nil {
		fmt.Printf("Error updating score for user %v: %v\n", r.RemoteAddr, err)
		http.Error(w, "Error updating score", http.StatusInternalServerError)
		return
	}

	// Log success and respond
	fmt.Printf("User %v logged in successfully and updated their score\n", r.RemoteAddr)
	w.WriteHeader(http.StatusOK) // Ensure only one WriteHeader call
	w.Write([]byte("Login successful"))
}
func register(w http.ResponseWriter, r *http.Request, db *DatabaseManager) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	username := r.Form.Get("Name")
	password := r.Form.Get("Password")
	cpassword := r.Form.Get("CPassword")
	if password != cpassword {
		http.Error(w, "password and confirm password mismatch", http.StatusBadRequest)
		log.Printf("%v tried to register but mismatched password and confirm password\n", r.RemoteAddr)
		return
	}
	if err := db.InsertUser(username, password); err != nil {
		if err == ErrUserAlreadyExists {
			http.Error(w, err.Error(), http.StatusConflict)
		}
		log.Println(err.Error())
		http.Error(w, "Database Error", http.StatusInternalServerError)
		log.Println("Database Error for register")
		return
	}
	log.Printf("New user of name %v registered\n", username)
}
func GiveLeaderboardData(w http.ResponseWriter, r *http.Request, db *DatabaseManager) {
	leaderboardData, err := db.GetLeaderboardData()
	if err != nil {
		http.Error(w, "Database error fetching leaderboard data;error is:"+err.Error(), http.StatusInternalServerError)
		return
	}
	response, err := json.Marshal(leaderboardData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
func pickQuestion(subject string) QuizItem {
	rand.NewSource(time.Now().UnixNano())
	var data []byte
	var err error
	// Read the JSON file
	switch subject {
	case "IV":
		data, err = os.ReadFile("quiz_data.json")
		if err != nil {
			log.Fatalf("Failed to read JSON file: %v", err)
		}

	case "vocab":
		data, err = os.ReadFile("vocab_data.json")
		if err != nil {
			log.Fatalf("Failed to read JSON file: %v", err)
		}
	}
	// Unmarshal the JSON data into a slice of QuizItem
	var quizItems []QuizItem
	err = json.Unmarshal(data, &quizItems)
	if err != nil {
		log.Fatalf("Failed to unmarshal JSON data: %v", err)
	}
	// Select a random index of the json array
	randomIndex := rand.Intn(len(quizItems))
	selectedQuizItem := quizItems[randomIndex]

	return *randomizeOptions(&selectedQuizItem)
}
func randomizeOptions(selectedQuizItem *QuizItem) *QuizItem {
	//Randomizes the answers so that they are not at the same index all the time
	for {
		var index1 = rand.Intn(3)
		var index2 = rand.Intn(3)
		var index3 = rand.Intn(3)
		if index1 != index2 && index1 != index3 && index2 != index3 {
			selectedQuizItem.Options[index1], selectedQuizItem.Options[index2] = selectedQuizItem.Options[index2], selectedQuizItem.Options[index1]
			selectedQuizItem.Options[index3], selectedQuizItem.Options[index1] = selectedQuizItem.Options[index1], selectedQuizItem.Options[index3]
			fmt.Printf("Question: %s\n", selectedQuizItem.Question)
			for i, option := range selectedQuizItem.Options {
				fmt.Printf("%d. %s\n", i+1, option)
			}
			fmt.Println(selectedQuizItem.Meaning)
			return selectedQuizItem
		}
	}
}
func giveQuestion(w http.ResponseWriter, r *http.Request, subject string) {
	fmt.Printf("Getting %s question for:%v\n", subject, r.RemoteAddr)

	selectedQuizItem := pickQuestion(subject)
	response, err := json.Marshal(selectedQuizItem)
	if err != nil {
		http.Error(w, "Failed to marshal JSON response", http.StatusInternalServerError)
		return
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response
	w.Write(response)
}
func handleFeedback(w http.ResponseWriter, r *http.Request, db *DatabaseManager) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	answer := r.Form.Get("answer")
	suggestion := r.Form.Get("suggestion")
	if answer == "" {
		http.Error(w, "Please select an option for the answer.", http.StatusBadRequest)
	}
	// I'm gonna have to refactor this to something more serious
	fmt.Println("FEEDBACK:{answer:", answer, "suggestion:", string(suggestion), "}")
	if err = db.InsertFeedback(answer, suggestion); err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	fmt.Println("Feedback submitted successfully")
	w.Write([]byte("Feedback submitted successfully"))
}
func main() {
	fs := http.FileServer(http.Dir("../frontend"))
	http.Handle("GET /static/", http.StripPrefix("/static/", fs))

	dsn := os.Getenv("DB_URI")
	if dsn == "" {
		log.Fatal("db uri not found")
	}
	db, err := InitialiseDB(dsn)
	if err != nil {
		log.Fatal(err)
	}
	//the routes
	http.HandleFunc("GET /", loadPage)
	http.HandleFunc("GET /vocab", vocabPage)
	http.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) {
		login(w, r, db)
	})
	http.HandleFunc("POST /register", func(w http.ResponseWriter, r *http.Request) {
		register(w, r, db)
	})
	http.HandleFunc("POST /feedback", func(w http.ResponseWriter, r *http.Request) {
		handleFeedback(w, r, db)
	})
	http.HandleFunc("GET /leaderboard", leaderboardPage)
	http.HandleFunc("GET /getivquestion", func(w http.ResponseWriter, r *http.Request) {
		giveQuestion(w, r, "IV")
	})
	http.HandleFunc("GET /getvocabquestion", func(w http.ResponseWriter, r *http.Request) {
		giveQuestion(w, r, "vocab")
	})
	http.HandleFunc("GET /getleaderboarddata", func(w http.ResponseWriter, r *http.Request) {
		GiveLeaderboardData(w, r, db)
	})
	http.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {})
	fmt.Println("Server Started ...")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := http.ListenAndServe("0.0.0.0:"+port, nil); err != nil {
		log.Printf("error:%s", err)
	}
}
