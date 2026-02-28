package s3Repo

import (
	"context"
	"errors"
	"io"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3Repository struct {
	S3Client   *s3.Client
	S3Uploader *manager.Uploader
}

func NewS3Repository(s3Client *s3.Client, s3Uploader *manager.Uploader) *S3Repository {
	return &S3Repository{S3Client: s3Client, S3Uploader: s3Uploader}
}

// UploadObject uses the S3 upload manager to upload an object to a bucket.
func (actor S3Repository) UploadObject(ctx context.Context, key string, fileBody io.Reader) (string, error) {
	// get bucket name from env
	bucket := os.Getenv("S3_BUCKET_NAME")

	var outKey string
	input := &s3.PutObjectInput{
		Bucket:            aws.String(bucket),
		Key:               aws.String(key),
		Body:              fileBody,
		ChecksumAlgorithm: types.ChecksumAlgorithmSha256,
	}
	output, err := actor.S3Uploader.Upload(ctx, input)
	if err != nil {
		var noBucket *types.NoSuchBucket
		if errors.As(err, &noBucket) {
			log.Printf("Bucket %s does not exist.\n", bucket)
			err = noBucket
		}
	} else {
		err := s3.NewObjectExistsWaiter(actor.S3Client).Wait(ctx, &s3.HeadObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		}, time.Minute)
		if err != nil {
			log.Printf("Failed attempt to wait for object %s to exist in %s.\n", key, bucket)
		} else {
			outKey = *output.Key
		}
	}
	return outKey, err
}
