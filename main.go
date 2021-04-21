package main

import (
	"fmt"
	"net/http"
	"os"

	"cw/app"
	"cw/controllers"
	"cw/models"

	neasted "github.com/antonfisher/nested-logrus-formatter"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&neasted.Formatter{
		HideKeys: true,
	})
	log.SetOutput(os.Stdout)
}

func main() {

	router := mux.NewRouter()
	db := models.CreateNewDBConnection()
	env := &controllers.Env{db, db}

	router.HandleFunc("/", controllers.TestingToken)
	router.HandleFunc("/register", env.RegistrationHandler)
	router.HandleFunc("/login", env.PasswordAuthentification)

	router.Use(app.JWTAuthentication)
	router.Use(app.LogNewConnection)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Info(port)
	fmt.Println(port)

	err := http.ListenAndServe(":"+port, router)

	if err != nil {
		log.Fatalln(err)
	}

}
