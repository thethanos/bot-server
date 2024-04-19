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
		Secure: false,
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

	return nil
}

func (m *MinIOAdapter) PutObject(bucketName, objectName string, file io.Reader, size int64) error {
	options := minio.PutObjectOptions{}
	if _, err := m.client.PutObject(context.Background(), bucketName, objectName, file, size, options); err != nil {
		return err
	}

	return nil
}