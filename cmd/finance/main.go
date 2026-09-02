package main

import (
	"fmt"
	"log"
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/jurgisjaska/binbogami/internal"
	"github.com/jurgisjaska/binbogami/internal/api"
	"github.com/jurgisjaska/binbogami/internal/api/token"
	"github.com/jurgisjaska/binbogami/internal/handlers/v1/finance"
	audithandler "github.com/jurgisjaska/binbogami/internal/service/log"
	echojwt "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

// Finance service provides finance related functionality.

func main() {
	log.Println("starting finance service")

	config, err := internal.CreateConfig()
	if err != nil {
		log.Fatalln("configuration load failure")
	}

	logger := slog.New(audithandler.CreateLoki(config.Loki))
	logger = logger.With("service", "finance").WithGroup(audithandler.GroupSystem)
	slog.SetDefault(logger)
	slog.Info("starting finance service")
	defer slog.Warn("stopping finance service")

	auditlog := slog.New(audithandler.CreateLoki(config.Loki))
	auditlog = auditlog.With("service", "finance").WithGroup(audithandler.GroupAudit)

	database, err := internal.ConnectDatabase(config.Database)
	if err != nil {
		slog.Error("database connection failure", "error", err, "group", "system")
		log.Fatalln("database connection failure")
	}
	defer func() { _ = database.Close() }()

	e := echo.New()
	e.Use(middleware.RequestLogger())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))
	e.HTTPErrorHandler = api.CustomHTTPErrorHandler
	e.Validator = &api.Validator{Validator: validator.New()}

	// @todo if this ever goes to production it needs to have proper values!
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"*"},
	}))

	g := e.Group("/v1")
	g.Use(echojwt.WithConfig(token.CreateJWTConfig(config.Secret)))

	finance.CreateFinance(g, database, auditlog)
	finance.CreateCategory(g, database, auditlog)
	finance.CreateLocation(g, database, auditlog)
	// book
	// entry

	if err := e.Start(fmt.Sprintf(":%d", config.Finance.Port)); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
