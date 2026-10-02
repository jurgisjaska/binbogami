package queue

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateOCR(t *testing.T) {
	tests := []struct {
		name     string
		resource string
	}{
		{
			name:     "stargate gate glyphs image",
			resource: "sgc://cheyenne-mountain/sub-level-28/gate-glyphs.png",
		},
		{
			name:     "example resource",
			resource: "this is an example",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := CreateOCR(tt.resource)
			require.NoError(t, err)
			require.NotNil(t, task)

			assert.Equal(t, TypeOCR, task.Type())

			var payload OCR
			err = json.Unmarshal(task.Payload(), &payload)
			require.NoError(t, err)
			assert.Equal(t, tt.resource, payload.Resource)
		})
	}
}

func TestHandleOCR(t *testing.T) {
	task, err := CreateOCR("sgc://dhd-symbols.png")
	require.NoError(t, err)

	q := &Queue{auditlog: slog.New(slog.NewTextHandler(io.Discard, nil))}
	err = q.ocr(context.Background(), task)
	assert.NoError(t, err)
}
