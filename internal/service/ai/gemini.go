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
	"google.golang.org/genai/interactions/models/interactions"
)

const ModelGemini35FlashLite interactions.Model = "gemini-3.5-flash-lite"

type (
	Gemini struct {
		model    interactions.Model
		client   *genai.Client
		auditlog *slog.Logger
	}
)

func (g *Gemini) Extract(p string) (string, error) {
	g.auditlog.Info("extracting text from image", "engine", "gemini", "path", p)

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

	resp, err := g.client.Models.GenerateContent(context.Background(), string(g.model), contents, nil)
	if err != nil {
		return "", err
	}

	if resp.UsageMetadata != nil {
		g.auditlog.Info("text extracted from image",
			"path", p,
			"service", "engine",
			"prompt_tokens", resp.UsageMetadata.PromptTokenCount,
			"candidate_tokens", resp.UsageMetadata.CandidatesTokenCount,
			"total_tokens", resp.UsageMetadata.TotalTokenCount,
		)
	}

	return resp.Text(), nil
}

// CreateGemini initializes a Gemini AI instance with the specified model, GenAI configuration, and audit logger.
func CreateGemini(model interactions.Model, c *internal.GenAI, auditlog *slog.Logger) (*Gemini, error) {
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

	return &Gemini{model: model, client: client, auditlog: auditlog}, nil
}
