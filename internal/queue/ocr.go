package queue

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
)

type (
	OCR struct {
		Resource string
	}
)

func CreateOCR(resource string) (*asynq.Task, error) {
	payload, err := json.Marshal(OCR{Resource: resource})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeOCR, payload), nil
}

func (q *Queue) ocr(c context.Context, t *asynq.Task) error {
	q.auditlog.Info("handling OCR task")

	var p OCR
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		q.auditlog.Error("queue ai error: json unmarshal failed", "error", err.Error())
		return err
	}

	q.auditlog.Info("AI task handled", "resource", p.Resource)

	return nil
}
