package config

import (
	"flag"
	"io"
	"os"
	"testing"
	"time"

	"github.com/aykhans/dodo/types"
	"github.com/stretchr/testify/assert"
)

func TestReadCLI(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectFile     types.ConfigFile
		expectError    bool
		expectedConfig *Config
	}{
		{
			name:        "simple url and duration",
			args:        []string{"-u", "https://example.com", "-o", "1m"},
			expectFile:  "",
			expectError: false,
			expectedConfig: &Config{
				URL:      &types.RequestURL{},
				Duration: &types.Duration{Duration: time.Minute},
			},
		},
		{
			name:           "config file only",
			args:           []string{"-f", "/path/to/config.json"},
			expectFile:     "/path/to/config.json",
			expectError:    false,
			expectedConfig: &Config{},
		},
		{
			name:        "all flags",
			args:        []string{"-f", "/path/to/config.json", "-u", "https://example.com", "-m", "POST", "-d", "10", "-r", "1000", "-o", "3m", "-t", "3s", "-b", "body1", "-H", "header1:value1", "-p", "param1=value1", "-c", "cookie1=value1", "-x", "http://proxy.example.com:8080", "-y"},
			expectFile:  "/path/to/config.json",
			expectError: false,
			expectedConfig: &Config{
				Method:       stringPtr("POST"),
				URL:          &types.RequestURL{},
				DodosCount:   uintPtr(10),
				RequestCount: uintPtr(1000),
				Duration:     &types.Duration{Duration: 3 * time.Minute},
				Timeout:      &types.Timeout{Duration: 3 * time.Second},
				Yes:          boolPtr(true),
			},
		},
		{
			name:           "unexpected arguments",
			args:           []string{"-u", "https://example.com", "extraArg"},
			expectFile:     "",
			expectError:    true,
			expectedConfig: &Config{},
		},
	}

	// Save original command-line arguments
	origArgs := os.Args
	origFlagCommandLine := flag.CommandLine

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flag.CommandLine to its original state
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			// Override os.Args for the test
			os.Args = append([]string{"dodo"}, tt.args...)

			// Initialize a new config
			config := NewConfig()

			// Mock URL to avoid actual URL parsing issues in tests
			if tt.expectedConfig.URL != nil {
				urlObj := types.RequestURL{}
				urlObj.Set("https://example.com")
				tt.expectedConfig.URL = &urlObj
			}

			// Call the function being tested
			file, err := config.ReadCLI()

			// Reset os.Args after test
			os.Args = origArgs

			// Assert expected results
			assert.Equal(t, tt.expectFile, file)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Check expected config values if URL is set
				if tt.expectedConfig.URL != nil {
					assert.NotNil(t, config.URL)
					assert.Equal(t, "https://example.com", config.URL.String())
				}

				// Check duration if expected
				if tt.expectedConfig.Duration != nil {
					assert.NotNil(t, config.Duration)
					assert.Equal(t, tt.expectedConfig.Duration.Duration, config.Duration.Duration)
				}

				// Check other values as needed
				if tt.expectedConfig.Method != nil {
					assert.Equal(t, *tt.expectedConfig.Method, *config.Method)
				}
				if tt.expectedConfig.DodosCount != nil {
					assert.Equal(t, *tt.expectedConfig.DodosCount, *config.DodosCount)
				}
				if tt.expectedConfig.RequestCount != nil {
					assert.Equal(t, *tt.expectedConfig.RequestCount, *config.RequestCount)
				}
				if tt.expectedConfig.Timeout != nil {
					assert.Equal(t, tt.expectedConfig.Timeout.Duration, config.Timeout.Duration)
				}
				if tt.expectedConfig.Yes != nil {
					assert.Equal(t, *tt.expectedConfig.Yes, *config.Yes)
				}
			}
		})
	}

	// Restore original flag.CommandLine
	flag.CommandLine = origFlagCommandLine
}

// Skip the prompt tests as they require interactive input/output handling
// which is difficult to test reliably in unit tests
func TestCLIYesOrNoReaderBasic(t *testing.T) {
	// We're just going to verify the function exists and returns the default value
	// when called with "\n" as input (which should trigger the default path)
	result := func() bool {
		// Save original standard input
		origStdin := os.Stdin
		origStdout := os.Stdout

		// Create a pipe to mock standard input
		r, w, _ := os.Pipe()
		os.Stdin = r

		// Redirect stdout to null device
		devNull, _ := os.Open(os.DevNull)
		os.Stdout = devNull

		// Write newline to mock stdin (should trigger default behavior)
		io.WriteString(w, "\n")
		w.Close()

		// Call the function being tested with default=true
		result := CLIYesOrNoReader("Test message", true)

		// Restore original stdin and stdout
		os.Stdin = origStdin
		os.Stdout = origStdout

		return result
	}()

	// Default value should be returned
	assert.True(t, result)
}

// Helper types and functions for testing
func stringPtr(s string) *string {
	return &s
}

func uintPtr(u uint) *uint {
	return &u
}

func boolPtr(b bool) *bool {
	return &b
}
