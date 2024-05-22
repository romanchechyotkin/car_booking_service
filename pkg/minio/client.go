package minio

import (
	"context"
	"fmt"
	"log"

	"github.com/romanchechyotkin/car_booking_service/pkg/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
)

const bucketName = "images"

type Client struct {
	log         *zap.Logger
	minioClient *minio.Client
}

func New(cfg *config.Config) *minio.Client {
	endpoint := fmt.Sprintf("%s:%s", cfg.Minio.Host, cfg.Minio.Port)
	accessKeyID := cfg.Minio.User
	secretAccessKey := cfg.Minio.Password

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
	})
	if err != nil {
		log.Println("failed to init minio", zap.Error(err))
		return nil
	}

	location := "BLR"
	ctx := context.Background()

	err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Println("We already own", zap.String("bucket", bucketName))
		} else {
			log.Println("failed to init minio", zap.Error(err))
			return nil
		}
	} else {
		log.Println("Successfully created")
	}

	return minioClient
}

func NewClient(cfg *config.Config, log *zap.Logger) *Client {
	endpoint := fmt.Sprintf("%s:%s", cfg.Minio.Host, cfg.Minio.Port)
	accessKeyID := cfg.Minio.User
	secretAccessKey := cfg.Minio.Password

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
	})
	if err != nil {
		log.Error("failed to init minio", zap.Error(err))
		return nil
	}

	location := "BLR"
	ctx := context.Background()

	err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Info("We already own", zap.String("bucket", bucketName))
		} else {
			log.Error("failed to init minio", zap.Error(err))
			return nil
		}
	} else {
		log.Info("Successfully created")
	}

	return &Client{
		log:         log,
		minioClient: minioClient,
	}
}

func (c *Client) Images() *minio.Client {
	return c.minioClient
}
