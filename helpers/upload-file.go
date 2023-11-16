package helpers

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"time"

	"cloud.google.com/go/storage"
)

type cloudStorage struct {
	projectID  string
	bucketName string
	uploadPath string
}

func NewCloudStorage(projectID string, bucketName string, uploadPath string) CloudStorageInterface {
	return &cloudStorage{
		projectID:  projectID,
		bucketName: bucketName,
		uploadPath: uploadPath,
	}
}

func (uf *cloudStorage) UploadFile(file multipart.File, object string) error {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	client, err := storage.NewClient(ctx)

	if err != nil {
		return err
	}

	// Upload an object with storage.Writer.
	wc := client.Bucket(uf.bucketName).Object(uf.uploadPath + object).NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}

	return nil
}
