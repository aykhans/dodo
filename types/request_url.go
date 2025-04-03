package types

import (
	"encoding/json"
	"errors"
	"net/url"
)

type RequestURL struct {
	url.URL
}

func (requestURL *RequestURL) UnmarshalJSON(data []byte) error {
	var urlStr string
	if err := json.Unmarshal(data, &urlStr); err != nil {
		return err
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return errors.New("request URL is invalid")
	}

	requestURL.URL = *parsedURL
	return nil
}

func (requestURL *RequestURL) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var urlStr string
	if err := unmarshal(&urlStr); err != nil {
		return err
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return errors.New("request URL is invalid")
	}

	requestURL.URL = *parsedURL
	return nil
}

func (requestURL RequestURL) MarshalJSON() ([]byte, error) {
	return json.Marshal(requestURL.URL.String())
}

func (requestURL RequestURL) String() string {
	return requestURL.URL.String()
}

func (requestURL *RequestURL) Set(value string) error {
	parsedURL, err := url.Parse(value)
	if err != nil {
		return err
	}

	requestURL.URL = *parsedURL
	return nil
}
