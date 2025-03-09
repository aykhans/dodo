package readers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"io"

	"github.com/aykhans/dodo/config"
	customerrors "github.com/aykhans/dodo/custom_errors"
)

func JSONConfigReader(filePath string) (*config.JSONConfig, error) {
	var (
		data []byte
		err  error
	)

	if strings.HasPrefix(filePath, "http://") || strings.HasPrefix(filePath, "https://") {
		client := &http.Client{
			Timeout: 10 * time.Second,
		}

		resp, err := client.Get(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch JSON config from %s", filePath)
		}
		defer resp.Body.Close()

		data, err = io.ReadAll(io.Reader(resp.Body))
		if err != nil {
			return nil, fmt.Errorf("failed to read JSON config from %s", filePath)
		}
	} else {
		data, err = os.ReadFile(filePath)
		if err != nil {
			return nil, customerrors.OSErrorFormater(err)
		}
	}

	jsonConf := config.NewJSONConfig(
		config.NewConfig("", 0, 0, 0, nil),
		nil, nil, nil, nil, nil,
	)
	err = json.Unmarshal(data, &jsonConf)

	if err != nil {
		switch err := err.(type) {
		case *json.UnmarshalTypeError:
			return nil,
				customerrors.NewTypeError(
					err.Type.String(),
					err.Value,
					err.Field,
					err,
				)
		}
		return nil, customerrors.NewInvalidFileError(filePath, err)
	}

	return jsonConf, nil
}
