package minio

import (
	"fmt"
	"strings"
)

type URLBuilder struct {
	baseURL string
}

func NewURLBuilder(publicBaseURL string) *URLBuilder {
	return &URLBuilder{baseURL: strings.TrimRight(publicBaseURL, "/")}
}

func (b *URLBuilder) Build(bucket, objectKey string) string {
	return fmt.Sprintf("%s/%s/%s", b.baseURL, bucket, objectKey)
}
