package server

import (
	"fmt"
	"net/http"
	"os"

	"database/sql"

	_ "github.com/jackc/pgx/stdlib"

	"cw/auth"
	auth_delivery "cw/auth/delivery"
	"cw/auth/postgres"
	"cw/auth/usecase"

	"cw/app"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type App struct {
	authController auth.UseCase
}

func NewApp() *App {
	db := initDB()

	userRepo := postgres.NewUserRepository(db)

	return &App{
		authController: usecase.NewAuthUseCase(
			userRepo,
			[]byte(os.Getenv("KEYPASSWORD"))),
	}
}

func initDB() *sql.DB {
	lib_db, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(fmt.Errorf("open db is failed: %v", err))
	}

	if err := lib_db.Ping(); err != nil {
		panic(fmt.Errorf("Pinging db is failed: %v", err))
	}

	return lib_db
}

func (a *App) Run(port string) {
	router := mux.NewRouter()

	auth_delivery.RegisterEndpoints(router, a.authController)

	router.Use(app.JWTAuthentication)
	router.Use(app.CheckAccessRight)
	router.Use(app.LogNewConnection)

	log.Fatalln(http.ListenAndServe(":"+port, router))
}
