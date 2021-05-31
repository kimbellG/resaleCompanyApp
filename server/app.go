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

	"cw/product"
	prddlvr "cw/product/delivery"
	prdrep "cw/product/postgres"
	prduse "cw/product/usecase"

	"cw/prdtoffer"
	offerdlvr "cw/prdtoffer/delivery"
	offerrep "cw/prdtoffer/postgres"
	offercase "cw/prdtoffer/usecase"

	"cw/order"
	orderdlvr "cw/order/delivery"
	orderrep "cw/order/postgres"
	ordercase "cw/order/usecase"

	"cw/rang"
	rangdlvr "cw/rang/delivery"
	rangrep "cw/rang/postgres"
	rangcase "cw/rang/usecase"

	"cw/app"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type App struct {
	authController auth.UseCase
	admin          auth.AdminUseCase
	provider       provider.UseCase
	clt            client.UseCase
	prd            product.UseCase
	offer          prdtoffer.UseCase
	ord            order.UseCase
	rank           rang.Usecase
}

func NewApp() *App {
	db := initDB()

	userRepo := postgres.NewUserRepository(db)
	adminRepo := postgres.NewAdminRepo(db)
	providerRepo := prv_rep.NewProviderRepository(db)
	cltRepo := client_rep.NewClientRepository(db)
	prdRepo := prdrep.NewProductRepository(db)
	offerRepo := offerrep.NewOfferRepository(db)
	orderRepo := orderrep.NewOrderRepository(db)
	rankRepo := rangrep.NewRangPostgres(db)

	return &App{
		authController: usecase.NewAuthUseCase(
			userRepo,
			[]byte(os.Getenv("KEYPASSWORD"))),
		provider: prv_usecase.NewProviderUseCase(providerRepo),
		clt:      client_usecase.NewClientUseCase(cltRepo),
		prd:      prduse.NewProductUseCase(prdRepo, rankRepo),
		offer:    offercase.NewProductOfferUseCase(offerRepo, providerRepo, prdRepo, rangcase.NewRangUseCase(rankRepo, userRepo)),
		ord:      ordercase.NewOrderUseCase(orderRepo, client_usecase.NewClientUseCase(cltRepo), offercase.NewProductOfferUseCase(offerRepo, providerRepo, prdRepo, rangcase.NewRangUseCase(rankRepo, userRepo)), userRepo),
		rank:     rangcase.NewRangUseCase(rankRepo, userRepo),
		admin:    usecase.NewAdminUseCase(adminRepo),
	}
}

func initDB() *sql.DB {
	lib_db, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(fmt.Errorf("open db is failed: %v", err))
	}

	if err := lib_db.Ping(); err != nil {
		panic(fmt.Errorf("pinging db is failed: %v", err))
	}

	return lib_db
}

func (a *App) Run(port string) {
	router := mux.NewRouter()

	auth_delivery.RegisterEndpoints(router, a.authController, a.admin)
	provider_dlvr.RegisterEndpoints(router, a.provider)
	client_dlvr.RegisterEndpoints(router, a.clt)
	prddlvr.RegisterEndpoints(router, a.prd)
	offerdlvr.RegisterEndpoints(router, a.offer)
	orderdlvr.RegisterEndpoints(router, a.ord)
	rangdlvr.RegisterEndpoints(router, a.rank)

	router.Use(app.JWTAuthentication)
	router.Use(app.CheckAccessRight)
	router.Use(app.LogNewConnection)

	log.Fatalln(http.ListenAndServe(":"+port, router))
}
