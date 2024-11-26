package readers

import (
	"encoding/json"
	"os"

	"github.com/aykhans/dodo/config"
	customerrors "github.com/aykhans/dodo/custom_errors"
)

func JSONConfigReader(filePath string) (*config.JSONConfig, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, customerrors.OSErrorFormater(err)
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
