package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

// TestHandler serve as a testing endpoint
func TestHandler(w http.ResponseWriter, r *http.Request) {
	// Log the headers
	//for name, values := range r.Header {
	//	for _, value := range values {
	//		log.Printf("%s: %s", name, value)
	//	}
	//}

	// Calculate and log total header size
	totalSize := calculateHeaderSize(r.Header)
	// Gather additional information
	method := r.Method
	url := r.URL.String()
	ip := r.RemoteAddr
	userAgent := r.UserAgent()
	timestamp := time.Now().Format(time.RFC3339)

	// Log the information
	log.Printf("Request received: %s %s from %s (User-Agent: %s) at %s", method, url, ip, userAgent, timestamp)
	log.Printf("Total header size: %d bytes", totalSize)

	// Prepare the response
	response := map[string]string{
		"message":           "This is a test endpoint for middleware.",
		"total header size": strconv.Itoa(totalSize),
		"method":            method,
		"url":               url,
		"client_ip":         ip,
		"user_agent":        userAgent,
		"timestamp":         timestamp,
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}

func calculateHeaderSize(headers http.Header) int {
	totalSize := 0
	for name, values := range headers {
		totalSize += len(name) + 2 // for ": " (colon and space)
		for _, value := range values {
			totalSize += len(value)
		}
	}
	return totalSize
}
