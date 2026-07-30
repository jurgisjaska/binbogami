package main

import (
	"fmt"
	"log"
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/jurgisjaska/binbogami/internal"
	"github.com/jurgisjaska/binbogami/internal/api"
	"github.com/jurgisjaska/binbogami/internal/handlers/auth"
	audithandler "github.com/jurgisjaska/binbogami/internal/service/log"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

// Auth service provides authentication and authorization functionality,
// including user login, registration, and JWT token management.

func main() {
	log.Println("starting auth service")

	config, err := internal.CreateConfig()
	if err != nil {
		log.Fatalln("configuration load failure")
	}

	logger := slog.New(audithandler.CreateLoki(config.Loki))
	logger = logger.With("service", "auth").WithGroup(audithandler.GroupSystem)
	slog.SetDefault(logger)
	slog.Info("starting auth service")
	defer slog.Warn("stopping auth service")

	auditlog := slog.New(audithandler.CreateLoki(config.Loki))
	auditlog = auditlog.With("service", "auth").WithGroup(audithandler.GroupAudit)

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

	dialer := internal.CreateDialer(config.Mail)
	auth.CreateAuth(e, database, config, dialer, auditlog)

	if err := e.Start(fmt.Sprintf(":%d", config.Auth.Port)); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
