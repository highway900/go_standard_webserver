package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// User struct defines a user
type User struct {
	ID      int     `json:"ID"`
	Balance float32 `json:"Balance"`
	// automatically omitted from JSON because it is hidden/private
	accountID string
}

func handler1(w http.ResponseWriter, r *http.Request) {
	var user User
	json.NewDecoder(r.Body).Decode(&user)

	log.Printf("Got some data User: %d has a balance: %f\n",
		user.ID, user.Balance)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func handler2(w http.ResponseWriter, r *http.Request) {
	u := User{ID: 123, Balance: 43.0, accountID: "ABCDEF123"}
	log.Printf("User: %d, AccountID: %s", u.ID, u.accountID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func main() {
	http.HandleFunc("/handler1", handler1)
	http.HandleFunc("/handler2", handler2)
	http.HandleFunc("/health", healthCheckHandler)

	log.Println("Service started on :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
