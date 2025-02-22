package main

import (
	"AIcodegenerator/database"
	"context"
	"fmt"
)

func main() {

	//Fetch data from a Supabase table (replace "your_table" with the actual table name)

	database.FetchDataUsingAPIKey("users")

	conn, err := database.ConnectToDatabase()
	if err != nil {
		fmt.Println(" Error:", err)
		return
	}
	defer conn.Close(context.Background())
}
