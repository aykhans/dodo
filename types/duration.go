package types

import (
	"encoding/json"
	"errors"
	"time"
)

type Timeout struct {
	time.Duration
}

func (timeout *Timeout) UnmarshalJSON(b []byte) error {
	var v any
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		timeout.Duration = time.Duration(value)
		return nil
	case string:
		var err error
		timeout.Duration, err = time.ParseDuration(value)
		if err != nil {
			return errors.New("Timeout is invalid (e.g. 400ms, 1s, 5m, 1h)")
		}
		return nil
	default:
		return errors.New("Timeout is invalid (e.g. 400ms, 1s, 5m, 1h)")
	}
}

func (timeout Timeout) MarshalJSON() ([]byte, error) {
	return json.Marshal(timeout.Duration.String())
}
