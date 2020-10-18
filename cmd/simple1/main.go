package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprintf(w, "hello")
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8000", nil)
}
