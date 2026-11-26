package minio

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"image-uploader/internal/s3"
)

type Client struct {
	client    *minio.Client
	publicURL *url.URL
	opTimeout time.Duration
}

func NewClient(ctx context.Context, cfg Config) (*Client, error) {
	var publicURL *url.URL
	if cfg.PublicEndpoint != "" {
		scheme := "http"
		if cfg.PublicSecure {
			scheme = "https"
		}
		parsed, err := url.Parse(scheme + "://" + cfg.PublicEndpoint)
		if err != nil || parsed.Hostname() == "" || parsed.User != nil || parsed.Path != "" || parsed.RawQuery != "" || parsed.ForceQuery || parsed.Fragment != "" {
			return nil, fmt.Errorf("public MinIO endpoint must be a host with an optional port")
		}
		publicURL = parsed
	}
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.Secure,
		Region: cfg.Region,
	})
	if err != nil {
		return &Client{}, fmt.Errorf("initializing minio client: %w", err)
	}
	return &Client{
		client:    client,
		publicURL: publicURL,
		opTimeout: cfg.Timeout,
	}, nil
}

func (c *Client) PresignedPostPolicy(ctx context.Context, bucket string, objectKey string, opts s3.PostPolicyOptions) (*url.URL, map[string]string, error) {
	ctx, cancel := context.WithTimeout(ctx, c.opTimeout)
	defer cancel()

	policy := minio.NewPostPolicy()
	if err := policy.SetBucket(bucket); err != nil {
		return nil, nil, fmt.Errorf("set policy bucket: %w", err)
	}
	if err := policy.SetKey(objectKey); err != nil {
		return nil, nil, fmt.Errorf("set policy object key: %w", err)
	}
	expires := time.Now().Add(opts.Expiry)
	if err := policy.SetExpires(expires); err != nil {
		return nil, nil, fmt.Errorf("set policy expires: %w", err)
	}
	if err := policy.SetContentLengthRange(opts.MinSizeBytes, opts.MaxSizeBytes); err != nil {
		return nil, nil, fmt.Errorf("set policy content length range: %w", err)
	}
	if opts.ContentTypePrefix != "" {
		if err := policy.SetContentTypeStartsWith(opts.ContentTypePrefix); err != nil {
			return nil, nil, fmt.Errorf("set policy content type: %w", err)
		}
	}

	uploadURL, formData, err := c.client.PresignedPostPolicy(ctx, policy)
	if err != nil {
		return nil, nil, fmt.Errorf("create presigned post policy: %w", err)
	}

	if c.publicURL != nil {
		public := *c.publicURL
		public.Path = "/" + bucket + "/"
		uploadURL = &public
	}

	return uploadURL, formData, nil
}

func (c *Client) StatObject(ctx context.Context, bucket string, objectKey string) (s3.ObjectInfo, error) {
	ctx, cancel := context.WithTimeout(ctx, c.opTimeout)
	defer cancel()

	stat, err := c.client.StatObject(ctx, bucket, objectKey, minio.GetObjectOptions{})
	objectInfo := s3.ObjectInfo{
		Size: stat.Size,
	}
	return objectInfo, err
}

func (c *Client) ObjectContentType(ctx context.Context, bucket string, objectKey string) (string, error) {
	opts := minio.GetObjectOptions{}
	if err := opts.SetRange(0, 511); err != nil {
		return "", err
	}
	head, err := c.client.GetObject(ctx, bucket, objectKey, opts)
	if err != nil {
		return "", fmt.Errorf("get head: %w", err)
	}
	defer func() {
		_ = head.Close()
	}()

	buf := make([]byte, 512)
	n, err := io.ReadFull(head, buf)
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
		return "", fmt.Errorf("read head: %w", err)
	}

	mime := http.DetectContentType(buf[:n])
	return mime, nil
}

func (c *Client) GetObject(ctx context.Context, bucket string, objectKey string) (io.ReadCloser, error) {
	return c.client.GetObject(ctx, bucket, objectKey, minio.GetObjectOptions{})
}

func (c *Client) PutObject(
	ctx context.Context,
	bucket string,
	objectKey string,
	r io.Reader,
	size int64,
	contentType string,
) error {
	_, err := c.client.PutObject(ctx, bucket, objectKey, r, size, minio.PutObjectOptions{ContentType: contentType})
	return err
}
