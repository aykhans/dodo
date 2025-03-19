package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jedib0t/go-pretty/v6/text"
)

type Headers []KeyValue[string, []string]

func (headers Headers) String() string {
	var buffer bytes.Buffer
	if len(headers) == 0 {
		return string(buffer.Bytes())
	}

	indent := "  "

	displayLimit := 3

	for i, item := range headers[:min(len(headers), displayLimit)] {
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
	if remainingPairs := len(headers) - displayLimit; remainingPairs > 0 {
		buffer.WriteString(",\n" + text.FgGreen.Sprintf("+%d headers", remainingPairs))
	}

	return string(buffer.Bytes())
}

func (headers *Headers) UnmarshalJSON(b []byte) error {
	var data []map[string]any
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	for _, item := range data {
		for key, value := range item {
			switch parsedValue := value.(type) {
			case string:
				*headers = append(*headers, KeyValue[string, []string]{Key: key, Value: []string{parsedValue}})
			case []any:
				parsedStr := make([]string, len(parsedValue))
				for i, item := range parsedValue {
					parsedStr[i] = fmt.Sprintf("%v", item)
				}
				*headers = append(*headers, KeyValue[string, []string]{Key: key, Value: parsedStr})
			default:
				return fmt.Errorf("unsupported type for headers expected string or []string, got %T", parsedValue)
			}
		}
	}

	return nil
}

func (headers *Headers) Set(value string) error {
	parts := strings.SplitN(value, ":", 2)
	switch len(parts) {
	case 0:
		headers.AppendByKey("", "")
	case 1:
		headers.AppendByKey(parts[0], "")
	case 2:
		headers.AppendByKey(parts[0], parts[1])
	}

	return nil
}

func (headers *Headers) AppendByKey(key, value string) {
	if item := headers.GetValue(key); item != nil {
		*item = append(*item, value)
	} else {
		*headers = append(*headers, KeyValue[string, []string]{Key: key, Value: []string{value}})
	}
}

func (headers Headers) GetValue(key string) *[]string {
	for i := range headers {
		if headers[i].Key == key {
			return &headers[i].Value
		}
	}
	return nil
}
