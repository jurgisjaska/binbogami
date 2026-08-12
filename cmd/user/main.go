package main

import (
	"fmt"
	"log"
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/jurgisjaska/binbogami/internal"
	"github.com/jurgisjaska/binbogami/internal/api"
	"github.com/jurgisjaska/binbogami/internal/api/token"
	"github.com/jurgisjaska/binbogami/internal/handlers/v1/user"
	audithandler "github.com/jurgisjaska/binbogami/internal/service/log"
	echojwt "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

//

func main() {
	log.Println("starting user service")

	config, err := internal.CreateConfig()
	if err != nil {
		log.Fatalln("configuration load failure")
	}

	logger := slog.New(audithandler.CreateLoki(config.Loki))
	logger = logger.With("service", "user").WithGroup(audithandler.GroupSystem)
	slog.SetDefault(logger)
	slog.Info("starting user service")
	defer slog.Warn("stopping user service")

	auditlog := slog.New(audithandler.CreateLoki(config.Loki))
	auditlog = auditlog.With("service", "user").WithGroup(audithandler.GroupAudit)

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

	user.CreateUser(g, database, auditlog)

	if err := e.Start(fmt.Sprintf(":%d", config.User.Port)); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
