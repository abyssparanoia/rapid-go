package requestlog

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/abyssparanoia/rapid-go/internal/pkg/grpcerror"
	"github.com/abyssparanoia/rapid-go/internal/pkg/log"
	"github.com/blendle/zapdriver"
	"github.com/google/uuid"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const producerID = "rapid-go"

// UnaryServerInterceptor ...
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (
		interface{},
		error,
	) {

		l := log.Must(ctx)
		operationID := uuid.New()

		l.Info("call start", zap.Reflect("request", req), zapdriver.OperationStart(operationID.String(), producerID))

		ctxzap.AddFields(
			ctx,
			zapdriver.OperationCont(operationID.String(), producerID),
		)

		resp, err := handler(ctx, req)

		if err != nil {
			code := grpcerror.ErrToCode(err)
			zapcoreLevel := grpcerror.CodeToLevel(code)
			l.Check(zapcoreLevel, "call end").Write(
				zapdriver.OperationEnd(operationID.String(), producerID),
				zap.String("grpc.code", code.String()),
				zap.Error(err),
			)
			err = status.Errorf(code, "%s", err.Error())
			return nil, err
		}

		l.Info(
			"call end",
			zapdriver.OperationEnd(operationID.String(), producerID),
			zap.String("grpc.code", codes.OK.String()),
			zap.Reflect("response", resp),
		)

		return resp, nil
	}
}
