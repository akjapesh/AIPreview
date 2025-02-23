package database

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/jackc/pgx/v5"
)

func ConnectToDatabase() (*pgx.Conn, error) {
	// Load environment variables
	loadEnvFile(".env")

	// Get database connection string from environment variable
	dbURL := os.Getenv("DB_CONNECTION_STRING")
	if dbURL == "" {
		return nil, fmt.Errorf("DB_CONNECTION_STRING environment variable is not set")
	}

	conn, err := pgx.Connect(context.Background(), dbURL) // Connect to the database
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	fmt.Println("Successfully connected to the database!")
	return conn, nil
}

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

func FetchDataUsingAPIKey(tableName string) {
	loadEnvFile(".env") // Load environment variables from .env file
	dbUrl := os.Getenv("DATABASE_URL")
	dbAPIKey := os.Getenv("DATABASE_API_KEY")

	if dbUrl == "" || dbAPIKey == "" {
		log.Fatal("Please set DATABASE_API_KEY and DATABASE_API_KEY environment variables")
	}

	apiURL := fmt.Sprintf("%s/rest/v1/%s?select=*", dbUrl, tableName) // Fetch data from database table

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

	fmt.Println("Fetched Data:", string(body)) // Print the response

}

func FetchAllUsers(conn *pgx.Conn) ([]User, error) {
	query := "SELECT * FROM users"

	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close() // Close rows when function exits

	var users []User

	for rows.Next() { // Iterate over each row
		var user User
		err := rows.Scan(&user.ID, &user.Username)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		users = append(users, user) // Add user to slice
	}

	if err := rows.Err(); err != nil { // Check for iteration errors
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return users, nil // Return all users
}
