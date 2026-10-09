package main

import (
	"fmt"
	"log"
	"log/slog"

	"github.com/hibiken/asynq"
	"github.com/jurgisjaska/binbogami/internal"
	"github.com/jurgisjaska/binbogami/internal/queue"
	"github.com/jurgisjaska/binbogami/internal/service/ai"
	ls "github.com/jurgisjaska/binbogami/internal/service/log"
	"github.com/jurgisjaska/binbogami/internal/service/ocr"
)

func main() {
	log.Println("starting worker service")

	config, err := internal.CreateConfig()
	if err != nil {
		log.Fatalln("configuration load failure")
	}

	logger := slog.New(ls.CreateLoki(config.Loki))
	logger = logger.With("service", "_").WithGroup(ls.GroupSystem)
	slog.SetDefault(logger)
	slog.Info("starting worker service")
	defer slog.Warn("stopping worker service")

	auditlog := slog.New(ls.CreateLoki(config.Loki))
	auditlog = auditlog.With("service", "worker").WithGroup(ls.GroupAudit)

	database, err := internal.ConnectDatabase(config.Database)
	if err != nil {
		slog.Error("database connection failure", "error", err, "group", "system")
		log.Fatalln("database connection failure")
	}
	defer func() { _ = database.Close() }()

	server := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr: fmt.Sprintf("%s:%d", config.Redis.Connection.Hostname, config.Redis.Connection.Port),
		},
		asynq.Config{
			Concurrency: config.Redis.Concurrency,
			Logger:      ls.CreateAsynq(auditlog),
		},
	)

	// Tesseract OCR engine
	tesseract := ocr.CreateTesseract(auditlog)
	defer func() { _ = tesseract.Close() }()

	gemini, err := ai.CreateGemini(config.GenAI, auditlog)

	mux := asynq.NewServeMux()
	queue.CreateQueue(mux, database, auditlog, tesseract, gemini)

	if err := server.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
