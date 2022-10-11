package debug

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"firebase.google.com/go/auth"
	debug_apiv1 "github.com/abyssparanoia/rapid-go/schema/proto/pb/rapid/debug_api/v1"
)

type DebugHander struct {
	debug_apiv1.UnimplementedDebugV1ServiceServer
	firebaseAuthCli      *auth.Client
	firebaseClientApiKey string
}

func NewDebugHandler(
	firebaseAuthCli *auth.Client,
	firebaseClientApiKey string,
) debug_apiv1.DebugV1ServiceServer {
	return &DebugHander{
		firebaseAuthCli:      firebaseAuthCli,
		firebaseClientApiKey: firebaseClientApiKey,
	}
}

type verifyCustomTokenResponse struct {
	IDToken string `json:"idToken"`
}

func (h *DebugHander) CreateIDToken(ctx context.Context, req *debug_apiv1.CreateIDTokenRequest) (*debug_apiv1.CreateIDTokenResponse, error) {
	customToken, err := h.firebaseAuthCli.CustomToken(ctx, req.GetAuthUid())
	if err != nil {
		return nil, err
	}

	values := url.Values{}
	values.Add("token", customToken)
	values.Add("returnSecureToken", "true")
	values.Add("key", h.firebaseClientApiKey)

	resp, err := http.Post(
		"https://www.googleapis.com/identitytoolkit/v3/relyingparty/verifyCustomToken",
		"application/x-www-form-urlencoded",
		strings.NewReader(values.Encode()),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var res verifyCustomTokenResponse
	if err := json.Unmarshal(b, &res); err != nil {
		return nil, err
	}

	return &debug_apiv1.CreateIDTokenResponse{
		IdToken: res.IDToken,
	}, nil
}
