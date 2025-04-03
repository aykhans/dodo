package config

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/aykhans/dodo/types"
	"github.com/stretchr/testify/assert"
)

func TestReadFile(t *testing.T) {
	// Create temporary files for testing
	tempDir := t.TempDir()

	// Create a valid JSON config file
	validJSONFile := filepath.Join(tempDir, "valid.json")
	jsonContent := `{
		"method": "POST",
		"url": "https://example.com",
		"timeout": "5s",
		"dodos": 5,
		"requests": 100,
		"duration": "1m",
		"yes": true,
		"headers": [{"Content-Type": "application/json"}]
	}`
	err := os.WriteFile(validJSONFile, []byte(jsonContent), 0644)
	assert.NoError(t, err)

	// Create a valid YAML config file
	validYAMLFile := filepath.Join(tempDir, "valid.yaml")
	yamlContent := `
method: POST
url: https://example.com
timeout: 5s
dodos: 5
requests: 100
duration: 1m
yes: true
headers:
  - Content-Type: application/json
`
	err = os.WriteFile(validYAMLFile, []byte(yamlContent), 0644)
	assert.NoError(t, err)

	// Create an invalid JSON config file
	invalidJSONFile := filepath.Join(tempDir, "invalid.json")
	invalidJSONContent := `{
		"method": "POST",
		"url": "https://example.com",
		syntax error
	}`
	err = os.WriteFile(invalidJSONFile, []byte(invalidJSONContent), 0644)
	assert.NoError(t, err)

	// Create a file with unsupported extension
	unsupportedFile := filepath.Join(tempDir, "config.txt")
	err = os.WriteFile(unsupportedFile, []byte("some content"), 0644)
	assert.NoError(t, err)

	// Setup HTTP test server for remote config
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/valid.json":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, jsonContent)
		case "/invalid.json":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, invalidJSONContent)
		case "/valid.yaml":
			w.Header().Set("Content-Type", "application/yaml")
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, yamlContent)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	tests := []struct {
		name      string
		filePath  types.ConfigFile
		expectErr bool
		validate  func(t *testing.T, config *Config)
	}{
		{
			name:      "valid local JSON file",
			filePath:  types.ConfigFile(validJSONFile),
			expectErr: false,
			validate: func(t *testing.T, config *Config) {
				assert.Equal(t, "POST", *config.Method)
				assert.Equal(t, "https://example.com", config.URL.String())
				assert.Equal(t, int64(5000000000), config.Timeout.Nanoseconds())
				assert.Equal(t, uint(5), *config.DodosCount)
				assert.Equal(t, uint(100), *config.RequestCount)
				assert.Equal(t, int64(60000000000), config.Duration.Nanoseconds())
				assert.True(t, *config.Yes)
				assert.Equal(t, 1, len(config.Headers))
				assert.Equal(t, "Content-Type", config.Headers[0].Key)
				assert.Equal(t, "application/json", config.Headers[0].Value[0])
			},
		},
		{
			name:      "valid local YAML file",
			filePath:  types.ConfigFile(validYAMLFile),
			expectErr: false,
			validate: func(t *testing.T, config *Config) {
				assert.Equal(t, "POST", *config.Method)
				assert.Equal(t, "https://example.com", config.URL.String())
				assert.Equal(t, int64(5000000000), config.Timeout.Nanoseconds())
				assert.Equal(t, uint(5), *config.DodosCount)
				assert.Equal(t, uint(100), *config.RequestCount)
				assert.Equal(t, int64(60000000000), config.Duration.Nanoseconds())
				assert.True(t, *config.Yes)
				assert.Equal(t, 1, len(config.Headers))
				assert.Equal(t, "Content-Type", config.Headers[0].Key)
				assert.Equal(t, "application/json", config.Headers[0].Value[0])
			},
		},
		{
			name:      "invalid local JSON file",
			filePath:  types.ConfigFile(invalidJSONFile),
			expectErr: true,
			validate:  func(t *testing.T, config *Config) {},
		},
		{
			name:      "unsupported file extension",
			filePath:  types.ConfigFile(unsupportedFile),
			expectErr: true,
			validate:  func(t *testing.T, config *Config) {},
		},
		{
			name:      "non-existent file",
			filePath:  types.ConfigFile(filepath.Join(tempDir, "nonexistent.json")),
			expectErr: true,
			validate:  func(t *testing.T, config *Config) {},
		},
		{
			name:      "valid remote JSON file",
			filePath:  types.ConfigFile(server.URL + "/valid.json"),
			expectErr: false,
			validate: func(t *testing.T, config *Config) {
				assert.Equal(t, "POST", *config.Method)
				assert.Equal(t, "https://example.com", config.URL.String())
				assert.Equal(t, int64(5000000000), config.Timeout.Nanoseconds())
				assert.Equal(t, uint(5), *config.DodosCount)
				assert.Equal(t, uint(100), *config.RequestCount)
				assert.Equal(t, int64(60000000000), config.Duration.Nanoseconds())
				assert.True(t, *config.Yes)
				assert.Equal(t, 1, len(config.Headers))
				assert.Equal(t, "Content-Type", config.Headers[0].Key)
				assert.Equal(t, "application/json", config.Headers[0].Value[0])
			},
		},
		{
			name:      "valid remote YAML file",
			filePath:  types.ConfigFile(server.URL + "/valid.yaml"),
			expectErr: false,
			validate: func(t *testing.T, config *Config) {
				assert.Equal(t, "POST", *config.Method)
				assert.Equal(t, "https://example.com", config.URL.String())
				assert.Equal(t, int64(5000000000), config.Timeout.Nanoseconds())
				assert.Equal(t, uint(5), *config.DodosCount)
				assert.Equal(t, uint(100), *config.RequestCount)
				assert.Equal(t, int64(60000000000), config.Duration.Nanoseconds())
				assert.True(t, *config.Yes)
				assert.Equal(t, 1, len(config.Headers))
				assert.Equal(t, "Content-Type", config.Headers[0].Key)
				assert.Equal(t, "application/json", config.Headers[0].Value[0])
			},
		},
		{
			name:      "invalid remote JSON file",
			filePath:  types.ConfigFile(server.URL + "/invalid.json"),
			expectErr: true,
			validate:  func(t *testing.T, config *Config) {},
		},
		{
			name:      "non-existent remote file",
			filePath:  types.ConfigFile(server.URL + "/nonexistent.json"),
			expectErr: true,
			validate:  func(t *testing.T, config *Config) {},
		},
		{
			name:      "invalid URL",
			filePath:  types.ConfigFile("http://nonexistent.example.com/config.json"),
			expectErr: true,
			validate:  func(t *testing.T, config *Config) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := NewConfig()
			err := config.ReadFile(tt.filePath)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				tt.validate(t, config)
			}
		})
	}
}

