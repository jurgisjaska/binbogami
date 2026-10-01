package main

import (
	"fmt"
	"log"
	"log/slog"

	"github.com/hibiken/asynq"
	"github.com/jurgisjaska/binbogami/internal"
	"github.com/jurgisjaska/binbogami/internal/queue"
	audithandler "github.com/jurgisjaska/binbogami/internal/service/log"
)

func main() {
	log.Println("starting worker service")

	config, err := internal.CreateConfig()
	if err != nil {
		log.Fatalln("configuration load failure")
	}

	logger := slog.New(audithandler.CreateLoki(config.Loki))
	logger = logger.With("service", "_").WithGroup(audithandler.GroupSystem)
	slog.SetDefault(logger)
	slog.Info("starting worker service")
	defer slog.Warn("stopping worker service")

	auditlog := slog.New(audithandler.CreateLoki(config.Loki))
	auditlog = auditlog.With("service", "worker").WithGroup(audithandler.GroupAudit)

	server := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr: fmt.Sprintf("%s:%d", config.Redis.Connection.Hostname, config.Redis.Connection.Port),
		},
		asynq.Config{
			Concurrency: config.Redis.Concurrency,
		},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc(queue.TypeOCR, queue.HandleOCR)

	if err := server.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
