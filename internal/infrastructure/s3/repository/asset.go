package repository

import (
	"context"
	"time"

	"github.com/abyssparanoia/rapid-go/internal/domain/errors"
	"github.com/abyssparanoia/rapid-go/internal/domain/repository"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type asset struct {
	cli        *s3.Client
	presignCli *s3.PresignClient
	bucketName string
}

func NewAsset(
	cli *s3.Client,
	bucketName string,
) repository.Asset {
	presignCli := s3.NewPresignClient(cli)
	return &asset{
		cli,
		presignCli,
		bucketName,
	}
}

func (r *asset) GenerateWritePresignedURL(
	ctx context.Context,
	contentType string,
	path string,
	expires time.Duration,
) (string, error) {
	req, err := r.presignCli.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(r.bucketName),
		Key:         aws.String(path),
		ContentType: aws.String(contentType),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expires
	})
	if err != nil {
		return "", errors.InternalErr.Wrap(err)
	}
	return req.URL, nil
}

func (r *asset) GenerateReadPresignedURL(
	ctx context.Context,
	path string,
	expires time.Duration,
) (string, error) {
	req, err := r.presignCli.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(path),
	})
	if err != nil {
		return "", errors.InternalErr.Wrap(err)
	}
	return req.URL, nil
}
