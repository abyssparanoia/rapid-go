package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

func NewSession(region string, localhost string) *session.Session {
	if localhost != "" {
		return session.Must(session.NewSessionWithOptions(session.Options{
			Config: aws.Config{Region: aws.String(region), Endpoint: aws.String(localhost)},
		}))
	}
	return session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{Region: aws.String(region)},
	}))
}
