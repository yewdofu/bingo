package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"html"
	"io/fs"
	"log"
	"net/http"

	"smw-bingo/bingo"
)

//go:embed static/*
var staticFiles embed.FS

func main() {
	bingo.InitData("./bingo.json")

	// Serve static files
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/", http.FileServer(http.FS(staticFS)))

	http.HandleFunc("/api/test", handleTest)
	http.HandleFunc("/api/create", handleCreateBingo)
	http.HandleFunc("/api/debug", handleDebug)

	log.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
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
