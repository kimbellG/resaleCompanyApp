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

	"cw/provider"
	provider_dlvr "cw/provider/delivery"
	prv_rep "cw/provider/postgres"
	prv_usecase "cw/provider/usecase"

	"cw/client"
	client_dlvr "cw/client/delivery"
	client_rep "cw/client/postgres"
	client_usecase "cw/client/usecase"

	"cw/app"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type App struct {
	authController auth.UseCase
	provider       provider.UseCase
	clt            client.UseCase
}

func NewApp() *App {
	db := initDB()

	userRepo := postgres.NewUserRepository(db)
	providerRepo := prv_rep.NewProviderRepository(db)
	cltRepo := client_rep.NewClientRepository(db)

	return &App{
		authController: usecase.NewAuthUseCase(
			userRepo,
			[]byte(os.Getenv("KEYPASSWORD"))),
		provider: prv_usecase.NewProviderUseCase(providerRepo),
		clt:      client_usecase.NewClientUseCase(cltRepo),
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
	provider_dlvr.RegisterEndpoints(router, a.provider)
	client_dlvr.RegisterEndpoints(router, a.clt)

	router.Use(app.JWTAuthentication)
	router.Use(app.CheckAccessRight)
	router.Use(app.LogNewConnection)

	log.Fatalln(http.ListenAndServe(":"+port, router))
}
