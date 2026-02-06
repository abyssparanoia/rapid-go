package s3

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func NewClient(cfg aws.Config, localhost string) *s3.Client {
	if localhost != "" {
		baseEndpoint := fmt.Sprintf("http://%s", localhost)
		return s3.NewFromConfig(
			cfg,
			func(o *s3.Options) {
				o.UsePathStyle = true
				o.BaseEndpoint = &baseEndpoint
			},
		)
	}
	return s3.NewFromConfig(cfg)
}
