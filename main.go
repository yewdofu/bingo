package main

import (
	"bingo/bingo"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
)

func main() {
	bingo.InitData("./bingo.json")

	http.HandleFunc("/test", handleTest)
	http.HandleFunc("/create", handleCreateBingo)
	http.HandleFunc("/debug", handleDebug)

	log.Fatal(http.ListenAndServe(":8080", nil))
	log.Println("Server started at http://localhost:8080")
}

func handleTest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func handleDebug(w http.ResponseWriter, r *http.Request) {
	bingo.InitData("./bingo.json")
}

func handleCreateBingo(w http.ResponseWriter, r *http.Request) {
	log.Println("Create bingo card.")

	seed := r.FormValue("seed")
	bc, err := bingo.CreateBingoCard(seed)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bc)
}
