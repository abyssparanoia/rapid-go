package firebase

import (
	"context"
	"os"
	"time"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func NewClient(projectID string, emulatorHost string) *auth.Client {
	ctx := context.Background()

	// Set emulator host if provided (Firebase Admin SDK reads FIREBASE_AUTH_EMULATOR_HOST env var)
	if emulatorHost != "" {
		os.Setenv("FIREBASE_AUTH_EMULATOR_HOST", emulatorHost)
	}

	var opts []option.ClientOption

	if emulatorHost == "" {
		// Only use aggressive keepalive for production (not needed for local HTTP emulator)
		gOpt := option.WithGRPCDialOption(grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                30 * time.Millisecond,
			Timeout:             20 * time.Millisecond,
			PermitWithoutStream: true,
		}))
		opts = append(opts, gOpt)
	}
	// Note: v4 SDK handles emulator authentication automatically via FIREBASE_AUTH_EMULATOR_HOST

	conf := &firebase.Config{ProjectID: projectID}
	app, err := firebase.NewApp(ctx, conf, opts...)
	if err != nil {
		panic(err)
	}
	cli, err := app.Auth(ctx)
	if err != nil {
		panic(err)
	}
	return cli
}
