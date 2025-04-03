package types

import (
	"encoding/json"
	"errors"
	"time"
)

type Duration struct {
	time.Duration
}

func (duration *Duration) UnmarshalJSON(b []byte) error {
	var v any
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		duration.Duration = time.Duration(value)
		return nil
	case string:
		var err error
		duration.Duration, err = time.ParseDuration(value)
		if err != nil {
			return errors.New("Duration is invalid (e.g. 400ms, 1s, 5m, 1h)")
		}
		return nil
	default:
		return errors.New("Duration is invalid (e.g. 400ms, 1s, 5m, 1h)")
	}
}

func (duration Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(duration.String())
}

func (duration *Duration) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var v any
	if err := unmarshal(&v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		duration.Duration = time.Duration(value)
		return nil
	case string:
		var err error
		duration.Duration, err = time.ParseDuration(value)
		if err != nil {
			return errors.New("Duration is invalid (e.g. 400ms, 1s, 5m, 1h)")
		}
		return nil
	default:
		return errors.New("Duration is invalid (e.g. 400ms, 1s, 5m, 1h)")
	}
}
