package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/abyssparanoia/rapid-go/internal/domain/errors"
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/domain/repository"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const (
	// urlRoundingDuration is the time rounding unit for presigned URL cache optimization
	urlRoundingDuration = 5 * time.Minute
)

type asset struct {
	cli                *s3.Client
	presignCli         *s3.PresignClient
	privateBucketName  string
	publicBucketName   string
	publicAssetBaseURL string
}

func NewAsset(
	cli *s3.Client,
	privateBucketName string,
	publicBucketName string,
	publicAssetBaseURL string,
) repository.Asset {
	presignCli := s3.NewPresignClient(cli)
	return &asset{
		cli:                cli,
		presignCli:         presignCli,
		privateBucketName:  privateBucketName,
		publicBucketName:   publicBucketName,
		publicAssetBaseURL: strings.TrimSuffix(publicAssetBaseURL, "/"),
	}
}

func (r *asset) GenerateWritePresignedURL(
	ctx context.Context,
	contentType model.ContentType,
	path string,
	expires time.Duration,
) (string, error) {
	// Select bucket based on path prefix
	bucketName := r.privateBucketName
	if strings.HasPrefix(path, "public/") {
		bucketName = r.publicBucketName
	}

	req, err := r.presignCli.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(path),
		ContentType: aws.String(contentType.String()),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expires
	})
	if err != nil {
		return "", errors.InternalErr.Wrap(err)
	}
	return req.URL, nil
}

func (r *asset) GenerateReadURL(
	ctx context.Context,
	path string,
	requestTime time.Time,
) (string, error) {
	// Public: return base URL + path (no signing)
	if strings.HasPrefix(path, "public/") {
		return fmt.Sprintf("%s/%s", r.publicAssetBaseURL, path), nil
	}

	// Private: generate presigned URL with rounded expiration for cache optimization
	// Round request time to 5-minute intervals so same URL is generated within window
	roundedTime := requestTime.Truncate(urlRoundingDuration)
	// Expiration: rounded time + 2 * rounding duration (ensures URL valid for full window)
	expiresAt := roundedTime.Add(urlRoundingDuration * 2)

	req, err := r.presignCli.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(r.privateBucketName),
		Key:    aws.String(path),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Until(expiresAt)
	})
	if err != nil {
		return "", errors.InternalErr.Wrap(err)
	}
	return req.URL, nil
}
