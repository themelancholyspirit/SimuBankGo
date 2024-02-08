package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := NewPosgreDB()

	if err != nil {
		log.Fatal(err)
	}

	if err := db.createAccountTable(); err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("PORT")

	fmt.Println("Server is running successfully on port:", port)

	s := NewAPIServer(fmt.Sprintf(":%s", port), db)

	s.Run()

}
