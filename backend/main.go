package main

import (
	"backend/database"
	"context"
	"log"
	"backend/pkg/mlServices"
	"os"
	"time"
)

func main() {

	//connect to database
	conn, err := database.ConnectToDatabase()

	// Initialize ML client
	mlClient := mlServices.NewMLClient(os.Getenv("ML_SERVICE_URL"),5,30*time.Second)

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer conn.Close(context.Background())

	//connect to graphqi
	server(conn,mlClient)

}
