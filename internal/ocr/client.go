package ocr

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const MaxImageBytes = 5 * 1024 * 1024
const MaxResponseBytes = 1024 * 1024
const ModelID = "easyocr:en,ru"

type Result struct {
	Text  string `json:"text"`
	Model string `json:"model"`
}

type Client struct {
	URL  string
	http *http.Client
}

func NewClient(cfg Config) (*Client, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &Client{
		URL:  strings.TrimRight(cfg.URL, "/"),
		http: &http.Client{Timeout: cfg.Timeout},
	}, nil
}

func (c *Client) RecognizeImage(ctx context.Context, image io.Reader) (Result, error) {
	data, err := io.ReadAll(io.LimitReader(image, MaxImageBytes+1))
	if err != nil {
		return Result{}, fmt.Errorf("read OCR image: %w", err)
	}

	if len(data) == 0 || len(data) > MaxImageBytes {
		return Result{}, fmt.Errorf("OCR image must contain 1 to %d bytes", MaxImageBytes)
	}

	return c.request(ctx, bytes.NewReader(data))
}

func (c *Client) request(
	ctx context.Context,
	body io.Reader,
) (Result, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.URL, body)
	if err != nil {
		return Result{}, fmt.Errorf("create OCR request: %w", err)
	}

	req.Header.Set("Content-Type", "application/octet-stream")
	resp, err := c.http.Do(req)
	if err != nil {
		return Result{}, fmt.Errorf("call OCR: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Result{}, fmt.Errorf("status: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(io.LimitReader(resp.Body, MaxResponseBytes+1))
	if err != nil {
		return Result{}, fmt.Errorf("read OCR response: %w", err)
	}

	if len(data) > MaxResponseBytes {
		return Result{}, fmt.Errorf("OCR response exceeds 1 MiB")
	}

	var result Result
	if err := json.Unmarshal(data, &result); err != nil {
		return Result{}, fmt.Errorf("decode OCR response: %w", err)
	}

	if result.Model != ModelID {
		return Result{}, fmt.Errorf("unexpected OCR model: got %q, expected %q", result.Model, ModelID)
	}

	return result, nil
}
