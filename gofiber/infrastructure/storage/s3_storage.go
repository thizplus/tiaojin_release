package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gofiber-template/domain/ports"
)

// S3Config holds configuration for S3-compatible storage
// Works with: iDrive e2, Cloudflare R2, AWS S3, DigitalOcean Spaces, MinIO, etc.
type S3Config struct {
	Endpoint        string // e.g., "https://xxxx.r2.cloudflarestorage.com" for Cloudflare R2
	Region          string // e.g., "us-east-1" or "auto"
	AccessKeyID     string
	SecretAccessKey string
	Bucket          string
	CDNUrl          string // Public URL prefix for accessing files
}

// S3Storage implements StoragePort for S3-compatible storage
type S3Storage struct {
	client        *s3.Client
	presignClient *s3.PresignClient
	bucket        string
	cdnURL        string
}

// NewS3Storage creates a new S3Storage instance
func NewS3Storage(cfg S3Config) (ports.StoragePort, error) {
	// Ensure endpoint has https:// prefix
	endpoint := cfg.Endpoint
	if endpoint != "" && !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") {
		endpoint = "https://" + endpoint
	}

	// Create custom resolver for S3-compatible endpoints
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if endpoint != "" {
			return aws.Endpoint{
				URL:               endpoint,
				SigningRegion:     cfg.Region,
				HostnameImmutable: true,
			}, nil
		}
		// Return default endpoint for AWS S3
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	// Load AWS config with custom credentials and endpoint
	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AccessKeyID,
			cfg.SecretAccessKey,
			"",
		)),
		config.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create S3 client with path-style addressing for compatibility
	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true // Required for most S3-compatible services
	})

	// Create presign client for generating presigned URLs
	presignClient := s3.NewPresignClient(client)

	return &S3Storage{
		client:        client,
		presignClient: presignClient,
		bucket:        cfg.Bucket,
		cdnURL:        strings.TrimSuffix(cfg.CDNUrl, "/"),
	}, nil
}

// UploadFile uploads a file to S3-compatible storage
// Returns the path only (not full URL) for database storage
// Use GetFileURL(path) to construct full URL when displaying
func (s *S3Storage) UploadFile(file io.Reader, path string, contentType string) (string, error) {
	// Ensure path doesn't start with /
	path = strings.TrimPrefix(path, "/")

	// Read file content
	content, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// Upload to S3 (private bucket - use presigned URLs for access)
	_, err = s.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(path),
		Body:        bytes.NewReader(content),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	// Return path only, not full URL
	// Full URL = S3_PUBLIC_URL + "/" + path (construct in DTO/response)
	return path, nil
}

// DeleteFile deletes a file from S3-compatible storage
func (s *S3Storage) DeleteFile(path string) error {
	// Ensure path doesn't start with /
	path = strings.TrimPrefix(path, "/")

	_, err := s.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file from S3: %w", err)
	}

	return nil
}

// GetFileURL returns the public URL for a file (for public buckets)
// For private buckets, use GetPresignedURL instead
func (s *S3Storage) GetFileURL(path string) string {
	// Ensure path starts with /
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return s.cdnURL + path
}

// GetPresignedURL generates a presigned URL for private bucket access
// expiration: how long the URL should be valid (e.g., 15*time.Minute, 1*time.Hour)
func (s *S3Storage) GetPresignedURL(path string, expiration time.Duration) (string, error) {
	// Ensure path doesn't start with /
	path = strings.TrimPrefix(path, "/")

	// Skip if path is empty
	if path == "" {
		return "", nil
	}

	// Generate presigned URL
	presignedReq, err := s.presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expiration
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return presignedReq.URL, nil
}
