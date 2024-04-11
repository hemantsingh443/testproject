package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/hemantsingh/research-paper-assistant/internal/summarize"
	"github.com/hemantsingh/research-paper-assistant/pkg/arxiv"
	"github.com/hemantsingh/research-paper-assistant/pkg/pdf"
)

// HandleIndex serves the index page
func HandleIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/index.html")
}

// HandleSearch performs a search and returns results as JSON
func HandleSearch(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("search_query")
	page, _ := strconv.Atoi(r.FormValue("page"))
	resultsPerPage := 10
	start := page * resultsPerPage

	papers, err := arxiv.SearchArxiv(query, start, resultsPerPage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"results":     papers,
		"currentPage": page,
	}

	json.NewEncoder(w).Encode(response)
}

// HandleSummarize generates a summary for a given PDF URL
func HandleSummarize(w http.ResponseWriter, r *http.Request) {
	pdfURL := r.FormValue("pdf_url")
	text, err := pdf.ExtractTextFromPDF(pdfURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	apiToken := "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" // Replace with your actual API token
	modelID := "facebook/bart-large-cnn"

	summary, err := summarize.SummarizeWithAPI(text, apiToken, modelID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, summary)
}

// HandleDownload serves the PDF file for download
func HandleDownload(w http.ResponseWriter, r *http.Request) {
	pdfURL := r.FormValue("pdf_url")

	// Get the PDF file
	resp, err := http.Get(pdfURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Create a temporary file
	file, err := os.CreateTemp("", "paper_*.pdf")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer os.Remove(file.Name())

	// Copy the PDF content to the temporary file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set headers for download
	w.Header().Set("Content-Disposition", "attachment; filename=paper.pdf")
	w.Header().Set("Content-Type", "application/pdf")

	// Serve the temporary PDF file
	http.ServeFile(w, r, file.Name())
}
