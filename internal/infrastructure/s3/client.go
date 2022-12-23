package s3

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func NewClient(session *session.Session) *s3.S3 {
	return s3.New(session)
}
