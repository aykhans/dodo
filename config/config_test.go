package config

import (
	"net/url"
	"os"
	"slices"
	"testing"
	"time"

	"github.com/aykhans/dodo/types"
	"github.com/aykhans/dodo/utils"
	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	config := NewConfig()
	assert.NotNil(t, config)
	assert.IsType(t, &Config{}, config)
}

func TestNewRequestConfig(t *testing.T) {
	// Create a sample Config object
	urlObj := types.RequestURL{}
	urlObj.Set("https://example.com")

	conf := &Config{
		Method:       utils.ToPtr("GET"),
		URL:          &urlObj,
		Timeout:      &types.Timeout{Duration: 5 * time.Second},
		DodosCount:   utils.ToPtr(uint(10)),
		RequestCount: utils.ToPtr(uint(100)),
		Duration:     &types.Duration{Duration: 1 * time.Minute},
		Yes:          utils.ToPtr(true),
		Params:       types.Params{{Key: "key1", Value: []string{"value1"}}},
		Headers:      types.Headers{{Key: "User-Agent", Value: []string{"TestAgent"}}},
		Cookies:      types.Cookies{{Key: "session", Value: []string{"123"}}},
		Body:         types.Body{"test body"},
		Proxies:      types.Proxies{url.URL{Scheme: "http", Host: "proxy.example.com:8080"}},
	}

	// Call the function being tested
	rc := NewRequestConfig(conf)

	// Assert the fields are correctly mapped
	assert.Equal(t, "GET", rc.Method)
	assert.Equal(t, "https://example.com", rc.URL.String())
	assert.Equal(t, 5*time.Second, rc.Timeout)
	assert.Equal(t, uint(10), rc.DodosCount)
	assert.Equal(t, uint(100), rc.RequestCount)
	assert.Equal(t, 1*time.Minute, rc.Duration)
	assert.True(t, rc.Yes)
	assert.Equal(t, types.Params{{Key: "key1", Value: []string{"value1"}}}, rc.Params)
	assert.Equal(t, types.Headers{{Key: "User-Agent", Value: []string{"TestAgent"}}}, rc.Headers)
	assert.Equal(t, types.Cookies{{Key: "session", Value: []string{"123"}}}, rc.Cookies)
	assert.Equal(t, types.Body{"test body"}, rc.Body)
	assert.Equal(t, types.Proxies{url.URL{Scheme: "http", Host: "proxy.example.com:8080"}}, rc.Proxies)
}

