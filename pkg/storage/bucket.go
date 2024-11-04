package storage

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	"strings"

	"cloud.google.com/go/storage"
)

func UploadFileToGCS(client *storage.Client, bucketName, objectName string, file *multipart.FileHeader) (string, error) {
	ctx := context.Background()

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	obj := client.Bucket(bucketName).Object(objectName)
	wc := obj.NewWriter(ctx)

	_, err = io.Copy(wc, src)
	if err != nil {
		wc.Close()
		return "", err
	}

	err = wc.Close()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, objectName), nil
}

func DeleteFileFromGCS(client *storage.Client, w io.Writer, bucket, object string) error {
	ctx := context.Background()
	o := client.Bucket(bucket).Object(object)

	attrs, err := o.Attrs(ctx)
	if err != nil {
		return fmt.Errorf("object.Attrs: %w", err)
	}
	o = o.If(storage.Conditions{GenerationMatch: attrs.Generation})

	if err := o.Delete(ctx); err != nil {
		return fmt.Errorf("Object(%q).Delete: %w", object, err)
	}
	fmt.Fprintf(w, "Blob %v deleted.\n", object)
	return nil
}

func DeleteFileFromGCSByURL(client *storage.Client, w io.Writer, fileURL string) error {
	parsedURL, err := url.Parse(fileURL)
	if err != nil {
		return fmt.Errorf("failed to parse URL: %w", err)
	}

	// Example URL: https://storage.googleapis.com/perpus-app/app/book/haiyo/doc_url/memory.pdf
	if parsedURL.Host != "storage.googleapis.com" {
		return fmt.Errorf("unsupported host: %s", parsedURL.Host)
	}

	path := strings.TrimPrefix(parsedURL.Path, "/")
	parts := strings.SplitN(path, "/", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid GCS URL format")
	}

	bucket := parts[0]
	object := parts[1]

	return DeleteFileFromGCS(client, w, bucket, object)
}
