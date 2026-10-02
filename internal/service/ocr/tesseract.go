package ocr

import "github.com/otiai10/gosseract/v2"

type Tesseract struct {
	client *gosseract.Client
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

func CreateTesseract() *Tesseract {
	return &Tesseract{client: gosseract.NewClient()}
}
