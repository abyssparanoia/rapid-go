package s3

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func NewClient(cfg aws.Config, localhost string) *s3.Client {
	if localhost != "" {
		return s3.NewFromConfig(
			cfg,
			func(o *s3.Options) {
				o.UsePathStyle = true
				o.BaseEndpoint = &localhost
			},
		)
	}
	return s3.NewFromConfig(cfg)
}
