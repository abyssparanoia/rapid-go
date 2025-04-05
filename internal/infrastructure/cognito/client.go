package cognito

import (
	"context"
	"fmt"

	"github.com/abyssparanoia/rapid-go/internal/domain/errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/lestrrat-go/jwx/v3/jwk"
)

func NewClient(cfg aws.Config, localhost string) *cognitoidentityprovider.Client {
	if localhost != "" {
		return cognitoidentityprovider.NewFromConfig(cfg, func(o *cognitoidentityprovider.Options) {
			o.BaseEndpoint = &localhost
		})
	}
	return cognitoidentityprovider.NewFromConfig(cfg)
}

func NewPublicKeySet(
	ctx context.Context,
	cognitoCli *cognitoidentityprovider.Client,
	userPoolID string,
	emulatorHost string,
	region string,
) (jwk.Set, error) {
	var endpoint string
	if emulatorHost != "" {
		endpoint = emulatorHost
	} else {
		endpoint = fmt.Sprintf("https://cognito-idp.%s.amazonaws.com", region)
	}
	publicKeysURL := fmt.Sprintf("%s/%s/.well-known/jwks.json", endpoint, userPoolID)
	publicKeySet, err := jwk.Fetch(ctx, publicKeysURL)
	if err != nil {
		return nil, errors.InternalErr.Wrap(err)
	}
	return publicKeySet, nil
}
