package queue

import (
	"log/slog"

	"github.com/hibiken/asynq"
)

const (
	TypeOCR string = "ocr"
	TypeAI  string = "ai"
)

type (
	Queue struct {
		mux      *asynq.ServeMux
		auditlog *slog.Logger
	}
)

func (q *Queue) initialize() *Queue {
	q.mux.HandleFunc(TypeOCR, q.ocr)
	q.mux.HandleFunc(TypeAI, q.ai)

	return q
}

func CreateQueue(m *asynq.ServeMux, auditlog *slog.Logger) *Queue {
	return (&Queue{mux: m, auditlog: auditlog}).initialize()
}
