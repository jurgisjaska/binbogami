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
	if err := t.client.SetImage(p); err != nil {
		return "", err
	}

	return t.client.Text()
}

func (t *Tesseract) Close() error {
	return t.client.Close()
}

func CreateTesseract(auditlog *slog.Logger) *Tesseract {
	client := gosseract.NewClient()
	_ = client.SetLanguage("lit", "eng")
	_ = client.SetPageSegMode(gosseract.PSM_SINGLE_COLUMN)
	_ = client.SetVariable("preserve_interword_spaces", "1")

	return &Tesseract{client: client, auditlog: auditlog}
}
