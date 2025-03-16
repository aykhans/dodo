package types

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/jedib0t/go-pretty/v6/text"
)

type Body []string

func (body Body) String() string {
	var buffer bytes.Buffer
	if len(body) == 0 {
		return string(buffer.Bytes())
	}

	if len(body) == 1 {
		buffer.WriteString(body[0])
		return string(buffer.Bytes())
	}

	buffer.WriteString(text.FgBlue.Sprint("Random") + "[\n")

	indent := "  "

	displayLimit := 5

	for i, item := range body[:min(len(body), displayLimit)] {
		if i > 0 {
			buffer.WriteString(",\n")
		}

		buffer.WriteString(indent + item)
	}

	// Add remaining count if there are more items
	if remainingValues := len(body) - displayLimit; remainingValues > 0 {
		buffer.WriteString(",\n" + indent + text.FgGreen.Sprintf("+%d bodies", remainingValues))
	}

	buffer.WriteString("\n]")
	return string(buffer.Bytes())
}

func (body *Body) UnmarshalJSON(b []byte) error {
	var data any
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	switch v := data.(type) {
	case string:
		*body = []string{v}
	case []any:
		var slice []string
		for _, item := range v {
			slice = append(slice, fmt.Sprintf("%v", item))
		}
		*body = slice
	default:
		return fmt.Errorf("invalid type for Body: %T (should be string or []string)", v)
	}

	return nil
}

func (body *Body) Set(value string) error {
	*body = append(*body, value)
	return nil
}
