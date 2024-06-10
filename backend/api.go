package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type QuizItem struct {
	Question string   `json:"question"`
	Options  []string `json:"options"`
	Answer   string   `json:"answer"`
}

// type Feedback struct {
// 	Answer     string `gorm:"not null"`
// 	Suggestion string `gorm:"type:text"`
// }

func loadPage(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Page loaded for:%v", r.RemoteAddr)
	http.ServeFile(w, r, "../frontend/index.html")
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

	// Select a random index
	randomIndex := rand.Intn(len(quizItems))

	// Print the selected quiz item
	selectedQuizItem := quizItems[randomIndex]
	fmt.Printf("Question: %s\n", selectedQuizItem.Question)
	for i, option := range selectedQuizItem.Options {
		fmt.Printf("%d. %s\n", i+1, option)
	}
	return selectedQuizItem
}
func giveQuestion(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Getting question for:%v", r.RemoteAddr)
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
func handleFeedback(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	answer := r.Form.Get("answer")
	suggestion := r.Form.Get("reason")

	if answer == "" {
		http.Error(w, "Please select an option for the answer.", http.StatusBadRequest)
	}

	fmt.Println("FEEDBACK:{answer:", answer, "suggestion:", suggestion, "}")
	//put in db
	fmt.Fprintln(w, "Thank you so much 🙏!")
}
func main() {
	http.HandleFunc("POST /feedback", handleFeedback)
	http.HandleFunc("GET /", loadPage)
	http.HandleFunc("GET/getquestion", giveQuestion)
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
