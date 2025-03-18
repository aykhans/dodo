package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/aykhans/dodo/types"
)

func (config *Config) ReadFile(filePath types.ConfigFile) error {
	var (
		data []byte
		err  error
	)

	if filePath.LocationType() == types.FileLocationTypeRemoteHTTP {
		client := &http.Client{
			Timeout: 10 * time.Second,
		}

		resp, err := client.Get(filePath.String())
		if err != nil {
			return fmt.Errorf("failed to fetch config file from %s", filePath)
		}
		defer resp.Body.Close()

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

	return parseJSONConfig(data, config)
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
