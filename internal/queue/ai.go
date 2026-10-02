package queue

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
)

type (
	AI struct {
		Resource string
	}
)

func CreateAI(resource string) (*asynq.Task, error) {
	payload, err := json.Marshal(AI{Resource: resource})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeAI, payload), nil
}

func (q *Queue) ai(c context.Context, t *asynq.Task) error {
	q.auditlog.Info("handling AI task")

	var p AI
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		q.auditlog.Error("queue ai error: json unmarshal failed", "error", err.Error())
		return err
	}

	q.auditlog.Info("AI task handled", "resource", p.Resource)

	return nil
}
