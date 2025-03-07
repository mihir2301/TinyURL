package main

import (
	"fmt"
	"log"
	"os"
	"tinyurl/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error in loading env files")
		return
	}
	r := gin.Default()
	routes.Routes(r)
	log.Fatal(r.Run(os.Getenv("APP_PORT")))
}
