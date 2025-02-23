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

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
)

func ConnectToDatabase() (*pgx.Conn, error) {
	// Load environment variables
	LoadEnvFile(".env")

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

func LoadEnvFile(filename string) {
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
	LoadEnvFile(".env") // Load environment variables from .env file
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

// FetchUserByEmail fetches a user by their email address.
func FetchUserByEmail(conn *pgx.Conn, email string) (*User, error) {
	query := "SELECT id, username, email, password_hash, created_at, updated_at, is_active FROM users WHERE email = $1"

	// Execute the query and fetch a single row
	row := conn.QueryRow(context.Background(), query, email)

	var user User
	// Scan the row into the User struct
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt, &user.IsActive)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	return &user, nil
}

// UserExists checks if a user with the given email or username exists.
func UserExists(conn *pgx.Conn, email, username string) (bool, error) {
	query := "SELECT id FROM users WHERE email = $1 OR username = $2"
	row := conn.QueryRow(context.Background(), query, email, username)

	var existingUserID uuid.UUID
	err := row.Scan(&existingUserID)
	if err == nil {
		return true, nil // User exists
	} else if err != pgx.ErrNoRows {
		return false, fmt.Errorf("error checking existing user: %w", err)
	}

	return false, nil // User does not exist
}

// SaveUser inserts a new user into the database.
func SaveUser(conn *pgx.Conn, newUser User) (*User, error) {
	insertQuery := `
		INSERT INTO users (username, email, password_hash, created_at, updated_at, is_active) 
		VALUES ($1, $2, $3, $4, $5, $6) 
		RETURNING id, created_at, updated_at, is_active`

	err := conn.QueryRow(context.Background(), insertQuery,
		newUser.Username, newUser.Email, newUser.PasswordHash,
		newUser.CreatedAt, newUser.UpdatedAt, newUser.IsActive,
	).Scan(&newUser.ID, &newUser.CreatedAt, &newUser.UpdatedAt, &newUser.IsActive)

	if err != nil {
		return nil, fmt.Errorf("failed to insert new user: %w", err)
	}

	return &newUser, nil
}
