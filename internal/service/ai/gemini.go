package ai

import (
	"context"
	"log/slog"
	"mime"
	"net/http"
	"os"
	"path/filepath"

	"github.com/jurgisjaska/binbogami/internal"
	"google.golang.org/genai"
)

const model = "gemini-3.1-flash-lite"

type (
	Gemini struct {
		client   *genai.Client
		auditlog *slog.Logger
	}
)

func (g *Gemini) Extract(p string) (string, error) {
	g.auditlog.Info("extracting text from image", "path", p)

	b, err := os.ReadFile(p)
	if err != nil {
		return "", err
	}

	mimeType := http.DetectContentType(b)
	if mimeType == "application/octet-stream" {
		if ext := filepath.Ext(p); ext != "" {
			if m := mime.TypeByExtension(ext); m != "" {
				mimeType = m
			}
		}
	}

	contents := []*genai.Content{
		genai.NewContentFromParts([]*genai.Part{
			genai.NewPartFromBytes(b, mimeType),
			genai.NewPartFromText("Extract all text from this receipt image accurately. Return only the extracted text without commentary."),
		}, genai.RoleUser),
	}

	resp, err := g.client.Models.GenerateContent(context.Background(), model, contents, nil)
	if err != nil {
		return "", err
	}

	return resp.Text(), nil
}

func CreateGemini(c *internal.GenAI, auditlog *slog.Logger) (*Gemini, error) {
	cfg := &genai.ClientConfig{
		APIKey: c.APIKey,
	}

	if c.UseEnterprise {
		cfg.Backend = genai.BackendEnterprise
		cfg.Project = c.Project
		cfg.Location = c.Location
	}

	client, err := genai.NewClient(context.Background(), cfg)
	if err != nil {
		return nil, err
	}

	return &Gemini{client: client, auditlog: auditlog}, nil
}
