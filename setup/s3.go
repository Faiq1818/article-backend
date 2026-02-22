package setup

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func EnsureBucketExists(ctx context.Context, client *s3.Client, bucketName string, logger *slog.Logger) error {
	input := &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	}

	_, err := client.CreateBucket(ctx, input)

	if err != nil {
		var ownedByYou *types.BucketAlreadyOwnedByYou
		var alreadyExists *types.BucketAlreadyExists

		if errors.As(err, &ownedByYou) || errors.As(err, &alreadyExists) {
			// log.Printf("Bucket '%s' already exist.\n", bucketName)
			logger.Info("Bucket already exist, skipping create bucket.")
		} else {
			return err
		}
	}

	policy := map[string]any{
		"Version": "2012-10-17",
		"Statement": []map[string]any{
			{
				"Effect":    "Allow",
				"Principal": "*",
				"Action":    "s3:GetObject",
				"Resource":  "arn:aws:s3:::" + bucketName + "/*",
			},
		},
	}

	policyBytes, err := json.Marshal(policy)
	if err != nil {
		return err
	}

	_, err = client.PutBucketPolicy(ctx, &s3.PutBucketPolicyInput{
		Bucket: aws.String(bucketName),
		Policy: aws.String(string(policyBytes)),
	})

	if err != nil {
		log.Printf("Failed to set public policy for bucket '%s': %v\n", bucketName, err)
		return err
	}

	return nil
}
