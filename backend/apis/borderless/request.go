package borderless

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// Borderless handles HTTP communication with headers and a client.
type Borderless struct {
	accessToken  string
	clientID     string
	clientSecret string
	accountID    string
	BaseUrl      string
	Client       *http.Client
	Headers      map[string]interface{}
	Timeout      time.Duration
}

// MakeRequest sends an HTTP request with retries and handles JSON responses.
func (hc *Borderless) MakeRequest(method, url string, data map[string]interface{}) (map[string]interface{}, error) {
	start := time.Now()

	req, err := hc.buildRequest(method, url, data)
	if err != nil {
		return nil, err
	}

	resp, err := hc.doWithRetry(req, 3)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return hc.handleErrorResponse(bodyBytes, resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	log.Printf("✅ %s %s succeeded in %v", method, url, time.Since(start))
	return result, nil
}

// buildRequest prepares an HTTP request with headers and optional JSON body.
func (hc *Borderless) buildRequest(method, url string, data map[string]interface{}) (*http.Request, error) {
	var body []byte
	var err error

	if data != nil {
		body, err = json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	for key, val := range hc.Headers {
		if strVal, ok := val.(string); ok {
			req.Header.Set(key, strVal)
		}
	}
	return req, nil
}

// doWithRetry attempts the request with a retry mechanism.
func (hc *Borderless) doWithRetry(req *http.Request, maxRetries int) (*http.Response, error) {
	if hc.Client == nil {
		hc.Client = &http.Client{}
	}

	var resp *http.Response
	var err error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		resp, err = hc.Client.Do(req)
		if err == nil {
			return resp, nil
		}

		if attempt < maxRetries {
			log.Printf("⚠️ Retry %d/%d after error: %v", attempt, maxRetries-1, err)
			time.Sleep(2 * time.Second)
		}
	}
	return nil, err
}

// handleErrorResponse parses and returns a meaningful error from a failed HTTP response.
func (hc *Borderless) handleErrorResponse(body []byte, statusCode int) (map[string]interface{}, error) {
	var errResp map[string]interface{}
	if err := json.Unmarshal(body, &errResp); err == nil {
		msg := "HTTP error"
		if m, ok := errResp["message"].(string); ok {
			msg = m
		}
		log.Printf("❌ %d error: %s", statusCode, msg)
		return errResp, fmt.Errorf("HTTP %d error: %s", statusCode, msg)
	}

	log.Printf("❌ %d error with unparseable body: %s", statusCode, string(body))
	return nil, fmt.Errorf("HTTP %d error: %s", statusCode, string(body))
}
