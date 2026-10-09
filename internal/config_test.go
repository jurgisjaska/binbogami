package internal

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPort(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{
			name:     "valid port",
			input:    "8080",
			expected: 8080,
		},
		{
			name:     "empty string",
			input:    "",
			expected: 0,
		},
		{
			name:     "invalid string",
			input:    "cheyenne-mountain",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, port(tt.input))
		})
	}
}

func TestBoolean(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "true lowercase",
			input:    "true",
			expected: true,
		},
		{
			name:     "1 numeric",
			input:    "1",
			expected: true,
		},
		{
			name:     "false lowercase",
			input:    "false",
			expected: false,
		},
		{
			name:     "0 numeric",
			input:    "0",
			expected: false,
		},
		{
			name:     "empty string",
			input:    "",
			expected: false,
		},
		{
			name:     "invalid string",
			input:    "stargate",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, boolean(tt.input))
		})
	}
}

func TestURI(t *testing.T) {
	t.Setenv("SGC_SCHEME", "https")
	t.Setenv("SGC_HOSTNAME", "sgc.cheyenne.af.mil")
	t.Setenv("SGC_PORT", "8443")

	res := uri("SGC")
	assert.Equal(t, "https", res.Scheme)
	assert.Equal(t, "sgc.cheyenne.af.mil", res.Hostname)
	assert.Equal(t, 8443, res.Port)
}

func TestCreateConfig(t *testing.T) {
	tempDir := t.TempDir()
	envContent := `APP_SECRET=sgc-classified-secret
APP_SCHEME=http
APP_HOSTNAME=localhost
APP_PORT=8101
AUTH_SERVICE_SECRET=sg1-auth-secret
AUTH_SERVICE_SCHEME=http
AUTH_SERVICE_HOSTNAME=auth
AUTH_SERVICE_PORT=8101
USER_SERVICE_SCHEME=http
USER_SERVICE_HOSTNAME=user
USER_SERVICE_PORT=8103
FINANCE_SERVICE_SCHEME=http
FINANCE_SERVICE_HOSTNAME=finance
FINANCE_SERVICE_PORT=8104
WEB_APPLICATION_SCHEME=http
WEB_APPLICATION_HOSTNAME=web
WEB_APPLICATION_PORT=8191
DATABASE_NAME=stargate
DATABASE_HOSTNAME=mariadb
DATABASE_PORT=3306
DATABASE_USERNAME=sg1
DATABASE_PASSWORD=dialthegate
MAIL_SENDER=sgc@cheyenne.af.mil
MAIL_HOSTNAME=mailcatcher
MAIL_PORT=1025
MAIL_USERNAME=sgc
MAIL_PASSWORD=dialthegate
LOKI_HOSTNAME=loki
LOKI_PORT=3100
LOKI_USERNAME=sgc
LOKI_PASSWORD=dialthegate
REDIS_HOSTNAME=redis
REDIS_PORT=6379
REDIS_USERNAME=sgc
REDIS_PASSWORD=dialthegate
GOOGLE_API_KEY=SG1-ANCIENT-GENAI-KEY
GOOGLE_GENAI_USE_ENTERPRISE=true
GOOGLE_CLOUD_PROJECT=stargate-command-project
GOOGLE_CLOUD_LOCATION=us-central1
`
	err := os.WriteFile(filepath.Join(tempDir, ".env"), []byte(envContent), 0600)
	require.NoError(t, err)

	t.Chdir(tempDir)

	config, err := CreateConfig()
	require.NoError(t, err)
	require.NotNil(t, config)

	assert.Equal(t, "sgc-classified-secret", config.Secret)
	require.NotNil(t, config.GenAI)
	assert.Equal(t, "SG1-ANCIENT-GENAI-KEY", config.GenAI.APIKey)
	assert.True(t, config.GenAI.UseEnterprise)
	assert.Equal(t, "stargate-command-project", config.GenAI.Project)
	assert.Equal(t, "us-central1", config.GenAI.Location)
}
