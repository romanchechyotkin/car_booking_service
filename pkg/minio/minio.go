package minio

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
)

func New(host, port string) *minio.Client {
	endpoint := fmt.Sprintf("%s:%s", host, port)
	accessKeyID := "bC2fbyLxLUsUHMtqUvDx"
	secretAccessKey := "rQ0EorX8bTLLo75xLn0lIeu9echhwQXEtwuOxhjA"

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
	})
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	log.Printf("%#v\n", minioClient) // minioClient is now set up

	bucketName := "test-bucket"
	location := "BLR"
	ctx := context.Background()

	err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
	}

	return minioClient
}
