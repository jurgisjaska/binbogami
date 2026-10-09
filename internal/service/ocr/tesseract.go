package ocr

import (
	"log/slog"

	"github.com/otiai10/gosseract/v2"
)

type Tesseract struct {
	client   *gosseract.Client
	auditlog *slog.Logger
}

func (t *Tesseract) Extract(p string) (string, error) {
	t.auditlog.Info("extracting text from image", "engine", "tesseract", "path", p)

	if err := t.client.SetImage(p); err != nil {
		return "", err
	}

	t.auditlog.Info("text extracted from image", "engine", "tesseract", "path", p)

	return t.client.Text()
}

func (t *Tesseract) Close() error {
	return t.client.Close()
}

// CreateTesseract initializes a Tesseract OCR engine with specified logging and default configurations.
func CreateTesseract(auditlog *slog.Logger) *Tesseract {
	client := gosseract.NewClient()
	_ = client.SetLanguage("lit", "eng")
	_ = client.SetPageSegMode(gosseract.PSM_SINGLE_COLUMN)
	_ = client.SetVariable("preserve_interword_spaces", "1")

	return &Tesseract{client: client, auditlog: auditlog}
}
