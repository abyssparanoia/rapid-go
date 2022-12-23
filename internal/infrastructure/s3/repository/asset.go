package repository

import (
	"context"
	"time"

	"github.com/abyssparanoia/rapid-go/internal/domain/repository"
	"github.com/abyssparanoia/rapid-go/internal/pkg/errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type asset struct {
	cli        *s3.S3
	bucketName string
}

func NewAsset(
	cli *s3.S3,
	bucketName string,
) repository.Asset {
	return &asset{
		cli,
		bucketName,
	}
}

func (r *asset) GenerateWritePresignedURL(
	ctx context.Context,
	contentType string,
	path string,
	expires time.Duration,
) (string, error) {
	req, _ := r.cli.PutObjectRequest(&s3.PutObjectInput{
		Bucket:      aws.String(r.bucketName),
		Key:         aws.String(path),
		ContentType: aws.String(contentType),
	})

	url, err := req.Presign(expires)
	if err != nil {
		return "", errors.InternalErr.Wrap(err)
	}
	return url, nil
}

func (r *asset) GenerateReadPresignedURL(
	ctx context.Context,
	path string,
	expires time.Duration,
) (string, error) {
	req, _ := r.cli.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(path),
	})

	url, err := req.Presign(expires)
	if err != nil {
		return "", errors.InternalErr.Wrap(err)
	}
	return url, nil
}
