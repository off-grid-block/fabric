package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// accepts a payload in any form
func SendRequest_POST(url, endpoint string, payload interface{}) (*http.Response, error) {

	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	err := encoder.Encode(payload)
	if err != nil {
		return nil, fmt.Errorf("Error while encoding: %v\n", err)
	}

	// create request for registering schema
	req, err := http.NewRequest(
		"POST", 
		url + endpoint,
		buffer,
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to create new post request: %v\n", err)
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// create http client & send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to execute http request: %v\n", err)
	}

	return resp, nil
}


// accepts a payload in any form
func SendRequestWithParams_POST(
	url, endpoint string,
	params map[string]string,
	payload interface{}) (*http.Response, error) {

	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	err := encoder.Encode(payload)
	if err != nil {
		return nil, fmt.Errorf("Error while encoding: %v\n", err)
	}

	// create request for registering schema
	req, err := http.NewRequest(
		"POST",
		url + endpoint,
		buffer,
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to create new post request: %v\n", err)
	}

	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()
	req.Header.Set("accept", "application/json")

	// create http client & send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to execute http request: %v\n", err)
	}

	return resp, nil
}


func SendRequest_GET(url, endpoint string, params map[string]string) (*http.Response, error) {

	req, err := http.NewRequest("GET", url + endpoint, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()
	req.Header.Set("accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}