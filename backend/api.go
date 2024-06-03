package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func loadPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Page loaded")
	http.ServeFile(w, r, "../frontend/index.html")
}
func pickQuestion() QuizItem {
	rand.Seed(time.Now().UnixNano())

	// Read the JSON file
	data, err := ioutil.ReadFile("quiz_data.json")
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
func main() {

	http.HandleFunc("/", loadPage)
	http.HandleFunc("/getquestion", giveQuestion)
	fmt.Println("Server listening on 192.168.1.20:8080")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := http.ListenAndServe("0.0.0.0:"+port, nil); err != nil {
		log.Panic("error:%s", err)
	}
}
