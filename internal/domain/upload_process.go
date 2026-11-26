package domain

type UploadProcess string

const (
	UploadProcessNormalization UploadProcess = "normalization"
	UploadProcessOCR           UploadProcess = "ocr"
	UploadProcessEmbedding     UploadProcess = "embedding"
)
