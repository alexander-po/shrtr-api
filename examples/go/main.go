// Minimal, zero-dependency Go client for the Shrtr URL-shortener API.
//
// Docs:  https://shrtr.top/api
// Spec:  https://shrtr.top/openapi.json
//
// The API is free and anonymous (no key, no signup). Errors are RFC 7807
// problem+json and surface here as *ProblemError. Standard library only.
//
// Run: go run .
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

var base = envOr("SHRTR_BASE", "https://shrtr.top/api/v1")

var client = &http.Client{Timeout: 10 * time.Second}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// ShortLink is the POST /shorten success payload.
type ShortLink struct {
	Code        string `json:"code"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	CreatedAt   string `json:"created_at"`
}

// Stats is the GET /stats/{code} success payload.
type Stats struct {
	Code          string `json:"code"`
	ClicksCount   int    `json:"clicks_count"`
	LastClickedAt string `json:"last_clicked_at"`
	CreatedAt     string `json:"created_at"`
	Enabled       bool   `json:"enabled"`
}

// ProblemError wraps a non-2xx RFC 7807 response.
type ProblemError struct {
	Status int
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

func (e *ProblemError) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("shrtr: HTTP %d: %s", e.Status, e.Detail)
	}
	return fmt.Sprintf("shrtr: HTTP %d: %s", e.Status, e.Title)
}

func do(method, path string, body any, out any) error {
	var reader io.Reader
	if body != nil {
		buf, _ := json.Marshal(body)
		reader = bytes.NewReader(buf)
	}
	req, err := http.NewRequest(method, base+path, reader)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 300 {
		pe := &ProblemError{Status: resp.StatusCode}
		_ = json.Unmarshal(data, pe)
		return pe
	}
	if out != nil {
		return json.Unmarshal(data, out)
	}
	return nil
}

// Shorten creates a short link. Pass an empty alias for a random code.
func Shorten(url, alias string) (*ShortLink, error) {
	body := map[string]string{"url": url}
	if alias != "" {
		body["alias"] = alias
	}
	var link ShortLink
	if err := do(http.MethodPost, "/shorten", body, &link); err != nil {
		return nil, err
	}
	return &link, nil
}

// GetStats returns aggregate click stats for a short link.
func GetStats(code string) (*Stats, error) {
	var s Stats
	if err := do(http.MethodGet, "/stats/"+code, nil, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

func main() {
	link, err := Shorten("https://example.com/some/long/path", "")
	if err != nil {
		panic(err)
	}
	fmt.Println("short_url:", link.ShortURL)

	s, err := GetStats(link.Code)
	if err != nil {
		panic(err)
	}
	fmt.Printf("stats:     %d clicks\n", s.ClicksCount)
}
