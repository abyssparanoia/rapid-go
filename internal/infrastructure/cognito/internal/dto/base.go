package dto

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

func NewUserAttribute(name, value string) *types.AttributeType {
	attr := &types.AttributeType{
		Name:  aws.String(name),
		Value: aws.String(value),
	}
	return attr
}
