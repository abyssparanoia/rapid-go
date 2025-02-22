package request_interceptor

import (
	"context"
	"fmt"

	"github.com/abyssparanoia/goerr"
	"github.com/abyssparanoia/rapid-go/internal/domain/errors"
	"github.com/abyssparanoia/rapid-go/internal/pkg/logger"
	"github.com/abyssparanoia/rapid-go/internal/pkg/logger/logger_field"
	"github.com/abyssparanoia/rapid-go/internal/pkg/now"
	"github.com/blendle/zapdriver"
	"github.com/google/uuid"
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

		ctx = logger.ToContext(
			ctx,
			i.logger.With(
				zapdriver.OperationCont(operationID.String(), producerID),
				zap.Time("RequestTime", now),
			))
		ctx = SetRequestTime(ctx, now)

		resp, err := handler(ctx, req)
		fields := []zap.Field{
			zapdriver.OperationEnd(operationID.String(), producerID),
			zap.String("method", info.FullMethod),
			zap.Any("request", req),
		}
		if err != nil {
			code := errToCode(err)
			fields = append(
				fields,
				zap.String("grpc.code", code.String()),
				logger_field.Error(err),
			)
			zapcoreLevel := codeToZapCoreLevel(code)
			logger.L(ctx).Check(zapcoreLevel, fmt.Sprintf("code: %s  rpc: %s", code.String(), info.FullMethod)).
				Write(fields...)

			errCode, errMessage := extractErrInfo(err)
			st, err := status.
				New(code, errCode).
				WithDetails(
					&errdetails.DebugInfo{
						Detail: errMessage,
					},
					&errdetails.RequestInfo{
						RequestId: operationID.String(),
					},
				)
			if err != nil {
				return nil, errors.InternalErr.Wrap(err)
			}

			return nil, st.Err()
		}

		var logResp interface{}
		if len(fmt.Sprintf("%v", resp)) < MAX_RESPONSE_LOG_SIZE {
			logResp = resp
		}

		fields = append(
			fields,
			zap.String("grpc.code", codes.OK.String()),
			zap.Any("response", logResp),
		)
		logger.L(ctx).Info(
			fmt.Sprintf("code: %s  rpc: %s", codes.OK.String(), info.FullMethod),
			fields...,
		)

		return resp, nil
	}
}

func extractErrInfo(err error) (string, string) {
	if goErr := goerr.Unwrap(err); goErr != nil {
		return goErr.Code(), goErr.Message()
	}
	return "E100001", "An internal error has occurred"
}
