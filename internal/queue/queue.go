package queue

import (
	"log/slog"

	"github.com/hibiken/asynq"
	"github.com/jmoiron/sqlx"
	"github.com/jurgisjaska/binbogami/internal/service/ai"
	"github.com/jurgisjaska/binbogami/internal/service/ocr"
)

const (
	TypeOCR string = "ocr"
	TypeAI  string = "ai"
)

type (
	Queue struct {
		mux       *asynq.ServeMux
		database  *sqlx.DB
		auditlog  *slog.Logger
		tesseract *ocr.Tesseract
		gemini    *ai.Gemini
	}
)

func (q *Queue) initialize() *Queue {
	q.mux.HandleFunc(TypeOCR, q.ocr)
	q.mux.HandleFunc(TypeAI, q.ai)

	return q
}

// CreateQueue initializes a Queue instance and sets up task handlers for OCR and AI processing tasks.
func CreateQueue(m *asynq.ServeMux, db *sqlx.DB, auditlog *slog.Logger, t *ocr.Tesseract, ai *ai.Gemini) *Queue {
	return (&Queue{mux: m, database: db, auditlog: auditlog, tesseract: t, gemini: ai}).initialize()
}
