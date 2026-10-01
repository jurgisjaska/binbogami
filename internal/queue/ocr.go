package queue

import (
	"context"
	"encoding/json"
	"log"

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

func HandleOCR(ctx context.Context, t *asynq.Task) error {
	log.Println("Handling OCR task")

	return nil
}
