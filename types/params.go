package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jedib0t/go-pretty/v6/text"
)

type Params []KeyValue[string, []string]

func (params Params) String() string {
	var buffer bytes.Buffer
	if len(params) == 0 {
		return string(buffer.Bytes())
	}

	indent := "  "

	displayLimit := 3

	for i, item := range params {
		if i > 0 {
			buffer.WriteString(",\n")
		}

		if len(item.Value) == 1 {
			buffer.WriteString(item.Key + ": " + item.Value[0])
			continue
		}
		buffer.WriteString(item.Key + ": " + text.FgBlue.Sprint("Random") + "[\n")

		for ii, v := range item.Value[:min(len(item.Value), displayLimit)] {
			if ii == len(item.Value)-1 {
				buffer.WriteString(indent + v + "\n")
			} else {
				buffer.WriteString(indent + v + ",\n")
			}
		}

		// Add remaining values count if needed
		if remainingValues := len(item.Value) - displayLimit; remainingValues > 0 {
			buffer.WriteString(indent + text.FgGreen.Sprintf("+%d values", remainingValues) + "\n")
		}

		buffer.WriteString("]")
	}

	// Add remaining key-value pairs count if needed
	if remainingPairs := len(params) - displayLimit; remainingPairs > 0 {
		buffer.WriteString(",\n" + text.FgGreen.Sprintf("+%d params", remainingPairs))
	}

	return string(buffer.Bytes())
}

func (params *Params) UnmarshalJSON(b []byte) error {
	var data []map[string]any
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	for _, item := range data {
		for key, value := range item {
			switch parsedValue := value.(type) {
			case string:
				*params = append(*params, KeyValue[string, []string]{Key: key, Value: []string{parsedValue}})
			case []any:
				parsedStr := make([]string, len(parsedValue))
				for i, item := range parsedValue {
					parsedStr[i] = fmt.Sprintf("%v", item)
				}
				*params = append(*params, KeyValue[string, []string]{Key: key, Value: parsedStr})
			default:
				return fmt.Errorf("unsupported type for params expected string or []string, got %T", parsedValue)
			}
		}
	}

	return nil
}

func (params *Params) Set(value string) error {
	parts := strings.SplitN(value, "=", 2)
	switch len(parts) {
	case 0:
		params.AppendByKey("", "")
	case 1:
		params.AppendByKey(parts[0], "")
	case 2:
		params.AppendByKey(parts[0], parts[1])
	}

	return nil
}

func (params *Params) AppendByKey(key string, value string) {
	if existingValue := params.GetValue(key); existingValue != nil {
		*params = append(*params, KeyValue[string, []string]{Key: key, Value: append(existingValue, value)})
	} else {
		*params = append(*params, KeyValue[string, []string]{Key: key, Value: []string{value}})
	}
}

func (params *Params) GetValue(key string) []string {
	for _, param := range *params {
		if param.Key == key {
			return param.Value
		}
	}
	return nil
}
