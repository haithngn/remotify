package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
)

func PostHTTP(endpoint, body string, header map[string][]string) ([]byte, map[string][]string, error) {
	log.Debug(endpoint, header, body)
	var jsonStr = []byte(body)
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonStr))
	req.Header = header

	client := &http.Client{}

	log.Debug("HTTP POST: ", endpoint, req.Header, req.Body)

	resp, err := client.Do(req)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to post HTTP (%s)", err.Error()))
		return nil, nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error("Failed to close response responseData")
		}
	}(resp.Body)

	log.Debug("Response Status:", resp.Status)
	log.Debug("Response Headers:", resp.Header)
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("Failed to read response responseData")
		return nil, nil, err
	}

	return responseData, resp.Header, nil
}

func GetHTTP(endpoint, params string, header map[string][]string) ([]byte, map[string][]string, error) {
	var parameters url.Values
	err := json.Unmarshal([]byte(params), &parameters)
	if err != nil {
		return nil, nil, err
	}

	u := endpoint + parameters.Encode()
	log.Debug("HTTP GET: ", u, header, params)

	// Make the GET request
	resp, err := http.Get(u)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to post HTTP (%s)", err.Error()))
		return nil, nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error("Failed to close response responseData")
		}
	}(resp.Body)

	log.Debug("Response Status:", resp.Status)
	log.Debug("Response Headers:", resp.Header)

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("Failed to read response responseData")
		return nil, nil, err
	}

	return responseData, resp.Header, nil
}
