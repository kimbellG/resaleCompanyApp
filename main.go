package main

import (
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

	port := os.Getenv("PORT")
	log.Info(port)
	if port == "" {
		port = "8000"
	}

	router := mux.NewRouter()
	db := models.CreateNewDBConnection()
	env := &controllers.Env{db, db, db}

	router.HandleFunc("/", controllers.TestingToken)
	router.HandleFunc("/register", env.RegistrationHandler)
	router.HandleFunc("/login", env.PasswordAuthentification)
	router.HandleFunc("/provider/create", env.CreateProviderController)

	router.Use(app.JWTAuthentication)
	router.Use(app.LogNewConnection)

	err := http.ListenAndServe(":"+port, router)

	if err != nil {
		log.Fatalln(err)
	}

}
