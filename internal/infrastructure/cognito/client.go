package cognito

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

func NewClient(cfg aws.Config, localhost string) *cognitoidentityprovider.Client {
	if localhost != "" {
		return cognitoidentityprovider.NewFromConfig(cfg, func(o *cognitoidentityprovider.Options) {
			o.BaseEndpoint = &localhost
		})
	}
	return cognitoidentityprovider.NewFromConfig(cfg)
}