func TestGetValidDodosCountForRequests(t *testing.T) {
	tests := []struct {
		name         string
		dodosCount   uint
		requestCount uint
		expected     uint
	}{
		{
			name:         "no request count limit",
			dodosCount:   10,
			requestCount: 0,
			expected:     10,
		},
		{
			name:         "dodos count less than request count",
			dodosCount:   5,
			requestCount: 100,
			expected:     5,
		},
		{
			name:         "dodos count greater than request count",
			dodosCount:   100,
			requestCount: 10,
			expected:     10,
		},
		{
			name:         "dodos count equal to request count",
			dodosCount:   50,
			requestCount: 50,
			expected:     50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc := &RequestConfig{
				DodosCount:   tt.dodosCount,
				RequestCount: tt.requestCount,
			}
			result := rc.GetValidDodosCountForRequests()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetMaxConns(t *testing.T) {
	tests := []struct {
		name         string
		dodosCount   uint
		requestCount uint
		minConns     uint
		expected     uint
	}{
		{
			name:         "min connections higher than valid dodos count",
			dodosCount:   10,
			requestCount: 0,
			minConns:     20,
			expected:     30, // 20 * 150%
		},
		{
			name:         "min connections lower than valid dodos count",
			dodosCount:   30,
			requestCount: 0,
			minConns:     10,
			expected:     45, // 30 * 150%
		},
		{
			name:         "request count limits dodos count",
			dodosCount:   100,
			requestCount: 20,
			minConns:     5,
			expected:     30, // 20 * 150%
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc := &RequestConfig{
				DodosCount:   tt.dodosCount,
				RequestCount: tt.requestCount,
			}
			result := rc.GetMaxConns(tt.minConns)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Skip the Print test as it's mainly a formatting function
// that uses external table rendering library
func TestRequestConfigPrint(t *testing.T) {
	// Create a sample RequestConfig
	rc := &RequestConfig{
		Method:       "GET",
		URL:          url.URL{Scheme: "https", Host: "example.com"},
		Timeout:      5 * time.Second,
		DodosCount:   10,
		RequestCount: 100,
		Duration:     1 * time.Minute,
		Params:       types.Params{{Key: "param1", Value: []string{"value1"}}},
		Headers:      types.Headers{{Key: "User-Agent", Value: []string{"TestAgent"}}},
		Cookies:      types.Cookies{{Key: "session", Value: []string{"123"}}},
		Body:         types.Body{"test body"},
		Proxies:      types.Proxies{url.URL{Scheme: "http", Host: "proxy.example.com:8080"}},
	}

	// We'll just call the function to ensure it doesn't panic
	// Redirect output to /dev/null
	origStdout := os.Stdout
	devNull, _ := os.Open(os.DevNull)
	os.Stdout = devNull

	// Call the function
	rc.Print()

	// Restore stdout
	os.Stdout = origStdout

	// No assertions needed, we're just checking that it doesn't panic
}

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name           string
		config         *Config
		expectedErrors int
		expectURLQuery bool
	}{
		{
			name: "valid config",
			config: &Config{
				Method:       utils.ToPtr("GET"),
				URL:          func() *types.RequestURL { u := &types.RequestURL{}; u.Set("https://example.com"); return u }(),
				Timeout:      &types.Timeout{Duration: 5 * time.Second},
				DodosCount:   utils.ToPtr(uint(10)),
				RequestCount: utils.ToPtr(uint(100)),
				Duration:     &types.Duration{Duration: 1 * time.Minute},
			},
			expectedErrors: 0,
			expectURLQuery: false,
		},
		{
			name: "missing URL",
			config: &Config{
				Method:       utils.ToPtr("GET"),
				URL:          nil,
				Timeout:      &types.Timeout{Duration: 5 * time.Second},
				DodosCount:   utils.ToPtr(uint(10)),
				RequestCount: utils.ToPtr(uint(100)),
				Duration:     &types.Duration{Duration: 1 * time.Minute},
			},
			expectedErrors: 1,
			expectURLQuery: false,
		},
		{
			name: "missing method",
			config: &Config{
				Method:       nil,
				URL:          func() *types.RequestURL { u := &types.RequestURL{}; u.Set("https://example.com"); return u }(),
				Timeout:      &types.Timeout{Duration: 5 * time.Second},
				DodosCount:   utils.ToPtr(uint(10)),
				RequestCount: utils.ToPtr(uint(100)),
				Duration:     &types.Duration{Duration: 1 * time.Minute},
			},
			expectedErrors: 1,
			expectURLQuery: false,
		},
		{
			name: "invalid URL scheme",
			config: &Config{
				Method:       utils.ToPtr("GET"),
				URL:          func() *types.RequestURL { u := &types.RequestURL{}; u.Scheme = "ftp"; u.Host = "example.com"; return u }(),
				Timeout:      &types.Timeout{Duration: 5 * time.Second},
				DodosCount:   utils.ToPtr(uint(10)),
				RequestCount: utils.ToPtr(uint(100)),
				Duration:     &types.Duration{Duration: 1 * time.Minute},
			},
			expectedErrors: 1,
			expectURLQuery: false,
		},
		{
			name: "missing both duration and request count",
			config: &Config{
				Method:       utils.ToPtr("GET"),
				URL:          func() *types.RequestURL { u := &types.RequestURL{}; u.Set("https://example.com"); return u }(),
				Timeout:      &types.Timeout{Duration: 5 * time.Second},
				DodosCount:   utils.ToPtr(uint(10)),
				RequestCount: utils.ToPtr(uint(0)),
				Duration:     &types.Duration{Duration: 0},
			},
			expectedErrors: 1,
			expectURLQuery: false,
		},
		{
			name: "URL with query parameters",
			config: &Config{
				Method: utils.ToPtr("GET"),
				URL: func() *types.RequestURL {
					u := &types.RequestURL{}
					u.Set("https://example.com?param1=value1&param2=value2")
					return u
				}(),
				Timeout:      &types.Timeout{Duration: 5 * time.Second},
				DodosCount:   utils.ToPtr(uint(10)),
				RequestCount: utils.ToPtr(uint(100)),
				Duration:     &types.Duration{Duration: 1 * time.Minute},
			},
			expectedErrors: 0,
			expectURLQuery: true,
		},
		{
			name: "invalid proxy scheme",
			config: &Config{
				Method:       utils.ToPtr("GET"),
				URL:          func() *types.RequestURL { u := &types.RequestURL{}; u.Set("https://example.com"); return u }(),
				Timeout:      &types.Timeout{Duration: 5 * time.Second},
				DodosCount:   utils.ToPtr(uint(10)),
				RequestCount: utils.ToPtr(uint(100)),
				Duration:     &types.Duration{Duration: 1 * time.Minute},
				Proxies:      types.Proxies{url.URL{Scheme: "invalid", Host: "proxy.example.com:8080"}},
			},
			expectedErrors: 1,
			expectURLQuery: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := tt.config.Validate()

			// Check number of errors
			assert.Len(t, errors, tt.expectedErrors)

			// Check if URL query parameters are extracted properly
			if tt.expectURLQuery {
				assert.Empty(t, tt.config.URL.RawQuery)
				assert.NotEmpty(t, tt.config.Params)
				found := false
				for _, param := range tt.config.Params {
					if param.Key == "param1" && len(param.Value) > 0 && param.Value[0] == "value1" {
						found = true
						break
					}
				}
				assert.True(t, found, "Expected param1=value1 in Params but not found")
			}
		})
	}
}

func TestMergeConfig(t *testing.T) {
	baseURL := func() *types.RequestURL { u := &types.RequestURL{}; u.Set("https://example.com"); return u }()
	newURL := func() *types.RequestURL { u := &types.RequestURL{}; u.Set("https://new-example.com"); return u }()

	baseConfig := &Config{
		Method:       utils.ToPtr("GET"),
		URL:          baseURL,
		Timeout:      &types.Timeout{Duration: 5 * time.Second},
		DodosCount:   utils.ToPtr(uint(10)),
		RequestCount: utils.ToPtr(uint(100)),
		Duration:     &types.Duration{Duration: 1 * time.Minute},
		Yes:          utils.ToPtr(false),
		Params:       types.Params{{Key: "base-param", Value: []string{"base-value"}}},
		Headers:      types.Headers{{Key: "base-header", Value: []string{"base-value"}}},
		Cookies:      types.Cookies{{Key: "base-cookie", Value: []string{"base-value"}}},
		Body:         types.Body{"base-body"},
		Proxies:      types.Proxies{url.URL{Scheme: "http", Host: "base-proxy.example.com:8080"}},
	}

	tests := []struct {
		name       string
		newConfig  *Config
		assertions func(t *testing.T, result *Config)
	}{
		{
			name: "merge all fields",
			newConfig: &Config{
				Method:       utils.ToPtr("POST"),
				URL:          newURL,
				Timeout:      &types.Timeout{Duration: 10 * time.Second},
				DodosCount:   utils.ToPtr(uint(20)),
				RequestCount: utils.ToPtr(uint(200)),
				Duration:     &types.Duration{Duration: 2 * time.Minute},
				Yes:          utils.ToPtr(true),
				Params:       types.Params{{Key: "new-param", Value: []string{"new-value"}}},
				Headers:      types.Headers{{Key: "new-header", Value: []string{"new-value"}}},
				Cookies:      types.Cookies{{Key: "new-cookie", Value: []string{"new-value"}}},
				Body:         types.Body{"new-body"},
				Proxies:      types.Proxies{url.URL{Scheme: "http", Host: "new-proxy.example.com:8080"}},
			},
			assertions: func(t *testing.T, result *Config) {
				assert.Equal(t, "POST", *result.Method)
				assert.Equal(t, "https://new-example.com", result.URL.String())
				assert.Equal(t, 10*time.Second, result.Timeout.Duration)
				assert.Equal(t, uint(20), *result.DodosCount)
				assert.Equal(t, uint(200), *result.RequestCount)
				assert.Equal(t, 2*time.Minute, result.Duration.Duration)
				assert.True(t, *result.Yes)
				assert.Equal(t, types.Params{{Key: "new-param", Value: []string{"new-value"}}}, result.Params)
				assert.Equal(t, types.Headers{{Key: "new-header", Value: []string{"new-value"}}}, result.Headers)
				assert.Equal(t, types.Cookies{{Key: "new-cookie", Value: []string{"new-value"}}}, result.Cookies)
				assert.Equal(t, types.Body{"new-body"}, result.Body)
				assert.Equal(t, types.Proxies{url.URL{Scheme: "http", Host: "new-proxy.example.com:8080"}}, result.Proxies)
			},
		},
		{
			name: "merge only specified fields",
			newConfig: &Config{
				Method: utils.ToPtr("POST"),
				URL:    newURL,
				Yes:    utils.ToPtr(true),
			},
			assertions: func(t *testing.T, result *Config) {
				assert.Equal(t, "POST", *result.Method)
				assert.Equal(t, "https://new-example.com", result.URL.String())
				assert.Equal(t, 5*time.Second, result.Timeout.Duration)                                                      // unchanged
				assert.Equal(t, uint(10), *result.DodosCount)                                                                // unchanged
				assert.Equal(t, uint(100), *result.RequestCount)                                                             // unchanged
				assert.Equal(t, 1*time.Minute, result.Duration.Duration)                                                     // unchanged
				assert.True(t, *result.Yes)                                                                                  // changed
				assert.Equal(t, types.Params{{Key: "base-param", Value: []string{"base-value"}}}, result.Params)             // unchanged
				assert.Equal(t, types.Headers{{Key: "base-header", Value: []string{"base-value"}}}, result.Headers)          // unchanged
				assert.Equal(t, types.Cookies{{Key: "base-cookie", Value: []string{"base-value"}}}, result.Cookies)          // unchanged
				assert.Equal(t, types.Body{"base-body"}, result.Body)                                                        // unchanged
				assert.Equal(t, types.Proxies{url.URL{Scheme: "http", Host: "base-proxy.example.com:8080"}}, result.Proxies) // unchanged
			},
		},
		{
			name:      "merge empty config",
			newConfig: &Config{},
			assertions: func(t *testing.T, result *Config) {
				// All fields should remain unchanged
				assert.Equal(t, "GET", *result.Method)
				assert.Equal(t, "https://example.com", result.URL.String())
				assert.Equal(t, 5*time.Second, result.Timeout.Duration)
				assert.Equal(t, uint(10), *result.DodosCount)
				assert.Equal(t, uint(100), *result.RequestCount)
				assert.Equal(t, 1*time.Minute, result.Duration.Duration)
				assert.False(t, *result.Yes)
				assert.Equal(t, types.Params{{Key: "base-param", Value: []string{"base-value"}}}, result.Params)
				assert.Equal(t, types.Headers{{Key: "base-header", Value: []string{"base-value"}}}, result.Headers)
				assert.Equal(t, types.Cookies{{Key: "base-cookie", Value: []string{"base-value"}}}, result.Cookies)
				assert.Equal(t, types.Body{"base-body"}, result.Body)
				assert.Equal(t, types.Proxies{url.URL{Scheme: "http", Host: "base-proxy.example.com:8080"}}, result.Proxies)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a copy of the base config for each test
			baseURL := func() *types.RequestURL { u := &types.RequestURL{}; u.Set("https://example.com"); return u }()
			testConfig := &Config{
				Method:       utils.ToPtr(*baseConfig.Method),
				URL:          baseURL,
				Timeout:      &types.Timeout{Duration: baseConfig.Timeout.Duration},
				DodosCount:   utils.ToPtr(*baseConfig.DodosCount),
				RequestCount: utils.ToPtr(*baseConfig.RequestCount),
				Duration:     &types.Duration{Duration: baseConfig.Duration.Duration},
				Yes:          utils.ToPtr(*baseConfig.Yes),
				Params:       slices.Clone(baseConfig.Params),
				Headers:      slices.Clone(baseConfig.Headers),
				Cookies:      slices.Clone(baseConfig.Cookies),
				Body:         slices.Clone(baseConfig.Body),
				Proxies:      slices.Clone(baseConfig.Proxies),
			}

			// Call the function being tested
			testConfig.MergeConfig(tt.newConfig)

			// Run assertions
			tt.assertions(t, testConfig)
		})
	}
}

func TestSetDefaults(t *testing.T) {
	tests := []struct {
		name     string
		config   *Config
		validate func(t *testing.T, config *Config)
	}{
		{
			name:   "empty config",
			config: &Config{},
			validate: func(t *testing.T, config *Config) {
				assert.Equal(t, DefaultMethod, *config.Method)
				assert.Equal(t, DefaultTimeout, config.Timeout.Duration)
				assert.Equal(t, DefaultDodosCount, *config.DodosCount)
				assert.Equal(t, DefaultRequestCount, *config.RequestCount)
				assert.Equal(t, DefaultDuration, config.Duration.Duration)
				assert.Equal(t, DefaultYes, *config.Yes)
				assert.True(t, config.Headers.Has("User-Agent"))
				userAgent := config.Headers.GetValue("User-Agent")
				assert.NotNil(t, userAgent)
				assert.Contains(t, (*userAgent)[0], DefaultUserAgent)
			},
		},
		{
			name: "partial config",
			config: &Config{
				Method:  utils.ToPtr("POST"),
				Timeout: &types.Timeout{Duration: 30 * time.Second},
				Headers: types.Headers{{Key: "Custom-Header", Value: []string{"value"}}},
			},
			validate: func(t *testing.T, config *Config) {
				assert.Equal(t, "POST", *config.Method)                    // should keep existing value
				assert.Equal(t, 30*time.Second, config.Timeout.Duration)   // should keep existing value
				assert.Equal(t, DefaultDodosCount, *config.DodosCount)     // should set default
				assert.Equal(t, DefaultRequestCount, *config.RequestCount) // should set default
				assert.Equal(t, DefaultDuration, config.Duration.Duration) // should set default
				assert.Equal(t, DefaultYes, *config.Yes)                   // should set default
				assert.True(t, config.Headers.Has("Custom-Header"))        // should keep existing header
				assert.True(t, config.Headers.Has("User-Agent"))           // should add User-Agent
				userAgent := config.Headers.GetValue("User-Agent")
				assert.NotNil(t, userAgent)
				assert.Contains(t, (*userAgent)[0], DefaultUserAgent)
			},
		},
		{
			name: "complete config",
			config: &Config{
				Method:       utils.ToPtr("DELETE"),
				URL:          func() *types.RequestURL { u := &types.RequestURL{}; u.Set("https://example.com"); return u }(),
				Timeout:      &types.Timeout{Duration: 15 * time.Second},
				DodosCount:   utils.ToPtr(uint(5)),
				RequestCount: utils.ToPtr(uint(500)),
				Duration:     &types.Duration{Duration: 5 * time.Minute},
				Yes:          utils.ToPtr(true),
				Headers:      types.Headers{{Key: "User-Agent", Value: []string{"CustomAgent"}}},
			},
			validate: func(t *testing.T, config *Config) {
				assert.Equal(t, "DELETE", *config.Method)
				assert.Equal(t, 15*time.Second, config.Timeout.Duration)
				assert.Equal(t, uint(5), *config.DodosCount)
				assert.Equal(t, uint(500), *config.RequestCount)
				assert.Equal(t, 5*time.Minute, config.Duration.Duration)
				assert.True(t, *config.Yes)
				assert.True(t, config.Headers.Has("User-Agent"))
				userAgent := config.Headers.GetValue("User-Agent")
				assert.NotNil(t, userAgent)
				assert.Equal(t, "CustomAgent", (*userAgent)[0])       // should keep custom user agent
				assert.NotEqual(t, DefaultUserAgent, (*userAgent)[0]) // should not overwrite
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the function being tested
			tt.config.SetDefaults()

			// Validate the result
			tt.validate(t, tt.config)
		})
	}
}
