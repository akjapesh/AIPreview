package main

import (
	"backend/database"
	"context"
	"log"
)

func main() {

	//connect to database
	conn, err := database.ConnectToDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer conn.Close(context.Background())

	//connect to graphqi
	server(conn)

}
