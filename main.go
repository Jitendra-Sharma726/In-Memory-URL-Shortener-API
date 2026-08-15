package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
)

// 1. Our simple in-memory database
var urlStore = make(map[string]string)

// 2. Structs to define our JSON input and output
type ShortenRequest struct {
	OriginalURL string `json:"originalUrl"`
}

type ShortenResponse struct {
	ShortCode string `json:"shortCode"`
}


// generateCode creates a random 6-character string
func generateCode() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	
	// Create an empty array of 6 bytes
	code := make([]byte, 6)
	
	// Loop 6 times, picking a random byte from our charset each time
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	
	// Convert the final byte array back into a printable string
	return string(code)
}


// shortenHandler takes a long URL and creates a short code
func shortenHandler(w http.ResponseWriter, r *http.Request) {
	var req ShortenRequest
	
	// 1. Create a new decoder that reads from the incoming request body
	decoder := json.NewDecoder(r.Body)

	// 2. Tell the decoder to translate the JSON and inject it into our 'req' struct
	err := decoder.Decode(&req)

	// 3. Check if something went wrong during the translation
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if req.OriginalURL == "" {
		http.Error(w, "originalUrl is required", http.StatusBadRequest)
		return
	}

	// Generate a random 6-character code
	code := generateCode()

	// Save the mapping to our in-memory dictionary
	urlStore[code] = req.OriginalURL

	// Prepare the response data
	res := ShortenResponse{
		ShortCode: code,
	}
	
	// Set the headers and success status code
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	
	// 4. Create an encoder and send the JSON back to the user
	encoder := json.NewEncoder(w)
	encoder.Encode(res)
}


// redirectHandler catches the short link and sends the browser to the real website
func redirectHandler(w http.ResponseWriter, r *http.Request) {
	// Easily grab the {code} wildcard from the URL path
	code := r.PathValue("code")

	// Look up the code in our map
	longURL, exists := urlStore[code]
	if !exists {
		http.Error(w, "404 - Short URL not found", http.StatusNotFound)
		return
	}

	// Redirect the user to the original website
	http.Redirect(w, r, longURL, http.StatusFound)
}


func main() {
	// Specifying "POST" here means we don't have to write manual method checks.
	http.HandleFunc("POST /api/shorten", shortenHandler)
	
	// The {code} wildcard automatically extracts the shortcode from the URL.
	http.HandleFunc("GET /{code}", redirectHandler)

	fmt.Println("🚀 Server running on port 8080")
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Printf("Server crashed: %v\n", err)
	}
}


