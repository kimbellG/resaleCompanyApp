package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"cw/app"
	"cw/controllers"
	"cw/models"

	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()
	db := models.CreateNewDBConnection()
	env := &controllers.Env{db}

	router.HandleFunc("/", controllers.TestingToken)
	router.HandleFunc("/login", env.RegistrationHandler)

	router.Use(app.JWTAuthentication)
	router.Use(app.LogNewConnection)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	fmt.Println(port)

	err := http.ListenAndServe(":"+port, router)

	if err != nil {
		log.Fatalln(err)
	}

}
