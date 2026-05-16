package log

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jurgisjaska/binbogami/internal"
)

const (
	GroupSystem string = "system"
	GroupAudit  string = "audit"
)

type (
	Loki struct {
		endpoint string
		username string
		password string
		service  string
		attrs    []slog.Attr
		group    string
		client   *http.Client
	}

	lokiStream struct {
		Stream map[string]string `json:"stream"`
		Values [][2]string       `json:"values"`
	}

	lokiPush struct {
		Streams []lokiStream `json:"streams"`
	}
)

func (h *Loki) Enabled(_ context.Context, _ slog.Level) bool {
	return true
}

func (h *Loki) Handle(_ context.Context, r slog.Record) error {
	labels := map[string]string{
		"level": r.Level.String(),
	}
	if h.service != "" {
		labels["service"] = h.service
	}
	if h.group != "" {
		labels["group"] = h.group
	}

	attrs := make(map[string]any, r.NumAttrs()+len(h.attrs))
	for _, a := range h.attrs {
		attrs[a.Key] = a.Value.Any()
	}
	r.Attrs(func(a slog.Attr) bool {
		attrs[h.prefixed(a.Key)] = a.Value.Any()
		return true
	})

	line := r.Message
	if len(attrs) > 0 {
		b, _ := json.Marshal(attrs)
		line = fmt.Sprintf("%s %s", r.Message, string(b))
	}

	t := r.Time
	if t.IsZero() {
		t = time.Now()
	}

	payload := lokiPush{
		Streams: []lokiStream{
			{
				Stream: labels,
				Values: [][2]string{{strconv.FormatInt(t.UnixNano(), 10), line}},
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, h.endpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if h.username != "" {
		req.SetBasicAuth(h.username, h.password)
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	return nil
}

func (h *Loki) WithAttrs(attrs []slog.Attr) slog.Handler {
	clone := *h
	var rest []slog.Attr
	for _, a := range attrs {
		if a.Key == "service" && clone.service == "" {
			clone.service = a.Value.String()
			continue
		}
		rest = append(rest, slog.Attr{Key: h.prefixed(a.Key), Value: a.Value})
	}
	clone.attrs = append(h.attrs[:len(h.attrs):len(h.attrs)], rest...)
	return &clone
}

func (h *Loki) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}

	clone := *h
	if h.group != "" {
		clone.group = strings.Join([]string{h.group, name}, ".")
	} else {
		clone.group = name
	}
	return &clone
}

func (h *Loki) prefixed(key string) string {
	if h.group == "" {
		return key
	}
	return h.group + "." + key
}

func CreateLoki(c *internal.Connection) *Loki {
	return &Loki{
		endpoint: fmt.Sprintf("http://%s:%d/loki/api/v1/push", c.Hostname, c.Port),
		username: c.Username,
		password: c.Password,
		client:   &http.Client{Timeout: 5 * time.Second},
	}
}
