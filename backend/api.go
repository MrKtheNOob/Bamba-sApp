package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
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
	http.ServeFile(w, r, "../frontend/index.html")
}
func leaderboardPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("leaderboard page loaded")
	http.ServeFile(w, r, "../frontend/leaderboard.html")
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
		fmt.Printf("%v tried to login but failed finding themselves in the DB\n", r.RemoteAddr)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if user == nil {
		fmt.Printf("%v tried to login with invalid credentials\n", r.RemoteAddr)
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	err = db.UpdateUserScore(db, user, intscore)
	if err != nil {
		fmt.Printf("Error updating score for user %v: %v\n", r.RemoteAddr, err)
		http.Error(w, "Error updating score", http.StatusInternalServerError)
		return
	}

	// Log success and respond
	log.Printf("User %v logged in successfully and updated their score\n", r.RemoteAddr)
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
		log.Printf("Database Error for register")
		return
	}
	fmt.Printf("New user of name %v registered\n", username)
}
func pickQuestion() QuizItem {
	rand.NewSource(time.Now().UnixNano())
	// Read the JSON file
	data, err := os.ReadFile("quiz_data.json")
	if err != nil {
		log.Fatalf("Failed to read JSON file: %v", err)
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
func giveQuestion(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Getting question for:%v\n", r.RemoteAddr)
	selectedQuizItem := pickQuestion()
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
	http.HandleFunc("GET /getquestion", giveQuestion)
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
