package cdn

import (
	"os"
	"strings"
)

// GetStorageURL generates URL through our API storage endpoint
// This is used for private S3 buckets that require presigned URLs
// path: "uploads/users/xxx/logo/yyy.jpg"
// returns: "/api/v1/storage/uploads/users/xxx/logo/yyy.jpg"
// The frontend will call this endpoint, and our server will redirect to a presigned URL
func GetStorageURL(path string) string {
	if path == "" {
		return ""
	}

	// If already a full URL, check if it's an old DO Spaces URL that needs migration
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		// Old DO Spaces URLs should be converted to path
		// For now, return as-is - migration script will handle this
		return path
	}

	// Ensure path doesn't start with /
	path = strings.TrimPrefix(path, "/")

	// Return API storage endpoint URL
	return "/api/v1/storage/" + path
}

// GetStorageURLPtr is a helper for pointer string fields
func GetStorageURLPtr(path *string) *string {
	if path == nil || *path == "" {
		return nil
	}
	url := GetStorageURL(*path)
	return &url
}

// GetFullFileURL constructs full URL from path using S3_PUBLIC_URL from env
// Use this ONLY for public S3 buckets
// For private buckets, use GetStorageURL instead
// path: "uploads/users/xxx/logo/yyy.jpg"
// returns: "https://cdn.example.com/uploads/users/xxx/logo/yyy.jpg"
func GetFullFileURL(path string) string {
	if path == "" {
		return ""
	}

	// If already a full URL, return as-is
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return path
	}

	// Check if we should use API storage endpoint (for private buckets)
	// Set S3_PRIVATE_BUCKET=true in .env to enable this
	if os.Getenv("S3_PRIVATE_BUCKET") == "true" {
		return GetStorageURL(path)
	}

	baseURL := strings.TrimSuffix(os.Getenv("S3_PUBLIC_URL"), "/")
	if baseURL == "" {
		return path // fallback to path if no base URL configured
	}

	// Ensure path doesn't start with /
	path = strings.TrimPrefix(path, "/")

	return baseURL + "/" + path
}

// GetFullFileURLPtr is a helper for pointer string fields
func GetFullFileURLPtr(path *string) *string {
	if path == nil || *path == "" {
		return nil
	}
	url := GetFullFileURL(*path)
	return &url
}
