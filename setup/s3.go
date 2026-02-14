package setup

import (
	"context"
	"errors"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func EnsureBucketExists(ctx context.Context, client *s3.Client, bucketName string) error {
	input := &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	}

	_, err := client.CreateBucket(ctx, input)

	if err != nil {
		var ownedByYou *types.BucketAlreadyOwnedByYou
		var alreadyExists *types.BucketAlreadyExists

		if errors.As(err, &ownedByYou) || errors.As(err, &alreadyExists) {
			log.Printf("Bucket '%s' already exist.\n", bucketName)
			return nil
		}

		return err
	}

	log.Printf("Bucket '%s' generated.\n", bucketName)
	return nil
}
