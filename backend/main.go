package main

import (
	"backend/database"
	"log"
)

func main() {

	err := database.ConnectToRedis()
	if err != nil {
		log.Fatalf("Redis connection error: %v", err)
	}

	pool, err := database.ConnectToDatabaseByPool()
	if err != nil {
		log.Fatalf("Failed to connect to database through pool: %v", err)
	}
	defer pool.Close()

	//connect to graphqi
	server(pool)

}
