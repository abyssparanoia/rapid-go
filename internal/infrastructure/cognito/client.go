package cognito

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

func NewClient(session *session.Session, localhost string) *cognitoidentityprovider.CognitoIdentityProvider {
	cfgs := []*aws.Config{}
	if localhost != "" {
		cfgs = append(cfgs, aws.NewConfig().WithEndpoint(localhost))
	}
	return cognitoidentityprovider.New(session, cfgs...)
}
