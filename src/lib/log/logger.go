package log

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type loggerKey struct{}

// NewLogger ... loggerを生成し、contextに仕込む。
func NewLogger(ctx context.Context, isDev bool) {

	// 開発環境
	if isDev {
		logger, err := zap.NewDevelopment()
		if err != nil {
			panic(err)
		}
		context.WithValue(ctx, loggerKey{}, logger)
		// production環境
	} else {
		logger, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}
		context.WithValue(ctx, loggerKey{}, logger)
	}

}

// Logger ... contextからloggerを取得する
func Logger(ctx context.Context) *zap.Logger {
	logger := ctx.Value(loggerKey{}).(*zap.Logger)

	if logger == nil {
		panic("no logger in context")
	}

	return logger
}

// Debugf ... Debugログを出力する
func Debugf(ctx context.Context, msg string, fields ...zapcore.Field) {
	Logger(ctx).Debug(msg, fields...)
}

// Infof ... Infoログを出力する
func Infof(ctx context.Context, msg string, fields ...zapcore.Field) {
	Logger(ctx).Info(msg, fields...)
}

// Warningf ... Warningログを出力する
func Warningf(ctx context.Context, msg string, fields ...zapcore.Field) {
	Logger(ctx).Warn(msg, fields...)
}

// Errorf ... Errorログを出力する
func Errorf(ctx context.Context, msg string, fields ...zapcore.Field) {
	Logger(ctx).Error(msg, fields...)
}
