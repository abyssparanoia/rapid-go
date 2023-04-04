package request_interceptor

import (
	"context"
	"fmt"

	"github.com/abyssparanoia/rapid-go/internal/domain/errors"
	"github.com/abyssparanoia/rapid-go/internal/pkg/now"
	"github.com/blendle/zapdriver"
	"github.com/google/uuid"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const producerID = "rapid"

const MAX_RESPONSE_LOG_SIZE = 150 * 1000

type RequestLog struct {
	logger *zap.Logger
}

func NewRequestLog(
	logger *zap.Logger,
) *RequestLog {
	return &RequestLog{
		logger,
	}
}

// UnaryServerInterceptor ...
func (i *RequestLog) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (
		interface{},
		error,
	) {

		operationID := uuid.New()
		now := now.Now()

		ctx = ctxzap.ToContext(ctx, i.logger.With(
			zapdriver.OperationCont(operationID.String(), producerID),
			zap.Time("RequestTime", now),
		))
		ctx = SetRequestTime(ctx, now)

		resp, err := handler(ctx, req)

		if err != nil {
			code := errToCode(err)
			zapcoreLevel := codeToZapCoreLevel(code)
			i.logger.Check(zapcoreLevel, fmt.Sprintf("code: %s  rpc: %s", code.String(), info.FullMethod)).Write(
				zapdriver.OperationEnd(operationID.String(), producerID),
				zap.String("method", info.FullMethod),
				zap.String("grpc.code", code.String()),
				zap.Reflect("request", req),
				zap.Error(err),
			)

			st, err := status.
				New(code, errors.ExtractPlaneErrMessage(err)).
				WithDetails(
					&errdetails.RequestInfo{
						RequestId: operationID.String(),
					},
				)
			if err != nil {
				return nil, err
			}

			return nil, st.Err()
		}

		var logResp interface{}
		if len(fmt.Sprintf("%v", resp)) < MAX_RESPONSE_LOG_SIZE {
			logResp = resp
		}

		i.logger.Info(
			fmt.Sprintf("code: %s  rpc: %s", codes.OK.String(), info.FullMethod),
			zapdriver.OperationEnd(operationID.String(), producerID),
			zap.String("method", info.FullMethod),
			zap.String("grpc.code", codes.OK.String()),
			zap.Reflect("request", req),
			zap.Reflect("response", logResp),
		)

		return resp, nil
	}
}
