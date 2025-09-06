// Package storage provides file and object storage infrastructure
// for managing artifacts, documents, and other binary data.
//
// The package includes:
// - S3-compatible object storage client
// - Local filesystem storage adapter
// - Multipart upload support
// - Pre-signed URL generation
// - Storage abstraction layer
package storage

import "time"

// Storage type constants
const (
	// TypeS3 indicates S3-compatible object storage
	TypeS3 = "s3"

	// TypeFilesystem indicates local filesystem storage
	TypeFilesystem = "filesystem"

	// TypeGCS indicates Google Cloud Storage
	TypeGCS = "gcs"

	// TypeAzure indicates Azure Blob Storage
	TypeAzure = "azure"
)

// S3 configuration constants
const (
	// DefaultRegion is the default AWS region
	DefaultRegion = "us-east-1"

	// DefaultBucket is the default storage bucket name
	DefaultBucket = "archesai"

	// DefaultEndpoint is the default S3 endpoint (empty for AWS)
	DefaultEndpoint = ""

	// DefaultUploadPartSize is the default multipart upload size (5MB)
	DefaultUploadPartSize = 5 * 1024 * 1024

	// DefaultUploadConcurrency is the default upload concurrency
	DefaultUploadConcurrency = 5

	// DefaultPresignExpiration is the default presigned URL expiration
	DefaultPresignExpiration = 15 * time.Minute
)

// Filesystem storage constants
const (
	// DefaultStoragePath is the default local storage path
	DefaultStoragePath = "./storage"

	// DefaultFilePermissions is the default file permissions
	DefaultFilePermissions = 0644

	// DefaultDirPermissions is the default directory permissions
	DefaultDirPermissions = 0755
)

// Storage limits
const (
	// MaxFileSize is the maximum allowed file size (100MB)
	MaxFileSize = 100 * 1024 * 1024

	// MaxFilenameLength is the maximum filename length
	MaxFilenameLength = 255

	// MaxPathLength is the maximum path length
	MaxPathLength = 4096
)

// Content types
const (
	// ContentTypeJSON is the JSON content type
	ContentTypeJSON = "application/json"

	// ContentTypeText is the plain text content type
	ContentTypeText = "text/plain"

	// ContentTypeBinary is the binary content type
	ContentTypeBinary = "application/octet-stream"

	// ContentTypePDF is the PDF content type
	ContentTypePDF = "application/pdf"

	// ContentTypeImage is the generic image content type
	ContentTypeImage = "image/*"
)
