package main

import (
	"encoding/json"
	"net/http"
)

type Task struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Details string `json:"details"`
}

var tasks = []Task{}
var nextID = 1

func handleTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Post only", http.StatusMethodNotAllowed)
	}

	var task Task
	task.ID = nextID
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	tasks = append(tasks, task)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(task)
}

func handelGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Get only", http.StatusMethodNotAllowed)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}
func main() {

	http.HandleFunc("/", handleTask)
	http.HandleFunc("/get", handelGet)
	http.ListenAndServe(":8080", nil)
}
