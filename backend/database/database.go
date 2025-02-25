package database

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
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

func ConnectToDatabaseByPool() (*pgxpool.Pool, error) {
	// Load environment variables
	LoadEnvFile(".env")

	// Get database connection string from environment variable
	dbURL := os.Getenv("DB_CONNECTION_STRING")
	if dbURL == "" {
		return nil, fmt.Errorf("DB_CONNECTION_STRING environment variable is not set")
	}

	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	maxConns := 5 // Default value
	if maxConnsEnv := os.Getenv("DB_MAX_CONNECTIONS"); maxConnsEnv != "" {
		if parsedMaxConns, err := strconv.Atoi(maxConnsEnv); err == nil {
			maxConns = parsedMaxConns
		}
	}

	config.MaxConns = int32(maxConns) // Apply max connections
	config.MinConns = 2               // Keep at least 2 open connections

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("unable to create database connection pool: %w", err)
	}

	fmt.Printf("Successfully connected to the database through pool with %v connections!\n", maxConns)
	return pool, nil
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
func FetchUserByEmail(conn *pgxpool.Pool, email string) (*User, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("user:email:%s", email)

	// Try to get the cached user data from Redis
	cachedUser, err := RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		// Cache hit: Unmarshal the JSON data into the User struct
		var user User
		if err := json.Unmarshal([]byte(cachedUser), &user); err == nil {
			return &user, nil
		}
		log.Println("Failed to unmarshal Redis data:", err) // Log if unmarshalling fails
	} else if err != redis.Nil {
		log.Println("Redis error:", err) // Log Redis errors (except key not found)
	}
	query := "SELECT id, username, email, password_hash, created_at, updated_at, is_active FROM users WHERE email = $1"

	// Execute the query and fetch a single row
	row := conn.QueryRow(context.Background(), query, email)

	var user User
	// Scan the row into the User struct
	err = row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt, &user.IsActive)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	// Store the fetched user data in Redis (cache for 10 minutes)
	userJSON, _ := json.Marshal(user) // Convert user struct to JSON
	RedisClient.Set(ctx, cacheKey, userJSON, 10*time.Minute)

	return &user, nil
}

// UserExists checks if a user with the given email or username exists.
func UserExists(conn *pgxpool.Pool, email, username string) (bool, error) {
	ctx := context.Background()

	// Create a Redis cache key based on email and username
	cacheKey := fmt.Sprintf("user:exists:%s:%s", email, username)
	cachedResult, err := RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		return cachedResult == "true", nil
	} else if err != redis.Nil {
		log.Println("Redis error:", err)
	}

	query := "SELECT id FROM users WHERE email = $1 OR username = $2"
	row := conn.QueryRow(context.Background(), query, email, username)

	var existingUserID uuid.UUID
	err = row.Scan(&existingUserID)
	if err == nil {

		RedisClient.Set(ctx, cacheKey, "true", 10*time.Minute)
		return true, nil // User exists
	} else if err != pgx.ErrNoRows {
		return false, fmt.Errorf("error checking existing user: %w", err)
	}

	return false, nil // User does not exist
}

// SaveUser inserts a new user into the database.
func SaveUser(pool *pgxpool.Pool, newUser User) (*User, error) {
	ctx := context.Background()

	// Begin a transaction
	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx) // Ensure rollback if anything goes wrong

	insertQuery := `
		INSERT INTO users (username, email, password_hash, created_at, updated_at, is_active) 
		VALUES ($1, $2, $3, $4, $5, $6) 
		RETURNING id, created_at, updated_at, is_active`

	err = tx.QueryRow(ctx, insertQuery,
		newUser.Username, newUser.Email, newUser.PasswordHash,
		newUser.CreatedAt, newUser.UpdatedAt, newUser.IsActive,
	).Scan(&newUser.ID, &newUser.CreatedAt, &newUser.UpdatedAt, &newUser.IsActive)

	if err != nil {
		return nil, fmt.Errorf("failed to insert new user: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &newUser, nil
}
