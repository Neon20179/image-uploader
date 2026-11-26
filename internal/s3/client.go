package s3

import (
	"context"
	"io"
	"net/url"
	"time"
)

type ObjectInfo struct {
	Size int64
}

type PostPolicyOptions struct {
	Expiry            time.Duration
	MinSizeBytes      int64
	MaxSizeBytes      int64
	ContentTypePrefix string
}

type Client interface {
	PresignedPostPolicy(
		ctx context.Context,
		bucket string,
		objectKey string,
		opts PostPolicyOptions,
	) (*url.URL, map[string]string, error)
	StatObject(ctx context.Context, bucket string, objectKey string) (ObjectInfo, error)
	ObjectContentType(ctx context.Context, bucket string, objectKey string) (string, error)
	GetObject(ctx context.Context, bucket string, objectKey string) (io.ReadCloser, error)
	PutObject(
		ctx context.Context,
		bucket string,
		objectKey string,
		r io.Reader,
		size int64,
		contentType string,
	) error
}
