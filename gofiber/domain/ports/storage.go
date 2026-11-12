package ports

import (
	"io"
	"time"
)

// StoragePort defines the interface for file storage operations
// Supports any S3-compatible storage: iDrive e2, Cloudflare R2, AWS S3, DO Spaces, etc.
type StoragePort interface {
	// UploadFile uploads a file to storage and returns the path (not full URL)
	UploadFile(file io.Reader, path string, contentType string) (string, error)

	// DeleteFile deletes a file from storage
	DeleteFile(path string) error

	// GetFileURL returns the public URL for a file path (for public buckets)
	GetFileURL(path string) string

	// GetPresignedURL generates a presigned URL for private bucket access
	// expiration: how long the URL should be valid (e.g., 15*time.Minute)
	GetPresignedURL(path string, expiration time.Duration) (string, error)
}
