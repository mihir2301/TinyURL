package main

import (
	"fmt"
	"tinyurl/routes"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error in loading env files")
		return
	}
	routes.Client()
}
