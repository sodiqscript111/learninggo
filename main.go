package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	// Check that the method is POST
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Decode the JSON body
	var user User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Respond with the user data
	fmt.Fprintf(w, "Received user: Name: %s, Email: %s", user.Name, user.Email)
}

func main() {
	http.HandleFunc("/register", handlePost)
	http.ListenAndServe(":8080", nil)
}
