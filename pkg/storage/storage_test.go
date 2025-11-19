package storage

import (
	"testing"
	"time"
)

func TestStorageTypeConstants(t *testing.T) {
	// Test storage type constants
	if TypeS3 != "s3" {
		t.Errorf("Expected TypeS3 to be 's3', got %s", TypeS3)
	}
	if TypeFilesystem != "filesystem" {
		t.Errorf("Expected TypeFilesystem to be 'filesystem', got %s", TypeFilesystem)
	}
	if TypeGCS != "gcs" {
		t.Errorf("Expected TypeGCS to be 'gcs', got %s", TypeGCS)
	}
	if TypeAzure != "azure" {
		t.Errorf("Expected TypeAzure to be 'azure', got %s", TypeAzure)
	}
}

func TestS3ConfigConstants(t *testing.T) {
	// Test S3 configuration constants
	if DefaultRegion != "us-east-1" {
		t.Errorf("Expected DefaultRegion to be 'us-east-1', got %s", DefaultRegion)
	}
	if DefaultBucket != "archesai" {
		t.Errorf("Expected DefaultBucket to be 'archesai', got %s", DefaultBucket)
	}
	if DefaultEndpoint != "" {
		t.Errorf("Expected DefaultEndpoint to be empty, got %s", DefaultEndpoint)
	}
	if DefaultUploadPartSize != 5*1024*1024 {
		t.Errorf("Expected DefaultUploadPartSize to be 5MB, got %d", DefaultUploadPartSize)
	}
	if DefaultUploadConcurrency != 5 {
		t.Errorf("Expected DefaultUploadConcurrency to be 5, got %d", DefaultUploadConcurrency)
	}
	if DefaultPresignExpiration != 15*time.Minute {
		t.Errorf(
			"Expected DefaultPresignExpiration to be 15 minutes, got %v",
			DefaultPresignExpiration,
		)
	}
}

func TestFilesystemConstants(t *testing.T) {
	// Test filesystem storage constants
	if DefaultStoragePath != "./storage" {
		t.Errorf("Expected DefaultStoragePath to be './storage', got %s", DefaultStoragePath)
	}
	if DefaultFilePermissions != 0644 {
		t.Errorf("Expected DefaultFilePermissions to be 0644, got %o", DefaultFilePermissions)
	}
	if DefaultDirPermissions != 0755 {
		t.Errorf("Expected DefaultDirPermissions to be 0755, got %o", DefaultDirPermissions)
	}
}

func TestStorageLimits(t *testing.T) {
	// Test storage limit constants
	if MaxFileSize != 100*1024*1024 {
		t.Errorf("Expected MaxFileSize to be 100MB, got %d", MaxFileSize)
	}
	if MaxFilenameLength != 255 {
		t.Errorf("Expected MaxFilenameLength to be 255, got %d", MaxFilenameLength)
	}
	if MaxPathLength != 4096 {
		t.Errorf("Expected MaxPathLength to be 4096, got %d", MaxPathLength)
	}
}

func TestContentTypes(t *testing.T) {
	// Test content type constants
	tests := []struct {
		name     string
		got      string
		expected string
	}{
		{
			name:     "JSON content type",
			got:      ContentTypeJSON,
			expected: "application/json",
		},
		{
			name:     "Text content type",
			got:      ContentTypeText,
			expected: "text/plain",
		},
		{
			name:     "Binary content type",
			got:      ContentTypeBinary,
			expected: "application/octet-stream",
		},
		{
			name:     "PDF content type",
			got:      ContentTypePDF,
			expected: "application/pdf",
		},
		{
			name:     "Image content type",
			got:      ContentTypeImage,
			expected: "image/*",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("Expected %s to be '%s', got '%s'", tt.name, tt.expected, tt.got)
			}
		})
	}
}

func TestFileSizeValidation(t *testing.T) {
	tests := []struct {
		name      string
		size      int
		wantValid bool
	}{
		{
			name:      "small file",
			size:      1024,
			wantValid: true,
		},
		{
			name:      "medium file",
			size:      10 * 1024 * 1024,
			wantValid: true,
		},
		{
			name:      "max size file",
			size:      MaxFileSize,
			wantValid: true,
		},
		{
			name:      "oversized file",
			size:      MaxFileSize + 1,
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := tt.size <= MaxFileSize
			if isValid != tt.wantValid {
				t.Errorf("File size %d validation = %v, want %v", tt.size, isValid, tt.wantValid)
			}
		})
	}
}

func TestFilenameValidation(t *testing.T) {
	tests := []struct {
		name      string
		filename  string
		wantValid bool
	}{
		{
			name:      "short filename",
			filename:  "test.txt",
			wantValid: true,
		},
		{
			name:      "max length filename",
			filename:  string(make([]byte, MaxFilenameLength)),
			wantValid: true,
		},
		{
			name:      "oversized filename",
			filename:  string(make([]byte, MaxFilenameLength+1)),
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := len(tt.filename) <= MaxFilenameLength
			if isValid != tt.wantValid {
				t.Errorf(
					"Filename length %d validation = %v, want %v",
					len(tt.filename),
					isValid,
					tt.wantValid,
				)
			}
		})
	}
}

func TestPathValidation(t *testing.T) {
	tests := []struct {
		name      string
		path      string
		wantValid bool
	}{
		{
			name:      "short path",
			path:      "/tmp/test.txt",
			wantValid: true,
		},
		{
			name:      "nested path",
			path:      "/var/lib/archesai/storage/artifacts/2024/01/file.txt",
			wantValid: true,
		},
		{
			name:      "max length path",
			path:      string(make([]byte, MaxPathLength)),
			wantValid: true,
		},
		{
			name:      "oversized path",
			path:      string(make([]byte, MaxPathLength+1)),
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := len(tt.path) <= MaxPathLength
			if isValid != tt.wantValid {
				t.Errorf(
					"Path length %d validation = %v, want %v",
					len(tt.path),
					isValid,
					tt.wantValid,
				)
			}
		})
	}
}

func TestPresignExpirationBounds(t *testing.T) {
	tests := []struct {
		name        string
		expiration  time.Duration
		wantValid   bool
		description string
	}{
		{
			name:        "default expiration",
			expiration:  DefaultPresignExpiration,
			wantValid:   true,
			description: "default 15 minutes",
		},
		{
			name:        "short expiration",
			expiration:  1 * time.Minute,
			wantValid:   true,
			description: "1 minute",
		},
		{
			name:        "long expiration",
			expiration:  7 * 24 * time.Hour,
			wantValid:   true,
			description: "7 days",
		},
		{
			name:        "zero expiration",
			expiration:  0,
			wantValid:   false,
			description: "invalid zero",
		},
		{
			name:        "negative expiration",
			expiration:  -1 * time.Hour,
			wantValid:   false,
			description: "invalid negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := tt.expiration > 0
			if isValid != tt.wantValid {
				t.Errorf("Expiration %v (%s) validation = %v, want %v",
					tt.expiration, tt.description, isValid, tt.wantValid)
			}
		})
	}
}

func BenchmarkFileSizeValidation(b *testing.B) {
	fileSize := 50 * 1024 * 1024 // 50MB

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = fileSize <= MaxFileSize
	}
}

func BenchmarkPathValidation(b *testing.B) {
	path := "/var/lib/archesai/storage/artifacts/2024/01/very/long/path/to/file.txt"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = len(path) <= MaxPathLength
	}
}
