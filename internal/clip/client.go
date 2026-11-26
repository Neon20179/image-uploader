package clip

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
	"net/http"
	"strings"
	"syscall"

	"image-uploader/internal/apperr"
)

const Dimensions = 512
const MaxImageBytes = 5 * 1024 * 1024
const MaxResponseBytes = 64 * 1024
const ModelID = "ViT-B-32:laion2b_s34b_b79k"

type Result struct {
	Model      string    `json:"model"`
	Dimensions int       `json:"dimensions"`
	Vector     []float32 `json:"vector"`
}

type Client struct {
	URL  string
	http *http.Client
}

func NewClient(cfg Config) (*Client, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &Client{URL: strings.TrimRight(cfg.URL, "/"), http: &http.Client{Timeout: cfg.Timeout}}, nil
}

func (c *Client) EncodeImage(ctx context.Context, image io.Reader) (Result, error) {
	data, err := io.ReadAll(io.LimitReader(image, MaxImageBytes+1))
	if err != nil {
		return Result{}, fmt.Errorf("read CLIP image: %w", err)
	}

	if len(data) == 0 || len(data) > MaxImageBytes {
		return Result{}, fmt.Errorf("CLIP image must contain 1 to %d bytes", MaxImageBytes)
	}

	return c.request(ctx, bytes.NewReader(data))
}

func (c *Client) request(
	ctx context.Context,
	body io.Reader,
) (Result, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.URL, body)
	if err != nil {
		return Result{}, fmt.Errorf("create CLIP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/octet-stream")
	resp, err := c.http.Do(req)
	if err != nil {
		if ctx.Err() == nil && noResponse(err) {
			err = apperr.NewUnavailableError(err)
		}
		return Result{}, fmt.Errorf("call CLIP: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Result{}, fmt.Errorf("status: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(io.LimitReader(resp.Body, MaxResponseBytes+1))
	if err != nil {
		if ctx.Err() == nil && noResponse(err) {
			err = apperr.NewUnavailableError(err)
		}
		return Result{}, fmt.Errorf("read CLIP response: %w", err)
	}
	if len(data) > MaxResponseBytes {
		return Result{}, fmt.Errorf("CLIP response exceeds 64 KiB")
	}

	var result Result
	if err := json.Unmarshal(data, &result); err != nil {
		return Result{}, fmt.Errorf("decode CLIP response: %w", err)
	}

	if result.Model != ModelID {
		return Result{}, fmt.Errorf("unexpected CLIP model")
	}

	if result.Dimensions != Dimensions || len(result.Vector) != Dimensions {
		return Result{}, fmt.Errorf("unexpected vector dimensions")
	}

	var norm float64
	for _, value := range result.Vector {
		v := float64(value)
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return Result{}, fmt.Errorf("CLIP vector contains a non-finite value")
		}
		norm += v * v
	}
	if math.Abs(norm-1) > 0.01 {
		return Result{}, fmt.Errorf("CLIP vector is not normalized")
	}

	return result, nil
}

func noResponse(err error) bool {
	var op *net.OpError
	var dns *net.DNSError
	var network net.Error
	return errors.As(err, &op) || errors.As(err, &dns) ||
		(errors.As(err, &network) && network.Timeout()) ||
		errors.Is(err, context.DeadlineExceeded) ||
		errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) ||
		errors.Is(err, syscall.ECONNRESET) || errors.Is(err, syscall.ECONNREFUSED) ||
		errors.Is(err, syscall.EPIPE)
}
