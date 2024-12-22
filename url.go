package main

import (
	"fmt"
	"net/url"
	"time"
)

func getQueryParams(rawURL string) (url.Values, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %v", err)
	}

	queryParams := parsedURL.Query()
	return queryParams, nil
}

func GetHashURL(rawURL string) (string, error) {
	urlParsed, _ := url.Parse(rawURL)
	queryParams := urlParsed.Query()
	hash := queryParams.Get("hm")
	if hash == "" {
		return "", fmt.Errorf("no hash found")
	}

	return hash, nil
}

func GetExURL(rawURL string) (*time.Time, error) {
	urlParsed, _ := url.Parse(rawURL)
	queryParams := urlParsed.Query()
	ex := queryParams.Get("ex")

	if ex == "" {
		return nil, fmt.Errorf("no expiration time found")
	}

	exTime, err := parseHexToUnix(ex)
	if err != nil {
		return nil, fmt.Errorf("failed to parse expiration time: %v", err)
	}

	return &exTime, nil
}

func GetIsURL(rawURL string) (*time.Time, error) {
	urlParsed, _ := url.Parse(rawURL)
	queryParams := urlParsed.Query()
	is := queryParams.Get("is")

	if is == "" {
		return nil, fmt.Errorf("no issue time found")
	}

	isTime, err := parseHexToUnix(is)
	if err != nil {
		return nil, fmt.Errorf("failed to parse issue time: %v", err)
	}

	return &isTime, nil
}
