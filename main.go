package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func init() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Could not load .env file!")
	}
}

func main() {

	fmt.Println("Hello!")
}
