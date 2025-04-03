package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/jedib0t/go-pretty/v6/text"
)

type Proxies []url.URL

func (proxies Proxies) String() string {
	var buffer bytes.Buffer
	if len(proxies) == 0 {
		return buffer.String()
	}

	if len(proxies) == 1 {
		buffer.WriteString(proxies[0].String())
		return buffer.String()
	}

	buffer.WriteString(text.FgBlue.Sprint("Random") + "[\n")

	indent := "  "

	displayLimit := 5

	for i, item := range proxies[:min(len(proxies), displayLimit)] {
		if i > 0 {
			buffer.WriteString(",\n")
		}

		buffer.WriteString(indent + item.String())
	}

	// Add remaining count if there are more items
	if remainingValues := len(proxies) - displayLimit; remainingValues > 0 {
		buffer.WriteString(",\n" + indent + text.FgGreen.Sprintf("+%d proxies", remainingValues))
	}

	buffer.WriteString("\n]")
	return buffer.String()
}

func (proxies *Proxies) UnmarshalJSON(b []byte) error {
	var data any
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	switch v := data.(type) {
	case string:
		parsed, err := url.Parse(v)
		if err != nil {
			return err
		}
		*proxies = []url.URL{*parsed}
	case []any:
		var urls []url.URL
		for _, item := range v {
			url, err := url.Parse(item.(string))
			if err != nil {
				return err
			}
			urls = append(urls, *url)
		}
		*proxies = urls
	default:
		return fmt.Errorf("invalid type for Body: %T (should be URL or []URL)", v)
	}

	return nil
}

func (proxies *Proxies) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var data any
	if err := unmarshal(&data); err != nil {
		return err
	}

	switch v := data.(type) {
	case string:
		parsed, err := url.Parse(v)
		if err != nil {
			return err
		}
		*proxies = []url.URL{*parsed}
	case []any:
		var urls []url.URL
		for _, item := range v {
			url, err := url.Parse(item.(string))
			if err != nil {
				return err
			}
			urls = append(urls, *url)
		}
		*proxies = urls
	default:
		return fmt.Errorf("invalid type for Body: %T (should be URL or []URL)", v)
	}

	return nil
}

func (proxies *Proxies) Set(value string) error {
	parsedURL, err := url.Parse(value)
	if err != nil {
		return err
	}

	*proxies = append(*proxies, *parsedURL)
	return nil
}
