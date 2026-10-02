package ocr

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fixtureDir is the path to the pre-rendered SGC test images relative to the
// repository root.  Tests are run from their package directory, so we walk up
// three levels (ocr → service → internal → repo root).
const fixtureDir = "../../../var/images"

// ---------------------------------------------------------------------------
// Unit tests — no Tesseract subprocess involved.
// ---------------------------------------------------------------------------

func TestCreateTesseract(t *testing.T) {
	svc := CreateTesseract()
	require.NotNil(t, svc)
	assert.NotNil(t, svc.client)

	// Cleanup must not error.
	assert.NoError(t, svc.Close())
}

func TestClose(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "close fresh client"},
		{name: "close second fresh client"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := CreateTesseract()
			require.NotNil(t, svc)
			assert.NoError(t, svc.Close())
		})
	}
}

// ---------------------------------------------------------------------------
// Integration tests — require a real Tesseract binary and the SGC fixture
// images produced by var/scripts/generate_ocr_fixtures.go.
// ---------------------------------------------------------------------------

func TestExtract_InvalidPath(t *testing.T) {
	tests := []struct {
		name string
		path string
	}{
		{
			name: "non-existent file",
			path: "/tmp/sg1_does_not_exist.png",
		},
		{
			name: "empty path",
			path: "",
		},
	}

	svc := CreateTesseract()
	require.NotNil(t, svc)
	t.Cleanup(func() { assert.NoError(t, svc.Close()) })

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			text, err := svc.Extract(tt.path)
			assert.Error(t, err)
			assert.Empty(t, text)
		})
	}
}

func TestExtract_MissionReport(t *testing.T) {
	path := fixtureDir + "/sg1_mission_report.png"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Skipf("fixture not found (%s) — run: go run var/scripts/generate_ocr_fixtures.go", path)
	}

	svc := CreateTesseract()
	require.NotNil(t, svc)
	t.Cleanup(func() { assert.NoError(t, svc.Close()) })

	text, err := svc.Extract(path)
	require.NoError(t, err)
	assert.NotEmpty(t, text)

	// Keywords are chosen to survive the predictable substitutions that
	// Tesseract 5 makes when reading the 5×7 pixel-font fixtures
	// (e.g. digit confusions 0→6, letter confusions N→W).
	tests := []struct {
		name    string
		keyword string
	}{
		{name: "document title", keyword: "MISSION REPORT"},
		{name: "classification label", keyword: "TOP SECRET"},
		{name: "document number prefix", keyword: "SGC-"},
		{name: "mission designation", keyword: "SG-1"},
		{name: "team leader title", keyword: "Colonel"},
		{name: "archaeologist name", keyword: "Jackson"},
		{name: "astrophysicist name", keyword: "Carter"},
		{name: "jaffa advisor name", keyword: "Teal"},
		{name: "mission status", keyword: "SUCCESS"},
		{name: "approver name", keyword: "Hammond"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.True(
				t,
				strings.Contains(text, tt.keyword),
				"expected extracted text to contain %q\ngot:\n%s", tt.keyword, text,
			)
		})
	}
}

func TestExtract_ExpenseReport(t *testing.T) {
	path := fixtureDir + "/sg1_expense_report.png"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Skipf("fixture not found (%s) — run: go run var/scripts/generate_ocr_fixtures.go", path)
	}

	svc := CreateTesseract()
	require.NotNil(t, svc)
	t.Cleanup(func() { assert.NoError(t, svc.Close()) })

	text, err := svc.Extract(path)
	require.NoError(t, err)
	assert.NotEmpty(t, text)

	tests := []struct {
		name    string
		keyword string
	}{
		{name: "document title", keyword: "EXPENSE REPORT"},
		{name: "document number prefix", keyword: "SGC-FIN"},
		{name: "submitter name", keyword: "Jackson"},
		{name: "first line item amount", keyword: "47.99"},
		{name: "net amount claimed", keyword: "378.24"},
		{name: "approver name", keyword: "Hammond"},
		{name: "finance reference prefix", keyword: "SGC-ACCTS"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.True(
				t,
				strings.Contains(text, tt.keyword),
				"expected extracted text to contain %q\ngot:\n%s", tt.keyword, text,
			)
		})
	}
}

// TestExtract_ReuseClient verifies that the same Tesseract client can process
// multiple images sequentially without error.
func TestExtract_ReuseClient(t *testing.T) {
	fixtures := []struct {
		name string
		path string
	}{
		{name: "mission report", path: fixtureDir + "/sg1_mission_report.png"},
		{name: "expense report", path: fixtureDir + "/sg1_expense_report.png"},
	}

	// Skip if neither fixture exists.
	anyFound := false
	for _, f := range fixtures {
		if _, err := os.Stat(f.path); err == nil {
			anyFound = true
			break
		}
	}
	if !anyFound {
		t.Skip("no fixtures found — run: go run var/scripts/generate_ocr_fixtures.go")
	}

	svc := CreateTesseract()
	require.NotNil(t, svc)
	t.Cleanup(func() { assert.NoError(t, svc.Close()) })

	for _, f := range fixtures {
		if _, err := os.Stat(f.path); os.IsNotExist(err) {
			t.Logf("skipping %s (file not found)", f.name)
			continue
		}

		t.Run(f.name, func(t *testing.T) {
			text, err := svc.Extract(f.path)
			assert.NoError(t, err)
			assert.NotEmpty(t, text)
		})
	}
}
