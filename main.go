package main

import (
	"cw/logger"
	"cw/server"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println(err)
	}

}

func main() {
	mainLogger := logger.NewLoggerWithFields(
		map[string]interface{}{"action": "start server"})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	mainLogger.Infof("port: %v", port)

	app := server.NewApp()
	app.Run(port)
}
