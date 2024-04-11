package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hemantsingh/research-paper-assistant/api"
)

func main() {
	router := mux.NewRouter()

	// Define API routes (we'll fill these in later)
	router.HandleFunc("/", api.HandleIndex).Methods("GET")
	router.HandleFunc("/search", api.HandleSearch).Methods("POST")
	router.HandleFunc("/summarize", api.HandleSummarize).Methods("POST")
	router.HandleFunc("/download", api.HandleDownload).Methods("POST")

	// Serve static files
	fs := http.FileServer(http.Dir("static/"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	fmt.Println("Server listening on port 8000")
	http.ListenAndServe(":8000", router)
}
