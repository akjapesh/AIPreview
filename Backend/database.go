package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func loadEnvFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || !strings.Contains(line, "=") {
			continue // Skip comments and invalid lines
		}

		parts := strings.SplitN(line, "=", 2)
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		os.Setenv(key, value) // Set environment variable
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("Error reading .env file:", err)
	}
}

func FetchData(tableName string) {
	loadEnvFile(".env") // Load environment variables from .env file
	dbUrl := os.Getenv("DATABASE_URL")
	dbAPIKey := os.Getenv("DATABASE_API_KEY")

	if dbUrl == "" || dbAPIKey == "" {
		log.Fatal("Please set DATABASE_API_KEY and DATABASE_API_KEY environment variables")
	}

	// Fetch data from database
	apiURL := fmt.Sprintf("%s/rest/v1/%s?select=*", dbUrl, tableName)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("apikey", dbAPIKey)
	req.Header.Set("Authorization", "Bearer "+dbAPIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to fetch data: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response: %v", err)
	}

	// Print the response
	fmt.Println("Fetched Data:", string(body))

}
