package main

import (
	"fmt"
	"github.com/IgnasiBosch/go-jwt-auth/api"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env: %v", err)
	} else {
		fmt.Println("Getting env variables")
	}

	api.Run()
}
