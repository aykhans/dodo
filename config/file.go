package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/aykhans/dodo/types"
	"gopkg.in/yaml.v3"
)

var supportedFileTypes = []string{"json", "yaml", "yml"}

func (config *Config) ReadFile(filePath types.ConfigFile) error {
	var (
		data []byte
		err  error
	)

	fileExt := filePath.Extension()
	if slices.Contains(supportedFileTypes, fileExt) {
		if filePath.LocationType() == types.FileLocationTypeRemoteHTTP {
			client := &http.Client{
				Timeout: 10 * time.Second,
			}

			resp, err := client.Get(filePath.String())
			if err != nil {
				return fmt.Errorf("failed to fetch config file from %s", filePath)
			}
			defer func() { _ = resp.Body.Close() }()

			data, err = io.ReadAll(io.Reader(resp.Body))
			if err != nil {
				return fmt.Errorf("failed to read config file from %s", filePath)
			}
		} else {
			data, err = os.ReadFile(filePath.String())
			if err != nil {
				return errors.New("failed to read config file from " + filePath.String())
			}
		}

		switch fileExt {
		case "json":
			return parseJSONConfig(data, config)
		case "yml", "yaml":
			return parseYAMLConfig(data, config)
		}
	}

	return fmt.Errorf("unsupported config file type (supported types: %v)", strings.Join(supportedFileTypes, ", "))
}

func parseJSONConfig(data []byte, config *Config) error {
	err := json.Unmarshal(data, &config)
	if err != nil {
		switch parsedErr := err.(type) {
		case *json.SyntaxError:
			return fmt.Errorf("JSON Config file: invalid syntax at byte offset %d", parsedErr.Offset)
		case *json.UnmarshalTypeError:
			return fmt.Errorf("JSON Config file: invalid type %v for field %s, expected %v", parsedErr.Value, parsedErr.Field, parsedErr.Type)
		default:
			return fmt.Errorf("JSON Config file: %s", err.Error())
		}
	}

	return nil
}

func parseYAMLConfig(data []byte, config *Config) error {
	err := yaml.Unmarshal(data, &config)
	if err != nil {
		return fmt.Errorf("YAML Config file: %s", err.Error())
	}

	return nil
}
