package main

import (
	"log"

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

	s := NewAPIServer(":8080", db)

	s.Run()

}
