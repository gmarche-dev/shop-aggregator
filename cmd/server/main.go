package main

import (
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"shop-aggregator/internal/config"
	"shop-aggregator/internal/db/postgresql"
	"shop-aggregator/internal/handler"
	"shop-aggregator/internal/router"
	"shop-aggregator/internal/usecase"
	"shop-aggregator/tools/migrations"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = zerolog.New(os.Stderr).With().Timestamp().Logger()

	cfg, err := config.LoadConfig("./config/config.yaml")
	if err != nil {
		log.Fatal().Caller().Err(err).Msg("Loading config failed")
	}

	db, err := postgresql.NewDB(context.Background(), &cfg.Database)
	if err != nil {
		log.Fatal().Caller().Err(err).Msg("NewDB error")
	}

	// run migrations
	if err = migrations.Run(context.Background(), db, "../../migrations/deploy"); err != nil {
		log.Fatal().Caller().Err(err).Msg("Migrations error")
	}
	e := gin.Default()

	// Middleware
	e.Use(gin.Logger())
	e.Use(gin.Recovery())
	e.Use(cors.Default())

	sqlAuth := postgresql.NewAuth(db)
	sqlUser := postgresql.NewUsers(db)
	sqlBrand := postgresql.NewBrand(db)
	sqlCompany := postgresql.NewCompany(db)
	sqlBill := postgresql.NewBill(db)
	sqlStore := postgresql.NewStore(db)
	sqlProduct := postgresql.NewProduct(db)
	sqlUserProduct := postgresql.NewUserProduct(db)

	useCaseAuth := usecase.NewAuth(sqlAuth, sqlUser)
	useCaseUser := usecase.NewUsers(sqlUser)
	useCaseBrand := usecase.NewBrand(sqlBrand)
	useCaseCompany := usecase.NewCompany(sqlCompany)
	useCaseBill := usecase.NewBill(sqlBill, sqlStore, sqlCompany, sqlUserProduct)
	useCaseStore := usecase.NewStore(sqlStore, sqlCompany)
	useCaseProduct := usecase.NewProduct(sqlProduct, sqlBrand)
	useCaseUserProduct := usecase.NewUserProduct(sqlUserProduct)

	handlerAuth := handler.NewAuth(useCaseAuth)
	handlerUser := handler.NewUser(useCaseUser)
	handlerBrand := handler.NewBrand(useCaseBrand)
	handlerCompany := handler.NewCompany(useCaseCompany)
	handlerBill := handler.NewBill(useCaseBill)
	handlerStore := handler.NewStore(useCaseStore)
	handlerProduct := handler.NewProduct(useCaseProduct)
	handlerUserProduct := handler.NewUserProduct(useCaseUserProduct)
	handlerInitialisation := handler.NewInitialisation()

	r := router.NewRouter(e, sqlAuth, handlerAuth, handlerUser, handlerBrand, handlerCompany, handlerBill, handlerStore, handlerProduct, handlerUserProduct, handlerInitialisation)
	log.Info().Caller().Msgf("Starting server on port %d", cfg.Server.Port)
	if err = r.Run(fmt.Sprintf(":%d", cfg.Server.Port)); err != nil {
		log.Fatal().Caller().Err(err).Msg("Loading router failed")
	}
}
