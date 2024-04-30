package minioadapter

import (
	"bot/internal/config"
	"bot/internal/logger"
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIOAdapter struct {
	logger logger.Logger
	cfg    *config.Config
	client *minio.Client
}

func NewMinIOAdapter(logger logger.Logger, cfg *config.Config) (*MinIOAdapter, error) {

	options := &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinIOUser, cfg.MinIOPass, ""),
		Secure: true,
	}

	client, err := minio.New(fmt.Sprintf("%s:%d", cfg.MinIOHost, cfg.MinIOPort), options)
	if err != nil {
		return nil, err
	}

	return &MinIOAdapter{logger: logger, cfg: cfg, client: client}, nil
}

func (m *MinIOAdapter) MakeBucket(bucketName string) error {

	if err := m.client.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{}); err != nil {
		return err
	}

	bucketPolicy := `
	{
		"Version": "2012-10-17",
		"Statement": [
		  {
			"Effect": "Allow",
			"Principal": {
                "AWS": [
                    "*"
                ]
			},
			"Action": [
			  "s3:GetBucketLocation",
			  "s3:ListBucket",
			  "s3:GetObject"
			],
			"Resource": [
			  "arn:aws:s3:::*"
			]
		  }
		]
	  }`

	if err := m.client.SetBucketPolicy(context.Background(), bucketName, bucketPolicy); err != nil {
		return err
	}

	m.logger.Infof("Bucket created: %s", bucketName)
	return nil
}

func (m *MinIOAdapter) PutObject(bucketName, objectName string, file io.Reader, size int64, contentType string) error {
	options := minio.PutObjectOptions{
		ContentType: contentType,
	}

	if _, err := m.client.PutObject(context.Background(), bucketName, objectName, file, size, options); err != nil {
		return err
	}

	m.logger.Infof("Object saved: %s %s", bucketName, objectName)
	return nil
}

func (m *MinIOAdapter) GetBucketObjectsURLs(bucketName string) []string {

	list := make([]string, 0)
	for object := range m.client.ListObjects(context.Background(), bucketName, minio.ListObjectsOptions{}) {
		list = append(list, fmt.Sprintf("%s/%s/%s", m.cfg.ImagePrefix, bucketName, object.Key))
	}
	return list
}

func (m *MinIOAdapter) DeleteObject(bucketName, objectName string) error {
	if err := m.client.RemoveObject(context.Background(), bucketName, objectName, minio.RemoveObjectOptions{}); err != nil {
		return err
	}

	m.logger.Infof("Object %s deletede from bucket %s", objectName, bucketName)
	return nil
}

func (m *MinIOAdapter) DeleteBucket(bucketName string) error {

	objects := m.client.ListObjects(context.Background(), bucketName, minio.ListObjectsOptions{})
	for err := range m.client.RemoveObjects(context.Background(), bucketName, objects, minio.RemoveObjectsOptions{}) {
		return err.Err
	}
	if err := m.client.RemoveBucket(context.Background(), bucketName); err != nil {
		return err
	}

	m.logger.Infof("Bucket deleted: %s", bucketName)
	return nil
}
