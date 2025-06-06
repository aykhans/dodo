package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jedib0t/go-pretty/v6/text"
)

type Cookies []KeyValue[string, []string]

func (cookies Cookies) String() string {
	var buffer bytes.Buffer
	if len(cookies) == 0 {
		return buffer.String()
	}

	indent := "  "

	displayLimit := 3

	for i, item := range cookies[:min(len(cookies), displayLimit)] {
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
	if remainingPairs := len(cookies) - displayLimit; remainingPairs > 0 {
		buffer.WriteString(",\n" + text.FgGreen.Sprintf("+%d cookies", remainingPairs))
	}

	return buffer.String()
}

func (cookies *Cookies) AppendByKey(key, value string) {
	if item := cookies.GetValue(key); item != nil {
		*item = append(*item, value)
	} else {
		*cookies = append(*cookies, KeyValue[string, []string]{Key: key, Value: []string{value}})
	}
}

func (cookies Cookies) GetValue(key string) *[]string {
	for i := range cookies {
		if cookies[i].Key == key {
			return &cookies[i].Value
		}
	}
	return nil
}

func (cookies *Cookies) UnmarshalJSON(b []byte) error {
	var data []map[string]any
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	for _, item := range data {
		for key, value := range item {
			switch parsedValue := value.(type) {
			case string:
				*cookies = append(*cookies, KeyValue[string, []string]{Key: key, Value: []string{parsedValue}})
			case []any:
				parsedStr := make([]string, len(parsedValue))
				for i, item := range parsedValue {
					parsedStr[i] = fmt.Sprintf("%v", item)
				}
				*cookies = append(*cookies, KeyValue[string, []string]{Key: key, Value: parsedStr})
			default:
				return fmt.Errorf("unsupported type for cookies expected string or []string, got %T", parsedValue)
			}
		}
	}

	return nil
}

func (cookies *Cookies) UnmarshalYAML(unmarshal func(any) error) error {
	var raw []map[string]any
	if err := unmarshal(&raw); err != nil {
		return err
	}

	for _, param := range raw {
		for key, value := range param {
			switch parsed := value.(type) {
			case string:
				*cookies = append(*cookies, KeyValue[string, []string]{Key: key, Value: []string{parsed}})
			case []any:
				var values []string
				for _, v := range parsed {
					if str, ok := v.(string); ok {
						values = append(values, str)
					}
				}
				*cookies = append(*cookies, KeyValue[string, []string]{Key: key, Value: values})
			}
		}
	}
	return nil
}

func (cookies *Cookies) Set(value string) error {
	parts := strings.SplitN(value, "=", 2)
	switch len(parts) {
	case 0:
		cookies.AppendByKey("", "")
	case 1:
		cookies.AppendByKey(parts[0], "")
	case 2:
		cookies.AppendByKey(parts[0], parts[1])
	}

	return nil
}
