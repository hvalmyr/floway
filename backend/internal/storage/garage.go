package storage

import (
	"context"
	"fmt"
	"io"
	"net/url"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Client struct {
	mc     *minio.Client
	bucket string
}

func NewClient(endpoint, region, accessKey, secretKey, bucket string) (*Client, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("parse garage endpoint: %w", err)
	}

	mc, err := minio.New(u.Host, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: u.Scheme == "https",
		Region: region,
	})
	if err != nil {
		return nil, fmt.Errorf("create garage client: %w", err)
	}

	return &Client{mc: mc, bucket: bucket}, nil
}

// EnsureBucket confirms the configured bucket actually exists and is
// reachable — NewClient/minio.New never make a network call (pure
// constructor), so a wrong bucket name or an unreachable Garage would
// otherwise only surface on the first admin upload, as a 500 with a raw S3
// error (architecture review finding #6). Call this once at startup to fail
// fast instead.
func (c *Client) EnsureBucket(ctx context.Context) error {
	ok, err := c.mc.BucketExists(ctx, c.bucket)
	if err != nil {
		return fmt.Errorf("check garage bucket %q: %w", c.bucket, err)
	}
	if !ok {
		return fmt.Errorf("garage bucket %q does not exist", c.bucket)
	}
	return nil
}

func (c *Client) Upload(ctx context.Context, key string, r io.Reader, size int64, contentType string) error {
	_, err := c.mc.PutObject(ctx, c.bucket, key, r, size, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return fmt.Errorf("upload %s: %w", key, err)
	}
	return nil
}

// Delete removes an object. Deleting a key that doesn't exist is not an
// error (Garage/S3 semantics) — callers doing best-effort cleanup don't need
// to special-case "already gone".
func (c *Client) Delete(ctx context.Context, key string) error {
	if err := c.mc.RemoveObject(ctx, c.bucket, key, minio.RemoveObjectOptions{}); err != nil {
		return fmt.Errorf("delete %s: %w", key, err)
	}
	return nil
}

// Download returns the object stream. Call Stat() on it before reading to
// check existence (a missing key only surfaces as an error there, not here).
func (c *Client) Download(ctx context.Context, key string) (*minio.Object, error) {
	obj, err := c.mc.GetObject(ctx, c.bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("download %s: %w", key, err)
	}
	return obj, nil
}

// ListKeys returns the key of every object in the bucket. Used by the
// content export feature to bundle every uploaded file (not just ones
// referenced by a current DB row — an orphaned upload is still something an
// admin backing up "everything" would expect back).
func (c *Client) ListKeys(ctx context.Context) ([]string, error) {
	var keys []string
	for obj := range c.mc.ListObjects(ctx, c.bucket, minio.ListObjectsOptions{Recursive: true}) {
		if obj.Err != nil {
			return nil, fmt.Errorf("list garage objects: %w", obj.Err)
		}
		keys = append(keys, obj.Key)
	}
	return keys, nil
}

func IsNotFound(err error) bool {
	resp := minio.ToErrorResponse(err)
	return resp.Code == "NoSuchKey"
}
