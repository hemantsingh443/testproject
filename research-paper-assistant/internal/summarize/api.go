package summarize

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// SummarizationRequest represents the request structure for the Hugging Face Inference API
type SummarizationRequest struct {
	Inputs string `json:"inputs"`
}

// SummarizationResponse represents the response structure from the Hugging Face Inference API
type SummarizationResponse struct {
	SummaryText string `json:"summary_text"`
}

// SummarizeWithAPI generates a summary using the Hugging Face Inference API
func SummarizeWithAPI(text, apiToken, modelID string) (string, error) {
	// Construct the API request
	reqBody, err := json.Marshal(SummarizationRequest{Inputs: text})
	if err != nil {
		return "", fmt.Errorf("error encoding request body: %v", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", "https://api-inference.huggingface.co/models/"+modelID, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+apiToken)
	req.Header.Set("Content-Type", "application/json")

	// Send request and get response
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Parse response
	var apiResponse SummarizationResponse
	err = json.NewDecoder(resp.Body).Decode(&apiResponse)
	if err != nil {
		return "", fmt.Errorf("error decoding response: %v", err)
	}

	return apiResponse.SummaryText, nil
}
