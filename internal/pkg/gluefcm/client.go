package gluefcm

import (
	"context"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"time"
)

// NewClient ... new fcm client
func NewClient(projectID string) *messaging.Client {
	ctx := context.Background()
	gOpt := option.WithGRPCDialOption(grpc.WithKeepaliveParams(keepalive.ClientParameters{
		Time:                30 * time.Millisecond,
		Timeout:             20 * time.Millisecond,
		PermitWithoutStream: true,
	}))
	conf := &firebase.Config{ProjectID: projectID}
	app, err := firebase.NewApp(ctx, conf, gOpt)
	if err != nil {
		panic(err)
	}
	cli, err := app.Messaging(ctx)
	if err != nil {
		panic(err)
	}
	return cli
}