func TestParseJSONConfig(t *testing.T) {
	tests := []struct {
		name      string
		jsonData  string
		expectErr bool
		validate  func(t *testing.T, config *Config)
	}{
		{
			name: "valid JSON config",
			jsonData: `{
				"method": "POST",
				"url": "https://example.com",
				"timeout": "5s",
				"dodos": 5,
				"requests": 100
			}`,
			expectErr: false,
			validate: func(t *testing.T, config *Config) {
				assert.Equal(t, "POST", *config.Method)
				assert.Equal(t, "https://example.com", config.URL.String())
				assert.Equal(t, int64(5000000000), config.Timeout.Nanoseconds())
				assert.Equal(t, uint(5), *config.DodosCount)
				assert.Equal(t, uint(100), *config.RequestCount)
			},
		},
		{
			name: "invalid JSON syntax",
			jsonData: `{
				"method": "POST",
				"url": "https://example.com",
				syntax error
			}`,
			expectErr: true,
			validate:  func(t *testing.T, config *Config) {},
		},
		{
			name: "invalid type for field",
			jsonData: `{
				"method": "POST",
				"url": "https://example.com",
				"dodos": "not-a-number"
			}`,
			expectErr: true,
			validate:  func(t *testing.T, config *Config) {},
		},
		{
			name:      "empty JSON object",
			jsonData:  `{}`,
			expectErr: false,
			validate: func(t *testing.T, config *Config) {
				assert.Nil(t, config.Method)
				assert.Nil(t, config.URL)
				assert.Nil(t, config.Timeout)
				assert.Nil(t, config.DodosCount)
				assert.Nil(t, config.RequestCount)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := NewConfig()
			err := parseJSONConfig([]byte(tt.jsonData), config)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				tt.validate(t, config)
			}
		})
	}
}

func TestParseYAMLConfig(t *testing.T) {
	tests := []struct {
		name      string
		yamlData  string
		expectErr bool
		validate  func(t *testing.T, config *Config)
	}{
		{
			name: "valid YAML config",
			yamlData: `
method: POST
url: https://example.com
timeout: 5s
dodos: 5
requests: 100
`,
			expectErr: false,
			validate: func(t *testing.T, config *Config) {
				assert.Equal(t, "POST", *config.Method)
				assert.Equal(t, "https://example.com", config.URL.String())
				assert.Equal(t, int64(5000000000), config.Timeout.Nanoseconds())
				assert.Equal(t, uint(5), *config.DodosCount)
				assert.Equal(t, uint(100), *config.RequestCount)
			},
		},
		{
			name: "invalid YAML syntax",
			yamlData: `
method: POST
url: https://example.com
dodos: 5
    invalid indentation
`,
			expectErr: true,
			validate:  func(t *testing.T, config *Config) {},
		},
		{
			name:      "empty YAML",
			yamlData:  ``,
			expectErr: false,
			validate: func(t *testing.T, config *Config) {
				assert.Nil(t, config.Method)
				assert.Nil(t, config.URL)
				assert.Nil(t, config.Timeout)
				assert.Nil(t, config.DodosCount)
				assert.Nil(t, config.RequestCount)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := NewConfig()
			err := parseYAMLConfig([]byte(tt.yamlData), config)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				tt.validate(t, config)
			}
		})
	}
}
